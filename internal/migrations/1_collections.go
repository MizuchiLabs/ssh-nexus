package migrations

import (
	"log/slog"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collections := []string{
			"machines",
			"users",
			"groups",
			"tags",
			"permissions",
			"requests",
			"auditlog",
			"settings",
			"providers",
		}

		// Create collections
		ids := make(map[string]string)
		for _, name := range collections {
			collection, _ := dao.FindCollectionByNameOrId(name)
			if collection == nil {
				collection = &models.Collection{
					Name: name,
					Type: models.CollectionTypeBase,
				}
			}
			if err := dao.SaveCollection(collection); err != nil {
				return err
			}
			ids[name] = collection.Id
		}

		// Machines Collection
		err := initCollection(
			dao,
			"machines",
			"@request.auth.id != ''", // List Rule
			"@request.auth.id != ''", // View Rule
			"@request.auth.permission.is_admin = true || "+
				"(@request.auth.permission.access_machines = true && "+
				"@request.auth.permission.can_create = true)", // Create Rule
			"@request.auth.permission.is_admin = true || "+
				"(@request.auth.permission.can_update = true && "+
				"(@request.auth.permission.access_machines = true || "+
				"@request.auth.permission.machines.id ?= id))", // Update Rule
			"@request.auth.permission.is_admin = true || "+
				"(@request.auth.permission.can_delete = true && "+
				"(@request.auth.permission.access_machines = true || "+
				"@request.auth.permission.machines.id ?= id))", // Delete Rule
			types.JsonArray[string]{
				"CREATE UNIQUE INDEX idx_machines_name_provider ON machines(name,provider)",
				"CREATE UNIQUE INDEX idx_machines_host ON machines(host)",
			},
			&schema.SchemaField{
				Name:        "name",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: true,
			},
			&schema.SchemaField{
				Name:     "host",
				Type:     schema.FieldTypeText,
				Required: true,
				Options: &schema.TextOptions{
					Pattern: `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$|^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+(?:[a-zA-Z]{2,})$`, // Checks if either ip or domain
				},
			},
			&schema.SchemaField{
				Name:     "port",
				Type:     schema.FieldTypeNumber,
				Required: true,
				Options: &schema.NumberOptions{
					Max:       types.Pointer[float64](65535),
					NoDecimal: true,
				},
			},
			&schema.SchemaField{
				Name:     "agent",
				Type:     schema.FieldTypeBool,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "uuid",
				Type:     schema.FieldTypeText,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "error",
				Type:     schema.FieldTypeText,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "users",
				Type:     schema.FieldTypeRelation,
				Required: false,
				Options: &schema.RelationOptions{
					CollectionId: ids["users"],
				},
			},
			&schema.SchemaField{
				Name:     "groups",
				Type:     schema.FieldTypeRelation,
				Required: false,
				Options: &schema.RelationOptions{
					CollectionId: ids["groups"],
				},
			},
			&schema.SchemaField{
				Name:     "tags",
				Type:     schema.FieldTypeRelation,
				Required: false,
				Options: &schema.RelationOptions{
					CollectionId: ids["tags"],
				},
			},
			&schema.SchemaField{
				Name:     "provider",
				Type:     schema.FieldTypeRelation,
				Required: false,
				Options: &schema.RelationOptions{
					CascadeDelete: true,
					MaxSelect:     types.Pointer(1),
					CollectionId:  ids["providers"],
				},
			},
		)
		if err != nil {
			slog.Error("failed to init machines collection", "err", err)
		}

		// Groups Collection
		err = initCollection(
			dao,
			"groups",
			"@request.auth.id != ''", // List Rule
			"@request.auth.id != ''", // View Rule
			"@request.auth.permission.is_admin = true || "+
				"(@request.auth.permission.access_groups = true && "+
				"@request.auth.permission.can_create = true)", // Create Rule
			"@request.auth.permission.is_admin = true || "+
				"(@request.auth.permission.can_update = true && "+
				"(@request.auth.permission.access_groups = true || "+
				"@request.auth.permission.groups.id ?= id))", // Update Rule
			"@request.auth.permission.is_admin = true || "+
				"(@request.auth.permission.can_delete = true && "+
				"(@request.auth.permission.access_groups = true || "+
				"@request.auth.permission.groups.id ?= id))", // Delete Rule
			types.JsonArray[string]{
				"CREATE UNIQUE INDEX idx_name ON groups(name)",
			},
			&schema.SchemaField{
				Name:        "name",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: true,
			},
			&schema.SchemaField{
				Name:     "description",
				Type:     schema.FieldTypeText,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "linux_username",
				Type:     schema.FieldTypeText,
				Required: true,
			},
		)
		if err != nil {
			slog.Error("failed to init groups collection", "err", err)
		}

		// Users Collection
		err = initCollection(
			dao,
			"users",
			"@request.auth.id = id || "+
				"@request.auth.permission.is_admin = true || "+
				"@request.auth.permission.access_users = true", // List Rule
			"@request.auth.id = id || "+
				"@request.auth.permission.is_admin = true || "+
				"@request.auth.permission.access_users = true", // View Rule
			"", // Create Rule
			"@request.auth.permission.is_admin = true || "+
				"@request.auth.permission.access_users = true || "+
				"(@request.auth.id = id && "+
				"@request.data.principal:isset = false && "+
				"@request.data.permission:isset = false)", // Update Rule
			"@request.auth.id = id || "+
				"@request.auth.permission.is_admin = true || "+
				"@request.auth.permission.access_users = true", // Delete Rule
			nil,
			&schema.SchemaField{
				Name:     "principal",
				Type:     schema.FieldTypeText,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "permission",
				Type:     schema.FieldTypeRelation,
				Required: false,
				Options: &schema.RelationOptions{
					CollectionId:  ids["permissions"],
					MaxSelect:     types.Pointer(1),
					CascadeDelete: false,
				},
			},
			&schema.SchemaField{
				Name:     "groups",
				Type:     schema.FieldTypeRelation,
				Required: false,
				Options: &schema.RelationOptions{
					CollectionId:  ids["groups"],
					CascadeDelete: false,
				},
			},
			&schema.SchemaField{
				Name:     "settings",
				Type:     schema.FieldTypeJson,
				Required: false,
				Options: &schema.JsonOptions{
					MaxSize: 2000000,
				},
			},
		)
		if err != nil {
			slog.Error("failed to init users collection", "err", err)
		}

		// Tags Collection
		err = initCollection(
			dao,
			"tags",
			"@request.auth.id != ''",
			"@request.auth.id != ''",
			"@request.auth.id != ''",
			"@request.auth.id != ''",
			"@request.auth.id != ''",
			types.JsonArray[string]{
				"CREATE UNIQUE INDEX idx_name ON tags(name)",
			},
			&schema.SchemaField{
				Name:        "name",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: true,
			},
			&schema.SchemaField{
				Name:     "description",
				Type:     schema.FieldTypeText,
				Required: false,
			},
		)
		if err != nil {
			slog.Error("failed to init tags collection", "err", err)
		}

		// Permissions Collection
		err = initCollection(
			dao,
			"permissions",
			"@request.auth.permission.is_admin = true || "+
				"@request.auth.permission.id ?= id", // List Rule
			"@request.auth.permission.is_admin = true || "+
				"@request.auth.permission.id ?= id", // View Rule
			"@request.auth.permission.is_admin = true", // Create Rule
			"@request.auth.permission.is_admin = true", // Update Rule
			"@request.auth.permission.is_admin = true", // Delete Rule
			types.JsonArray[string]{
				"CREATE UNIQUE INDEX idx_permissions_name ON permissions(name)",
			},
			&schema.SchemaField{
				Name:        "name",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: true,
			},
			&schema.SchemaField{
				Name:     "description",
				Type:     schema.FieldTypeText,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "can_create",
				Type:     schema.FieldTypeBool,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "can_update",
				Type:     schema.FieldTypeBool,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "can_delete",
				Type:     schema.FieldTypeBool,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "access_users",
				Type:     schema.FieldTypeBool,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "access_machines",
				Type:     schema.FieldTypeBool,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "access_groups",
				Type:     schema.FieldTypeBool,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "is_admin",
				Type:     schema.FieldTypeBool,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "machines",
				Type:     schema.FieldTypeRelation,
				Required: false,
				Options: &schema.RelationOptions{
					CollectionId:  ids["machines"],
					CascadeDelete: false,
				},
			},
			&schema.SchemaField{
				Name:     "groups",
				Type:     schema.FieldTypeRelation,
				Required: false,
				Options: &schema.RelationOptions{
					CollectionId:  ids["groups"],
					CascadeDelete: false,
				},
			},
			&schema.SchemaField{
				Name:     "users",
				Type:     schema.FieldTypeRelation,
				Required: false,
				Options: &schema.RelationOptions{
					CollectionId:  ids["users"],
					CascadeDelete: false,
				},
			},
		)
		if err != nil {
			slog.Error("failed to init permissions collection", "err", err)
		}

		// Users Collection
		err = initCollection(
			dao,
			"requests",
			"@request.auth.permission.is_admin = true || @request.auth.id = user.id", // List Rule
			"@request.auth.permission.is_admin = true || @request.auth.id = user.id", // View Rule
			"@request.auth.id != ''", // Create Rule
			"@request.auth.permission.is_admin = true || @request.auth.id = user.id", // Update Rule
			"@request.auth.permission.is_admin = true || @request.auth.id = user.id", // Delete Rule
			types.JsonArray[string]{
				"CREATE UNIQUE INDEX idx_user_machine_group ON requests(user, machine, group)",
			},
			&schema.SchemaField{
				Name:     "description",
				Type:     schema.FieldTypeText,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "user",
				Type:     schema.FieldTypeRelation,
				Required: true,
				Options: &schema.RelationOptions{
					CollectionId:  ids["users"],
					MaxSelect:     types.Pointer(1),
					CascadeDelete: false,
				},
			},
			&schema.SchemaField{
				Name:     "group",
				Type:     schema.FieldTypeRelation,
				Required: false,
				Options: &schema.RelationOptions{
					CollectionId:  ids["groups"],
					MaxSelect:     types.Pointer(1),
					CascadeDelete: false,
				},
			},
			&schema.SchemaField{
				Name:     "machine",
				Type:     schema.FieldTypeRelation,
				Required: false,
				Options: &schema.RelationOptions{
					CollectionId:  ids["machines"],
					MaxSelect:     types.Pointer(1),
					CascadeDelete: false,
				},
			},
		)
		if err != nil {
			slog.Error("failed to init requests collection", "err", err)
		}

		err = initCollection(
			dao,
			"auditlog",
			"@request.auth.permission.is_admin = true",
			"@request.auth.permission.is_admin = true",
			"@request.auth.permission.is_admin = true",
			"@request.auth.permission.is_admin = true",
			"@request.auth.permission.is_admin = true",
			nil,
			&schema.SchemaField{
				Name:     "collection",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "record",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "event",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "user",
				Type:     schema.FieldTypeRelation,
				Required: false,
				Options: &schema.RelationOptions{
					CollectionId:  ids["users"],
					MaxSelect:     types.Pointer(1),
					CascadeDelete: false,
				},
			},
			&schema.SchemaField{
				Name:     "admin",
				Type:     schema.FieldTypeText,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "data",
				Type:     schema.FieldTypeJson,
				Required: false,
				Options: &schema.JsonOptions{
					MaxSize: 2000000,
				},
			},
			&schema.SchemaField{
				Name:     "original",
				Type:     schema.FieldTypeJson,
				Required: false,
				Options: &schema.JsonOptions{
					MaxSize: 2000000,
				},
			},
		)
		if err != nil {
			slog.Error("failed to init auditlog collection", "err", err)
		}

		// Settings collection
		err = initCollection(
			dao,
			"settings",
			"@request.auth.permission.is_admin = true",
			"@request.auth.permission.is_admin = true",
			"@request.auth.permission.is_admin = true",
			"@request.auth.permission.is_admin = true",
			"@request.auth.permission.is_admin = true",
			types.JsonArray[string]{
				"CREATE UNIQUE INDEX idx_key ON settings(key)",
			},
			&schema.SchemaField{
				Name:     "key",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "value",
				Type:     schema.FieldTypeText,
				Required: false,
			},
		)
		if err != nil {
			slog.Error("failed to init settings collection", "err", err)
		}

		// Providers Collection
		err = initCollection(
			dao,
			"providers",
			"@request.auth.permission.is_admin = true",
			"@request.auth.permission.is_admin = true",
			"@request.auth.permission.is_admin = true",
			"@request.auth.permission.is_admin = true",
			"@request.auth.permission.is_admin = true",
			types.JsonArray[string]{
				"CREATE UNIQUE INDEX idx_name ON providers(name)",
			},
			&schema.SchemaField{
				Name:        "name",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: true,
			},
			&schema.SchemaField{
				Name:     "url",
				Type:     schema.FieldTypeText,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "username",
				Type:     schema.FieldTypeText,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "password",
				Type:     schema.FieldTypeText,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "token",
				Type:     schema.FieldTypeText,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "error",
				Type:     schema.FieldTypeText,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "last_sync",
				Type:     schema.FieldTypeDate,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "type",
				Type:     schema.FieldTypeText,
				Required: true,
			},
		)
		if err != nil {
			slog.Error("failed to init providers collection", "err", err)
		}

		return nil
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		machines, _ := dao.FindCollectionByNameOrId("machines")
		if machines != nil {
			dao.DeleteCollection(machines)
		}
		groups, _ := dao.FindCollectionByNameOrId("groups")
		if groups != nil {
			dao.DeleteCollection(groups)
		}
		permissions, _ := dao.FindCollectionByNameOrId("permissions")
		if permissions != nil {
			dao.DeleteCollection(permissions)
		}
		tags, _ := dao.FindCollectionByNameOrId("tags")
		if tags != nil {
			dao.DeleteCollection(tags)
		}
		settings, _ := dao.FindCollectionByNameOrId("settings")
		if settings != nil {
			dao.DeleteCollection(settings)
		}
		providers, _ := dao.FindCollectionByNameOrId("providers")
		if providers != nil {
			dao.DeleteCollection(providers)
		}
		users, _ := dao.FindCollectionByNameOrId("users")
		if users != nil {
			users.Schema.RemoveField("groups")
			users.Schema.RemoveField("principal")
			users.Schema.RemoveField("permissions")
		}

		return nil
	})
}

func initCollection(
	dao *daos.Dao,
	name string,
	listRule, viewRule, createRule, updateRule, deleteRule string,
	indexes types.JsonArray[string],
	schemaFields ...*schema.SchemaField,
) error {
	collection, err := dao.FindCollectionByNameOrId(name)
	if err != nil {
		return err
	}
	collection.ListRule = types.Pointer(listRule)
	collection.ViewRule = types.Pointer(viewRule)
	collection.CreateRule = types.Pointer(createRule)
	collection.UpdateRule = types.Pointer(updateRule)
	collection.DeleteRule = types.Pointer(deleteRule)
	collection.Indexes = types.JsonArray[string](indexes)
	for schemaField := range schemaFields {
		collection.Schema.AddField(schemaFields[schemaField])
	}

	return dao.SaveCollection(collection)
}
