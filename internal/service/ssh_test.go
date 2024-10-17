package service

import (
	"reflect"
	"testing"

	"github.com/MizuchiLabs/ssh-nexus/test"
	"github.com/pocketbase/pocketbase/models"
	"golang.org/x/crypto/ssh"
)

// func TestManualUpdate(t *testing.T) {
// 	type args struct {
// 		app     core.App
// 		machine *models.Record
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		{
// 			name: "Manual Update",
// 			args: args{
// 				app:     test.SetupApp(t),
// 				machine: test.GetRecord(t, "machines", "agent = false"),
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			ManualUpdate(tt.args.app, tt.args.machine)
// 		})
// 	}
// }
//
// func TestInstallAgent(t *testing.T) {
// 	type args struct {
// 		app     core.App
// 		machine *models.Record
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			InstallAgent(tt.args.app, tt.args.machine)
// 		})
// 	}
// }
//
// func TestSyncAgentToken(t *testing.T) {
// 	type args struct {
// 		app core.App
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			SyncAgentToken(tt.args.app)
// 		})
// 	}
// }
//
// func TestRestore(t *testing.T) {
// 	type args struct {
// 		machine *models.Record
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			Restore(tt.args.machine)
// 		})
// 	}
// }

func Test_uploadAgent(t *testing.T) {
	type args struct {
		conn *ssh.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Not connected",
			args: args{
				conn: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := uploadAgent(tt.args.conn); (err != nil) != tt.wantErr {
				t.Errorf("uploadAgent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setPrincipals(t *testing.T) {
	type args struct {
		conn   *ssh.Client
		groups map[string][]string
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
			if err := setPrincipals(tt.args.conn, tt.args.groups); (err != nil) != tt.wantErr {
				t.Errorf("setPrincipals() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_connect(t *testing.T) {
	type args struct {
		machine *models.Record
	}
	tests := []struct {
		name    string
		args    args
		want    *ssh.Client
		wantErr bool
	}{
		{
			name: "Not connected",
			args: args{
				machine: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Connected",
			args: args{
				machine: test.GetRecord(t, "machines", "agent = false"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := connect(tt.args.machine)
			if (err != nil) != tt.wantErr {
				t.Errorf("connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("connect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_run(t *testing.T) {
	type args struct {
		client  *ssh.Client
		command string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := run(tt.args.client, tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("run() = %v, want %v", got, tt.want)
			}
		})
	}
}
