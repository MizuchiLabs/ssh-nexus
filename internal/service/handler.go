// Package service contains all the service handlers
package service

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"slices"

	"github.com/MizuchiLabs/ssh-nexus/api/server"
	"github.com/MizuchiLabs/ssh-nexus/internal/config"
	"github.com/MizuchiLabs/ssh-nexus/tools/data"
	"github.com/MizuchiLabs/ssh-nexus/tools/updater"
	"github.com/MizuchiLabs/ssh-nexus/tools/util"
	"github.com/MizuchiLabs/ssh-nexus/web"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/spf13/cobra"
)

func Server() error {
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir:       "./pb_data",
		DefaultEncryptionEnv: "PB_ENCRYPTION_KEY",
	})

	app.RootCmd.Short = "SSH Nexus CLI"
	app.RootCmd.Version = fmt.Sprintf("%s (%s)", updater.Version, updater.BuildDate)
	app.RootCmd.AddCommand(&cobra.Command{
		Use:     "update",
		Aliases: []string{"up"},
		Short:   "Update the server to the latest version",
		Run: func(cmd *cobra.Command, args []string) {
			updater.UpdateSelf(updater.Version, true)
		},
	})

	if err := AppEventHandler(app.App); err != nil {
		return err
	}
	if err := AuditEventHandler(app.App); err != nil {
		return err
	}
	if err := UserEventHandler(app.App); err != nil {
		return err
	}
	if err := MachineEventHandler(app.App); err != nil {
		return err
	}

	if len(os.Args) <= 1 {
		os.Args = append(os.Args, "serve")
		os.Args = append(os.Args, "--http=0.0.0.0:8090")
	}
	if err := app.Start(); err != nil {
		slog.Error("backend server error", "err", err)
		return err
	}
	return nil
}

func KeyCheck(app core.App) error {
	settings, err := app.Dao().FindSettings(os.Getenv("PB_ENCRYPTION_KEY"))
	if err != nil {
		return err
	}

	if err := data.GenerateSSHKeys(false); err != nil {
		return fmt.Errorf("ssh key generation failed: %w", err)
	}
	if err := data.GenerateTLS(settings.Meta.AppUrl); err != nil {
		return fmt.Errorf("tls generation failed: %w", err)
	}
	if err := data.RegenerateKeys(settings.Meta.AppUrl); err != nil {
		return fmt.Errorf("tls regeneration failed: %w", err)
	}
	if err := data.GenerateToken(false); err != nil {
		return fmt.Errorf("token generation failed: %w", err)
	}

	return nil
}

func AppEventHandler(app core.App) error {
	initRoutes(app)

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Check for new version
		go updater.UpdateSelf(updater.Version, false)

		// Update default settings used by ssh-nexus
		if err := config.UpdateSettings(app); err != nil {
			return err
		}

		// Web UI
		e.Router.GET("/*", apis.StaticDirectoryHandler(web.Static, true))

		// Setup key checks
		if err := KeyCheck(app); err != nil {
			return err
		}

		// Start gRPC server
		server.Server(app)

		// Setup tasks and schedule them
		util.Execute(func() { syncProviders(app) })
		util.Execute(func() { cleanupAudit(app) })

		scheduler := cron.New()
		scheduler.MustAdd("Sync Providers", "0 * * * *", func() { // every hour
			util.Execute(func() { syncProviders(app) })
		})
		scheduler.MustAdd("Cleanup Auditlog", "0 0 * * 0", func() { // sunday midnight
			util.Execute(func() { cleanupAudit(app) })
		})
		scheduler.Start()
		return nil
	})

	// Dynamically regenerate cert in case the domain changes
	app.OnSettingsBeforeUpdateRequest().Add(func(e *core.SettingsUpdateEvent) error {
		if e.NewSettings.Meta.AppUrl != e.OldSettings.Meta.AppUrl {
			appURL, err := url.Parse(e.NewSettings.Meta.AppUrl)
			if err != nil {
				return err
			}
			if err = data.RegenerateKeys(appURL.Hostname()); err != nil {
				return err
			}
		}
		return nil
	})

	// Disconnect agents on shutdown
	// app.OnTerminate().Add(func(e *core.TerminateEvent) error {
	// 	machines, err := app.Dao().FindRecordsByFilter("machines", "agent = true", "", 0, 0, nil)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	for _, machine := range machines {
	// 		machine.Set("agent", false)
	// 		if err := app.Dao().SaveRecord(machine); err != nil {
	// 			slog.Error("failed to update agent status", "err", err)
	// 		}
	// 	}
	// 	return nil
	// })

	return nil
}

func AuditEventHandler(app core.App) error {
	app.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		return insertAuditlog(app, e.HttpContext, e.Record, "create")
	})

	app.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		return insertAuditlog(app, e.HttpContext, e.Record, "update")
	})

	app.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		return insertAuditlog(app, e.HttpContext, e.Record, "delete")
	})
	return nil
}

func UserEventHandler(app core.App) error {
	app.OnRecordAfterCreateRequest("users").Add(func(e *core.RecordCreateEvent) error {
		return setPrincipalUUID(app, e.Record)
	})
	app.OnRecordAfterUpdateRequest("users").Add(func(e *core.RecordUpdateEvent) error {
		if err := setPrincipalUUID(app, e.Record); err != nil {
			return err
		}

		// Update machines if groups or principal id changes
		groups := e.Record.GetStringSlice("groups")
		oldGroups := e.Record.OriginalCopy().GetStringSlice("groups")
		principal := e.Record.GetString("principal")
		oldPrincipal := e.Record.OriginalCopy().GetString("principal")
		if !slices.Equal(groups, oldGroups) || (principal != oldPrincipal) {
			machinesBefore, err := GetUserMachines(app, e.Record.OriginalCopy())
			if err != nil {
				return err
			}

			machinesAfter, err := GetUserMachines(app, e.Record)
			if err != nil {
				return err
			}
			machines := append(machinesBefore, machinesAfter...)
			for _, machine := range machines {
				util.Execute(func() { ManualUpdate(app, machine) })
			}
			fmt.Println()
		}

		return nil
	})
	app.OnRecordAfterAuthWithPasswordRequest("users").
		Add(func(e *core.RecordAuthWithPasswordEvent) error {
			return setPrincipalUUID(app, e.Record)
		})
	app.OnRecordAfterAuthWithOAuth2Request("users").
		Add(func(e *core.RecordAuthWithOAuth2Event) error {
			return setPrincipalUUID(app, e.Record)
		})
	return nil
}

func MachineEventHandler(app core.App) error {
	app.OnRecordBeforeCreateRequest("machines").
		Add(func(e *core.RecordCreateEvent) error {
			// User isn't allowed to manually switch agent status
			if e.Record.GetBool("agent") {
				e.Record.Set("agent", false)
				if err := app.Dao().SaveRecord(e.Record); err != nil {
					slog.Error("failed to save machine", "err", err)
				}
			}

			install, err := app.Dao().FindFirstRecordByData("settings", "key", "install_agent")
			if err != nil {
				return err
			}

			// Install agent if not already installed and enabled
			if install.GetBool("value") {
				util.Execute(func() { InstallAgent(app, e.Record) })
			}
			return nil
		})

	app.OnRecordAfterUpdateRequest("machines").
		Add(func(e *core.RecordUpdateEvent) error {
			// Manual update if agent is not connected
			util.Execute(func() { ManualUpdate(app, e.Record) })
			if err := cleanupTags(app); err != nil {
				return err
			}
			return nil
		})

	app.OnRecordBeforeDeleteRequest("machines").
		Add(func(e *core.RecordDeleteEvent) error {
			// Manual restore if agent is not connected
			util.Execute(func() { Restore(e.Record) })
			if err := cleanupTags(app); err != nil {
				return err
			}
			return nil
		})

	// Sync Providers on create/update
	app.OnRecordAfterCreateRequest("providers").
		Add(func(e *core.RecordCreateEvent) error {
			util.Execute(func() { syncProviders(app) })
			return nil
		})

	app.OnRecordAfterUpdateRequest("providers").
		Add(func(e *core.RecordUpdateEvent) error {
			util.Execute(func() { syncProviders(app) })
			return nil
		})

	return nil
}
