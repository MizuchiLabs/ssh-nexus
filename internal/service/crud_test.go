package service

import (
	"testing"

	"github.com/MizuchiLabs/ssh-nexus/test"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func Test_setPrincipalUUID(t *testing.T) {
	type args struct {
		app    *tests.TestApp
		record *models.Record
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Add Principal UUID",
			args: args{
				app:    test.SetupApp(t),
				record: test.GetRecord(t, "users", "principal = ''"),
			},
			wantErr: false,
		},
		{
			name: "Already has Principal UUID",
			args: args{
				app:    test.SetupApp(t),
				record: test.GetRecord(t, "users", "principal != ''"),
			},
			wantErr: false,
		},
		{
			name: "Wrong record type",
			args: args{
				app:    test.SetupApp(t),
				record: test.GetRecord(t, "machines", "id != ''"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := setPrincipalUUID(tt.args.app, tt.args.record); (err != nil) != tt.wantErr {
				t.Errorf("setPrincipalUUID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cleanupTags(t *testing.T) {
	type args struct {
		app *tests.TestApp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Cleanup Tags",
			args: args{
				app: test.SetupApp(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := cleanupTags(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("cleanupTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_insertAuditlog(t *testing.T) {
	type args struct {
		app         *tests.TestApp
		httpContext echo.Context
		record      *models.Record
		event       string
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
			if err := insertAuditlog(tt.args.app, tt.args.httpContext, tt.args.record, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("insertAuditlog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
