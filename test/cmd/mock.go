// This is just a script to generate mock data for the database
package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"syscall"
	"time"

	_ "github.com/MizuchiLabs/ssh-nexus/internal/migrations"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

var (
	numRecords  = 25
	collections = []string{
		"groups",
		"permissions",
		"users",
		"machines",
		"tags",
		"settings",
		"providers",
	}
	providers = []string{
		"aws",
		"azure",
		"google",
		"hetzner",
		"linode",
		"vultr",
		"proxmox",
		"digitalocean",
	}
)

func randomSubset(ids []string, min int, max int) []string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	n := rand.Intn(max-min+1) + min
	subset := make([]string, n)
	for i := 0; i < n; i++ {
		subset[i] = ids[rand.Intn(len(ids))]
	}
	return subset
}

func generateMockData(
	app *pocketbase.PocketBase,
	collectionName string,
	numRecords int,
) ([]string, error) {
	collection, err := app.Dao().FindCollectionByNameOrId(collectionName)
	if err != nil {
		return nil, err
	}

	var recordIDs []string
	for i := 0; i < numRecords; i++ {
		record := models.NewRecord(collection)
		switch collectionName {
		case "users":
			if err = record.SetPassword("testing1234"); err != nil {
				return nil, err
			}
			record.Set("email", gofakeit.Email())
			record.Set("emailVisibility", gofakeit.Bool())
			record.Set("name", gofakeit.Name())
			record.Set("passwordHash", gofakeit.UUID())
			record.Set("tokenKey", gofakeit.UUID())
			record.Set("username", gofakeit.Username())
			record.Set("verified", gofakeit.Bool())
			if rand.Intn(2) == 0 {
				record.Set("principal", gofakeit.UUID())
			} else {
				record.Set("principal", "")
			}
			record.Set("settings", gofakeit.Map())
		case "machines":
			record.Set("name", gofakeit.AppName())
			record.Set("host", gofakeit.IPv4Address())
			record.Set("port", gofakeit.Number(1, 65535))
			record.Set("agent", gofakeit.Bool())
			record.Set("error", gofakeit.Error())
		case "groups":
			record.Set("name", gofakeit.JobLevel()+gofakeit.UUID()[:8])
			record.Set("description", gofakeit.Sentence(10))
			record.Set("linux_username", strings.ToLower(gofakeit.Username()))
		case "tags":
			record.Set("name", gofakeit.Gamertag())
			record.Set("description", gofakeit.Sentence(10))
		case "permissions":
			record.Set("name", gofakeit.JobLevel()+gofakeit.UUID()[:8])
			record.Set("description", gofakeit.Sentence(10))
			record.Set("can_create", gofakeit.Bool())
			record.Set("can_update", gofakeit.Bool())
			record.Set("can_delete", gofakeit.Bool())
			record.Set("access_users", gofakeit.Bool())
			record.Set("access_machines", gofakeit.Bool())
			record.Set("access_groups", gofakeit.Bool())
			record.Set("is_admin", gofakeit.Bool())
		case "settings":
			record.Set("key", gofakeit.Name())
			record.Set("value", gofakeit.Word())
		case "providers":
			record.Set("name", gofakeit.AppName())
			record.Set("url", gofakeit.URL())
			record.Set("username", gofakeit.Username())
			record.Set("password", gofakeit.Password(true, true, true, true, true, 8))
			record.Set("token", gofakeit.UUID())
			record.Set("error", gofakeit.Error())
			record.Set("last_sync", time.Now().Format(time.RFC3339))
			record.Set("type", providers[rand.Intn(len(providers))])
		}

		err = app.Dao().Save(record)
		if err != nil {
			return nil, err
		}
		recordIDs = append(recordIDs, record.GetId())
	}
	return recordIDs, nil
}

func generateRelationships(app *pocketbase.PocketBase, recordIDs map[string][]string) error {
	// Updating relational fields with proper references
	for _, userID := range recordIDs["users"] {
		user, err := app.Dao().FindRecordById("users", userID)
		if err != nil {
			log.Fatal(err)
		}
		user.Set("permission", recordIDs["permissions"][rand.Intn(len(recordIDs["permissions"]))])
		user.Set("groups", randomSubset(recordIDs["groups"], 3, 5))

		if err := app.Dao().Save(user); err != nil {
			log.Fatal(err)
		}
	}

	for _, machineID := range recordIDs["machines"] {
		machine, err := app.Dao().FindRecordById("machines", machineID)
		if err != nil {
			log.Fatal(err)
		}
		machine.Set("users", randomSubset(recordIDs["users"], 1, 5))
		machine.Set("groups", randomSubset(recordIDs["groups"], 1, 5))
		machine.Set("tags", randomSubset(recordIDs["tags"], 1, 3))
		machine.Set("provider", recordIDs["providers"][rand.Intn(len(recordIDs["providers"]))])

		if err := app.Dao().Save(machine); err != nil {
			log.Fatal(err)
		}
	}

	for _, permissionID := range recordIDs["permissions"] {
		permission, err := app.Dao().FindRecordById("permissions", permissionID)
		if err != nil {
			log.Fatal(err)
		}
		permission.Set("machines", randomSubset(recordIDs["machines"], 1, 5))
		permission.Set("groups", randomSubset(recordIDs["groups"], 1, 5))
		permission.Set("users", randomSubset(recordIDs["users"], 1, 5))

		if err := app.Dao().Save(permission); err != nil {
			log.Fatal(err)
		}
	}

	for _, userID := range recordIDs["users"] {
		requests, err := app.Dao().FindCollectionByNameOrId("requests")
		if err != nil {
			log.Fatal(err)
		}
		record := models.NewRecord(requests)

		record.Set("user", userID)
		record.Set("description", gofakeit.Sentence(10))
		if rand.Intn(2) == 0 {
			record.Set("machine", recordIDs["machines"][rand.Intn(len(recordIDs["machines"]))])
		} else {
			record.Set("group", recordIDs["groups"][rand.Intn(len(recordIDs["groups"]))])
		}

		if err := app.Dao().Save(record); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Relational fields updated successfully!")

	return nil
}

func setupAdminAccount(app *pocketbase.PocketBase) error {
	adminEmail := "test@example.com"
	_, err := app.Dao().FindAdminByEmail(adminEmail)
	if err != nil {
		admin := &models.Admin{}
		admin.Email = adminEmail
		if err = admin.SetPassword("testing1234"); err != nil {
			return err
		}

		if err = app.Dao().SaveAdmin(admin); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	err := os.RemoveAll("./test_pb_data")
	if err != nil {
		log.Fatal(err)
	}

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: "./test_pb_data",
	})
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		recordIDs := make(map[string][]string)
		for _, collection := range collections {
			ids, err := generateMockData(app, collection, numRecords)
			if err != nil {
				log.Fatal(err)
			}
			recordIDs[collection] = ids
			fmt.Printf("Mock data for collection %s generated successfully!\n", collection)
		}

		err := generateRelationships(app, recordIDs)
		if err != nil {
			log.Fatal(err)
		}
		return setupAdminAccount(app)
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		return nil
	})

	os.Args = append(os.Args, "serve")
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
