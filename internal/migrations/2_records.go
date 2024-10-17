package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		// Create default admin permission and role
		permissions, err := dao.FindCollectionByNameOrId("permissions")
		if err != nil {
			return err
		}
		permissionsRecord := models.NewRecord(permissions)
		permissionsRecord.Set("name", "admin")
		permissionsRecord.Set("description", "Default admin permission")
		permissionsRecord.Set("is_admin", true)
		if err = dao.SaveRecord(permissionsRecord); err != nil {
			return err
		}

		// Create default group
		groups, err := dao.FindCollectionByNameOrId("groups")
		if err != nil {
			return err
		}
		groupsRecord := models.NewRecord(groups)
		groupsRecord.Set("name", "default")
		groupsRecord.Set("description", "Default group using root")
		groupsRecord.Set("linux_username", "root")
		if err = dao.SaveRecord(groupsRecord); err != nil {
			return err
		}

		return nil
	}, func(db dbx.Builder) error {
		dao := daos.New(db)
		permissions, _ := dao.FindCollectionByNameOrId("permissions")
		if permissions != nil {
			if err := dao.DeleteCollection(permissions); err != nil {
				return err
			}
		}

		return nil
	})
}
