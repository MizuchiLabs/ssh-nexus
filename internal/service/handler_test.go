// Package service contains all the service handlers
package service

import (
	"testing"

	"github.com/MizuchiLabs/ssh-nexus/test"
	"github.com/pocketbase/pocketbase/tests"
)

func TestServer(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "Valid Server", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := Server(); (err != nil) != tt.wantErr {
				t.Errorf("Server() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeyCheck(t *testing.T) {
	type args struct {
		app *tests.TestApp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Valid host", args: args{app: test.SetupApp(t)}, wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := KeyCheck(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("KeyCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppEventHandler(t *testing.T) {
	type args struct {
		app *tests.TestApp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Valid App", args: args{app: test.SetupApp(t)}, wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := AppEventHandler(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("AppEventHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuditEventHandler(t *testing.T) {
	type args struct {
		app *tests.TestApp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Valid App", args: args{app: test.SetupApp(t)}, wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := AuditEventHandler(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("AuditEventHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserEventHandler(t *testing.T) {
	type args struct {
		app *tests.TestApp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Valid App", args: args{app: test.SetupApp(t)}, wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := UserEventHandler(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("UserEventHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMachineEventHandler(t *testing.T) {
	type args struct {
		app *tests.TestApp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Valid App", args: args{app: test.SetupApp(t)}, wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := MachineEventHandler(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("MachineEventHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
