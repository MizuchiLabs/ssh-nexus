package server

import (
	"encoding/json"
	"fmt"
	"slices"

	agentv1 "github.com/MizuchiLabs/ssh-nexus/api/proto/gen/agent/v1"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func filterMachine(
	app core.App,
	host, agentID, hostname string,
) (*models.Record, error) {
	machine, err := app.Dao().
		FindFirstRecordByFilter("machines", "uuid = {:uuid}", dbx.Params{"uuid": agentID})
	if err != nil {
		machine, _ = app.Dao().
			FindFirstRecordByFilter("machines", "host = {:host}", dbx.Params{"host": host})
	}

	// Create machine if it doesn't exist
	if machine == nil {
		return addMachine(app, host, agentID, hostname)
	}
	// If uuid is not yet populated set it
	if machine.GetString("uuid") == "" {
		machine.Set("uuid", agentID)
		if err := app.Dao().SaveRecord(machine); err != nil {
			return nil, fmt.Errorf("failed to save machine: %v", err)
		}
	} else if machine.GetString("uuid") != agentID {
		// If uuid is already populated, don't let the agent connect
		machine.Set("error", "agent with uuid "+agentID+" tried to connect to machine with uuid "+machine.GetString("uuid"))
		if err := app.Dao().SaveRecord(machine); err != nil {
			return nil, fmt.Errorf("failed to save machine: %v", err)
		}
		return nil, fmt.Errorf("machine already exists with uuid %s", agentID)
	}
	if !machine.GetBool("agent") {
		machine.Set("agent", true)
		if err := app.Dao().SaveRecord(machine); err != nil {
			return nil, fmt.Errorf("failed to save machine: %v", err)
		}
	}
	if machine.GetString("error") != "" {
		machine.Set("error", "")
		if err := app.Dao().SaveRecord(machine); err != nil {
			return nil, fmt.Errorf("failed to save machine: %v", err)
		}
	}

	return machine, nil
}

// Create a new machine record
func addMachine(
	app core.App,
	host, agentID, hostname string,
) (*models.Record, error) {
	machines, err := app.Dao().FindCollectionByNameOrId("machines")
	if err != nil {
		return nil, err
	}
	machine := models.NewRecord(machines)
	machine.Set("name", hostname)
	machine.Set("host", host)
	machine.Set("uuid", agentID)
	machine.Set("port", 22) // Assuming default port
	machine.Set("agent", true)
	if err := app.Dao().SaveRecord(machine); err != nil {
		return nil, fmt.Errorf("failed to save machine: %v", err)
	}

	return machine, nil
}

func getPrincipals(
	app core.App,
	machine *models.Record,
) ([]*agentv1.StreamResponse_Principal, error) {
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

	// Serialize data
	var principals []*agentv1.StreamResponse_Principal
	for key, values := range data {
		principal := &agentv1.StreamResponse_Principal{
			Key:    key,
			Values: values,
		}
		principals = append(principals, principal)
	}

	return principals, nil
}

func getUserMachines(
	app core.App,
	user *models.Record,
) ([]*models.Record, error) {
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
