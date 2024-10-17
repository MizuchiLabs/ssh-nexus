package service

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type User struct {
	ID              string   `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	Username        string   `json:"username,omitempty"`
	Principal       string   `json:"principal,omitempty"`
	Avatar          string   `json:"avatar,omitempty"`
	Email           string   `json:"email,omitempty"`
	EmailVisibility bool     `json:"emailVisibility,omitempty"`
	Verified        bool     `json:"verified,omitempty"`
	Permission      string   `json:"permission,omitempty"`
	Groups          []string `json:"groups,omitempty"`
}

type Machine struct {
	ID       string   `json:"id,omitempty"`
	Name     string   `json:"name,omitempty"`
	Host     string   `json:"host,omitempty"`
	Port     int      `json:"port,omitempty"`
	Agent    bool     `json:"agent,omitempty"`
	Error    string   `json:"error,omitempty"`
	Provider string   `json:"provider,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Users    []string `json:"users,omitempty"`
	Groups   []string `json:"groups,omitempty"`
}

type Group struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Username    string `json:"linux_username,omitempty"`
}

type Tag struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type Permission struct {
	ID             string   `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description,omitempty"`
	CanCreate      bool     `json:"can_create"`
	CanUpdate      bool     `json:"can_update"`
	CanDelete      bool     `json:"can_delete"`
	AccessUsers    bool     `json:"access_users"`
	AccessMachines bool     `json:"access_machines"`
	AccessGroups   bool     `json:"access_groups"`
	IsAdmin        bool     `json:"is_admin"`
	Users          []string `json:"users,omitempty"`
	Groups         []string `json:"groups,omitempty"`
	Machines       []string `json:"machines,omitempty"`
}

func setPrincipalUUID(app core.App, record *models.Record) error {
	if record.Collection().Name != "users" {
		return fmt.Errorf("wrong collection %s", record.Collection().Name)
	}

	if record.GetString("principal") == "" {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		record.Set("principal", uuid.String())
		if err := app.Dao().SaveRecord(record); err != nil {
			return err
		}
	}

	return nil
}

func cleanupTags(app core.App) error {
	var tags []struct {
		ID string `json:"id"`
	}
	err := app.Dao().DB().NewQuery(`
		DELETE FROM tags 
		WHERE id NOT IN (
			SELECT DISTINCT json_each.value
			FROM machines
			JOIN json_each(machines.tags))
		`).
		All(&tags)
	if err != nil {
		return err
	}

	return nil
}

func insertAuditlog(
	app core.App,
	httpContext echo.Context,
	record *models.Record,
	event string,
) error {
	// ignore auditlog record changes
	if !slices.Contains(
		[]string{"users", "machines", "groups", "settings"},
		record.Collection().Name,
	) {
		return nil
	}

	auditlogCollection, err := app.Dao().FindCollectionByNameOrId("auditlog")
	if err != nil {
		return err
	}

	var adminID string
	var authRecordID string

	// get the authenticated admin
	admin, _ := httpContext.Get(apis.ContextAdminKey).(*models.Admin)
	if admin != nil {
		adminID = admin.Id
	}

	// get the authenticated user record
	authRecord, _ := httpContext.Get(apis.ContextAuthRecordKey).(*models.Record)
	if authRecord != nil {
		authRecordID = authRecord.Id
	}
	if event != "delete" {
		if record.Collection().Name == "users" {
			var user1, user2 User
			user1JSON, _ := record.CleanCopy().MarshalJSON()
			json.Unmarshal(user1JSON, &user1)
			user2JSON, _ := record.OriginalCopy().MarshalJSON()
			json.Unmarshal(user2JSON, &user2)

			// ignore auditlog record changes if same ignore updated and created
			if reflect.DeepEqual(user1, user2) {
				return nil
			}

			if errs := app.Dao().ExpandRecord(record, []string{"permission", "groups"}, nil); len(
				errs,
			) > 0 {
				return fmt.Errorf("failed to expand: %v", errs)
			}
			if errs := app.Dao().ExpandRecord(record.OriginalCopy(), []string{"permission", "groups"}, nil); len(
				errs,
			) > 0 {
				return fmt.Errorf("failed to expand: %v", errs)
			}
		} else if record.Collection().Name == "machines" {
			var machine1, machine2 Machine
			machine1JSON, _ := record.CleanCopy().MarshalJSON()
			json.Unmarshal(machine1JSON, &machine1)
			machine2JSON, _ := record.OriginalCopy().MarshalJSON()
			json.Unmarshal(machine2JSON, &machine2)

			if reflect.DeepEqual(machine1, machine2) {
				return nil
			}

			if errs := app.Dao().ExpandRecord(record, []string{"users", "groups"}, nil); len(errs) > 0 {
				return fmt.Errorf("failed to expand: %v", errs)
			}
			if errs := app.Dao().ExpandRecord(record.OriginalCopy(), []string{"users", "groups"}, nil); len(errs) > 0 {
				return fmt.Errorf("failed to expand: %v", errs)
			}
		} else if record.Collection().Name == "groups" {
			var group1, group2 Group
			group1JSON, _ := record.CleanCopy().MarshalJSON()
			json.Unmarshal(group1JSON, &group1)
			group2JSON, _ := record.OriginalCopy().MarshalJSON()
			json.Unmarshal(group2JSON, &group2)

			if reflect.DeepEqual(group1, group2) {
				return nil
			}
		} else if record.Collection().Name == "permissions" {
			var permission1, permission2 Permission
			permission1JSON, _ := record.CleanCopy().MarshalJSON()
			json.Unmarshal(permission1JSON, &permission1)
			permission2JSON, _ := record.OriginalCopy().MarshalJSON()
			json.Unmarshal(permission2JSON, &permission2)

			if reflect.DeepEqual(permission1, permission2) {
				return nil
			}

			if errs := app.Dao().ExpandRecord(record, []string{"users", "machines", "groups"}, nil); len(errs) > 0 {
				return fmt.Errorf("failed to expand: %v", errs)
			}
			if errs := app.Dao().ExpandRecord(record.OriginalCopy(), []string{"users", "machines", "groups"}, nil); len(errs) > 0 {
				return fmt.Errorf("failed to expand: %v", errs)
			}
		}
	}

	auditlog := models.NewRecord(auditlogCollection)
	auditlog.Set("collection", record.Collection().Name)
	auditlog.Set("record", record.Id)
	auditlog.Set("event", event)
	auditlog.Set("user", authRecordID)
	auditlog.Set("admin", adminID)
	auditlog.Set("data", record)
	auditlog.Set("original", record.OriginalCopy())

	return app.Dao().SaveRecord(auditlog)
}
