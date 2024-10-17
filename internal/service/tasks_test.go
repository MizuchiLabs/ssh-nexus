package service

import (
	"reflect"
	"testing"

	"github.com/MizuchiLabs/ssh-nexus/test"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func TestGetMachineUsers(t *testing.T) {
	type args struct {
		app     core.App
		machine *models.Record
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Valid Machine Users",
			args: args{
				app:     test.SetupApp(t),
				machine: test.GetRecord(t, "machines", "groups != '' && users != ''"),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Invalid Machine",
			args: args{
				app:     test.SetupApp(t),
				machine: nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetMachineUsers(tt.args.app, tt.args.machine)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMachineUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 && tt.want {
				t.Errorf("GetMachineUsers() = %v, want %v", got, tt.wantErr)
			}
			if len(got) > 0 && !tt.want {
				t.Errorf("GetMachineUsers() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestGetUserMachines(t *testing.T) {
	type args struct {
		app  core.App
		user *models.Record
	}
	tests := []struct {
		name    string
		args    args
		want    []*models.Record
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetUserMachines(tt.args.app, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserMachines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserMachines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_syncMachines(t *testing.T) {
	type args struct {
		app core.App
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := syncMachines(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("syncMachines() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_syncAgents(t *testing.T) {
	type args struct {
		app core.App
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := syncAgents(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("syncAgents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_syncAgentToken(t *testing.T) {
	type args struct {
		app core.App
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := syncAgentToken(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("syncAgentToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_syncProviders(t *testing.T) {
	type args struct {
		app core.App
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := syncProviders(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("syncProviders() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cleanupAudit(t *testing.T) {
	type args struct {
		app core.App
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Valid Cleanup", args: args{app: test.SetupApp(t)}, wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cleanupAudit(tt.args.app)
		})
	}
}
