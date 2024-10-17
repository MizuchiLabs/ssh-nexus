package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/MizuchiLabs/ssh-nexus/internal/provider"
	"github.com/MizuchiLabs/ssh-nexus/tools/data"
	"github.com/MizuchiLabs/ssh-nexus/tools/util"
	"github.com/pkg/sftp"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

// GetMachineUsers fetches all principals based on a machine
func GetMachineUsers(
	app core.App,
	machine *models.Record,
) (map[string][]string, error) {
	if machine == nil {
		return nil, fmt.Errorf("no machine provided")
	}

	data := make(map[string][]string)
	if err := app.Dao().ExpandRecord(machine, []string{"groups", "users"}, nil); len(err) > 0 {
		return nil, fmt.Errorf("failed to expand: %v", err)
	}

	for _, group := range machine.ExpandedAll("groups") {
		users, err := app.Dao().
			FindRecordsByFilter("users", "groups.id ?= {:group_id}", "", 0, 0, dbx.Params{"group_id": group.Id})
		if err != nil {
			return nil, fmt.Errorf("failed to find users: %v", err)
		}

		for _, user := range users {
			data[group.GetString("linux_username")] = append(
				data[group.GetString("linux_username")],
				user.GetString("principal"),
			)
		}
	}
	for _, user := range machine.ExpandedAll("users") {
		data["root"] = append(data["root"], user.GetString("principal"))
	}

	if !slices.Contains(data["root"], "root") {
		data["root"] = append([]string{"root"}, data["root"]...)
	}

	return data, nil
}

// GetUserMachines fetches all machines based on a user
func GetUserMachines(
	app core.App,
	user *models.Record,
) ([]*models.Record, error) {
	if user == nil {
		return nil, fmt.Errorf("no user provided")
	}

	groupsJSON, err := json.Marshal(user.GetStringSlice("groups"))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user groups: %v", err)
	}

	machineIds := []struct {
		ID string `db:"id" json:"id"`
	}{}
	err = app.Dao().DB().NewQuery(`
		SELECT DISTINCT machines.id
		FROM machines
		LEFT JOIN json_each(machines.users) AS machine_user
		LEFT JOIN json_each(machines.groups) AS machine_group
		WHERE machine_user.value = {:user_id}
		OR machine_group.value IN (SELECT value FROM json_each({:groups}));
    `).Bind(dbx.Params{
		"user_id": user.Id,
		"groups":  string(groupsJSON),
	}).All(&machineIds)
	if err != nil {
		return nil, fmt.Errorf("failed to find machines: %v", err)
	}

	if len(machineIds) == 0 {
		return nil, nil
	}

	ids := make([]string, len(machineIds))
	for i, machineID := range machineIds {
		ids[i] = machineID.ID
	}

	machines, err := app.Dao().FindRecordsByIds("machines", ids)
	if err != nil {
		return nil, fmt.Errorf("failed to find machines: %v", err)
	}
	return machines, nil
}

// syncMachines forces a manual update of all machines which are not connected to an agent
func syncMachines(app core.App) error {
	machines, err := app.Dao().
		FindRecordsByFilter("machines", "agent = false", "", 0, 0, nil)
	if err != nil {
		return err
	}
	if len(machines) == 0 {
		return nil
	}

	for _, machine := range machines {
		util.Execute(func() { ManualUpdate(app, machine) })
	}
	return nil
}

// syncAgents tries to install agents on all machines
func syncAgents(app core.App) error {
	machines, err := app.Dao().
		FindRecordsByFilter("machines", "agent = false", "", 0, 0, nil)
	if err != nil {
		return err
	}
	if len(machines) == 0 {
		return nil
	}
	for _, machine := range machines {
		util.Execute(func() { InstallAgent(app, machine) })
	}
	return nil
}

// syncAgentToken tries to send the agent token to the machine in case it was rotated
func syncAgentToken(app core.App) error {
	machines, err := app.Dao().FindRecordsByFilter("machines", "id != ''", "", 0, 0, nil)
	if err != nil {
		return err
	}
	for _, machine := range machines {
		defer func() {
			if err := app.Dao().SaveRecord(machine); err != nil {
				slog.Error("Failed to save machine", "name", machine.GetString("name"), "err", err)
			}
		}()

		conn, err := connect(machine)
		if err != nil {
			machine.Set("error", err.Error())
			return err
		}
		defer conn.Close()

		token, err := os.Open(data.Token)
		if err != nil {
			machine.Set("error", err.Error())
			return err
		}
		defer token.Close()

		client, err := sftp.NewClient(conn)
		if err != nil {
			machine.Set("error", err.Error())
			return err
		}
		defer client.Close()

		remoteToken, err := client.Create(data.Token)
		if err != nil {
			machine.Set("error", err.Error())
			return err
		}
		defer remoteToken.Close()

		_, err = io.Copy(remoteToken, token)
		if err != nil {
			machine.Set("error", err.Error())
			return err
		}
	}

	return nil
}

// syncProviders syncs cloud providers
func syncProviders(app core.App) error {
	providers, err := app.Dao().FindRecordsByFilter("providers", "type != ''", "", 0, 0)
	if err != nil {
		return err
	}

	if len(providers) == 0 {
		return nil
	}

	machines, err := app.Dao().FindCollectionByNameOrId("machines")
	if err != nil {
		return err
	}

	tags, err := app.Dao().FindCollectionByNameOrId("tags")
	if err != nil {
		return err
	}

	go func() {
		for _, p := range providers {
			lastSync := p.GetTime("last_sync")

			// Skip if last sync was less than 2 minutes ago
			if lastSync.After(time.Now().Add(-2 * time.Minute)) {
				continue
			}

			provider, err := provider.NewProvider(p)
			if err != nil {
				p.Set("error", err.Error())
				if err = app.Dao().SaveRecord(p); err != nil {
					slog.Error("failed to save provider", "err", err)
				}
				continue
			}

			pMachines, err := provider.Sync()
			if err != nil {
				p.Set("error", err.Error())
				if err := app.Dao().SaveRecord(p); err != nil {
					slog.Error("failed to save provider", "err", err)
				}
				continue
			} else if len(pMachines) == 0 {
				p.Set("error", "no machines found")
				if err := app.Dao().SaveRecord(p); err != nil {
					slog.Error("failed to save provider", "err", err)
				}
				continue
			}

			// Update last sync time
			p.Set("error", "")
			p.Set("last_sync", time.Now())
			if err := app.Dao().SaveRecord(p); err != nil {
				slog.Error("failed to save provider", "err", err)
				continue
			}

			// Add tag if it doesn't exist
			tag, _ := app.Dao().FindFirstRecordByData("tags", "name", p.GetString("type"))
			if tag == nil {
				tag = models.NewRecord(tags)
				tag.Set("name", p.GetString("type"))
				if err := app.Dao().SaveRecord(tag); err != nil {
					slog.Error("failed to save tag", "err", err)
				}
			}

			for _, pm := range pMachines {
				machine := models.NewRecord(machines)
				machine.Set("name", pm.Name)
				machine.Set("host", pm.Host)
				machine.Set("port", "22")
				machine.Set("provider", p.Id)
				machine.Set("tags", tag.Id)
				if err := app.Dao().SaveRecord(machine); err != nil {
					slog.Error("failed to save machine", "err", err)
					continue
				}
			}
		}
	}()

	return nil
}

// cleanupAudit cleans up audit log based on retention
func cleanupAudit(app core.App) error {
	settings, err := app.Dao().FindFirstRecordByData("settings", "key", "default_retention")
	if err != nil {
		return fmt.Errorf("failed to get settings: %v", err)
	}

	var retention time.Time
	seconds, _ := strconv.Atoi(settings.GetString("value"))
	if seconds == 0 {
		retention = time.Now().AddDate(0, 0, -30).UTC()
	} else {
		retention = time.Now().UTC().Add(-time.Duration(seconds) * time.Second)
	}

	records, _ := app.Dao().
		FindRecordsByFilter("auditlog", "created <= {:created}", "-created", 0, 0, dbx.Params{"created": retention})
	for _, record := range records {
		if err := app.Dao().DeleteRecord(record); err != nil {
			slog.Error("failed to delete record", "err", err)
			continue // ignore
		}
	}

	return nil
}
