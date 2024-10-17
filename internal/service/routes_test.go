package service

import (
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
)

func Test_initRoutes(t *testing.T) {
	type args struct {
		app core.App
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			initRoutes(tt.args.app)
		})
	}
}

func Test_getVersion(t *testing.T) {
	// TODO: Add test cases
	// scenarios := []tests.ApiScenario{
	// 	{
	// 		Name:           "try with GET method",
	// 		Method:         http.MethodGet,
	// 		Url:            "/version",
	// 		ExpectedStatus: 200,
	// 		TestAppFactory: test.SetupApp,
	// 	},
	// 	{
	// 		Name:           "try with POST method",
	// 		Method:         http.MethodPost,
	// 		Url:            "/version",
	// 		ExpectedStatus: 200,
	// 		TestAppFactory: test.SetupApp,
	// 	},
	// }
	//
	// for _, scenario := range scenarios {
	// 	scenario.Test(t)
	// }
}

func Test_getUserMachines(t *testing.T) {
	type args struct {
		c   echo.Context
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
			if err := getUserMachines(tt.args.c, tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("getUserMachines() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getPublicKey(t *testing.T) {
	type args struct {
		c func() ([]byte, error)
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
			if err := getPublicKey(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("getPublicUserKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setUserCA(t *testing.T) {
	type args struct {
		c echo.Context
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
			if err := setUserCA(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("setUserCA() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_rotateSSHKeys(t *testing.T) {
	type args struct {
		c echo.Context
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
			if err := rotateSSHKeys(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("rotateSSHKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getServerCertificate(t *testing.T) {
	type args struct {
		c echo.Context
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
			if err := getServerCertificate(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("getServerCertificate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getAgentToken(t *testing.T) {
	type args struct {
		c echo.Context
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
			if err := getAgentToken(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("getAgentToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_rotateAgentToken(t *testing.T) {
	type args struct {
		c echo.Context
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
			if err := rotateAgentToken(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("rotateAgentToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_signUserCertificate(t *testing.T) {
	type args struct {
		c   echo.Context
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
			if err := signUserCertificate(tt.args.c, tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("signUserCertificate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_signHostCertificate(t *testing.T) {
	type args struct {
		c   echo.Context
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
			if err := signHostCertificate(tt.args.c, tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("signHostCertificate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
