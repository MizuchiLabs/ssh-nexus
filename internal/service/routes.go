package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MizuchiLabs/ssh-nexus/tools/data"
	"github.com/MizuchiLabs/ssh-nexus/tools/updater"
	"github.com/MizuchiLabs/ssh-nexus/tools/util"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"golang.org/x/crypto/ssh"
)

func initRoutes(app core.App) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		api := e.Router.Group("/api")
		api.GET("/version", getVersion)

		authorized := api.Group("", apis.RequireAdminOrRecordAuth())
		authorized.GET(
			"/self/machines",
			func(c echo.Context) error { return getUserMachines(c, app) },
		)

		api.GET("/rpc/certificate", getServerCertificate)
		authorized.GET("/rpc/token", getAgentToken)
		authorized.POST("/rpc/token/rotate", rotateAgentToken)

		api.GET("/ssh/user/public", getPublicKey(data.GetPublicUserKey))
		api.GET("/ssh/host/public", getPublicKey(data.GetPublicHostKey))
		authorized.POST("/ssh/user/set", setUserCA)
		authorized.POST("/ssh/rotate", rotateSSHKeys)
		authorized.POST(
			"/ssh/user/sign",
			func(c echo.Context) error { return signUserCertificate(c, app) },
			apis.ActivityLogger(app),
		)
		api.POST(
			"/ssh/host/sign",
			func(c echo.Context) error { return signHostCertificate(c, app) },
		)

		authorized.POST(
			"/sync/agents",
			func(c echo.Context) error { return forceSync(c, app, syncAgents) },
		)
		authorized.POST(
			"/sync/token",
			func(c echo.Context) error { return forceSync(c, app, syncAgentToken) },
		)
		authorized.POST(
			"/sync/machines",
			func(c echo.Context) error { return forceSync(c, app, syncMachines) },
		)
		authorized.POST(
			"/sync/providers",
			func(c echo.Context) error { return forceSync(c, app, syncProviders) },
		)

		return nil
	})
}

func getVersion(c echo.Context) error {
	return c.JSON(http.StatusOK,
		map[string]string{
			"version": updater.Version,
			"build":   updater.BuildDate,
			"commit":  updater.Commit,
			"branch":  updater.Branch,
		},
	)
}

func getUserMachines(c echo.Context, app core.App) error {
	admin := apis.RequestInfo(c).Admin
	user := apis.RequestInfo(c).AuthRecord

	var machines []Machine

	if admin != nil {
		dbMachines, err := app.Dao().FindRecordsByFilter("machines", "id != ''", "", 0, 0, nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		for _, machine := range dbMachines {
			machineJSON, err := machine.MarshalJSON()
			if err != nil {
				return c.JSON(
					http.StatusInternalServerError,
					map[string]string{"error": err.Error()},
				)
			}

			var data Machine
			if err = json.Unmarshal(machineJSON, &data); err != nil {
				return c.JSON(
					http.StatusInternalServerError,
					map[string]string{"error": err.Error()},
				)
			}
			machines = append(machines, data)
		}
		return c.JSON(
			http.StatusOK,
			map[string]interface{}{"machines": machines},
		)
	}

	if user != nil {
		dbMachines, err := GetUserMachines(app, user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		for _, machine := range dbMachines {
			machineJSON, err := machine.MarshalJSON()
			if err != nil {
				return c.JSON(
					http.StatusInternalServerError,
					map[string]string{"error": err.Error()},
				)
			}

			var data Machine
			if err = json.Unmarshal(machineJSON, &data); err != nil {
				return c.JSON(
					http.StatusInternalServerError,
					map[string]string{"error": err.Error()},
				)
			}
			machines = append(machines, data)
		}

		return c.JSON(
			http.StatusOK,
			map[string]interface{}{"machines": machines},
		)
	}

	return c.JSON(
		http.StatusOK,
		map[string]interface{}{"machines": []Machine{}},
	)
}

func getPublicKey(fetchKey func() ([]byte, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		publicKey, err := fetchKey()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]string{
			"key": strings.TrimSuffix(string(publicKey), "\n"),
		})
	}
}

func setUserCA(c echo.Context) error {
	d := apis.RequestInfo(c).Data
	privateKey := strings.TrimSpace(d["key"].(string))

	_, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return err
	}

	err = os.WriteFile(data.UserKey, []byte(privateKey), 0600)
	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusOK,
		map[string]string{"status": "ok"},
	)
}

func rotateSSHKeys(c echo.Context) error {
	if err := data.GenerateSSHKeys(true); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func getServerCertificate(c echo.Context) error {
	cert, err := data.GetPublicServerCA()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, string(cert))
}

func getAgentToken(c echo.Context) error {
	token, err := data.GetToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func rotateAgentToken(c echo.Context) error {
	if err := data.GenerateToken(true); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func forceSync(c echo.Context, app core.App, syncFunc func(core.App) error) error {
	if err := syncFunc(app); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func signUserCertificate(c echo.Context, app core.App) error {
	admin := apis.RequestInfo(c).Admin
	user := apis.RequestInfo(c).AuthRecord
	d := apis.RequestInfo(c).Data

	// Fetch ttl settings from db
	userTTL, err := app.Dao().FindFirstRecordByData("settings", "key", "user_lease")
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": err.Error()},
		)
	}
	maxTTL, err := app.Dao().FindFirstRecordByData("settings", "key", "max_lease")
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": err.Error()},
		)
	}

	leaseDuration := util.GetLeaseDuration(
		d["ttl"],
		userTTL.GetInt("value"),
		maxTTL.GetInt("value"),
	)

	var principal string
	if admin != nil {
		principal = "root"
	}
	if user != nil {
		principal = user.GetString("principal")
	}

	cert, err := data.SignUserCertificate(
		d["publickey"].(string),
		principal,
		leaseDuration,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": err.Error()},
		)
	}

	return c.JSON(
		http.StatusOK,
		map[string]string{
			"certificate": strings.TrimSuffix(string(cert), "\n"),
			"expiry":      fmt.Sprintf("%d", time.Now().Add(leaseDuration).Unix()),
		},
	)
}

func signHostCertificate(c echo.Context, app core.App) error {
	d := apis.RequestInfo(c).Data

	// Fetch TTL settings from db
	hostTTL, err := app.Dao().FindFirstRecordByData("settings", "key", "host_lease")
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": err.Error()},
		)
	}
	maxTTL, err := app.Dao().FindFirstRecordByData("settings", "key", "max_lease")
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": err.Error()},
		)
	}

	leaseDuration := util.GetLeaseDuration(
		d["ttl"],
		hostTTL.GetInt("value"),
		maxTTL.GetInt("value"),
	)

	cert, err := data.SignUserCertificate(
		d["publickey"].(string),
		d["hostname"].(string),
		leaseDuration,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": err.Error()},
		)
	}

	return c.JSON(
		http.StatusOK,
		map[string]string{
			"certificate": strings.TrimSuffix(string(cert), "\n"),
			"expiry":      fmt.Sprintf("%d", time.Now().Add(leaseDuration).Unix()),
		},
	)
}
