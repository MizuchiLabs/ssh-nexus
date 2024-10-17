package test

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tokens"
)

const testDataDir = "./test_pb_data"

func SetupApp(t *testing.T) *tests.TestApp {
	testApp, err := tests.NewTestApp(testDataDir)
	if err != nil {
		t.Fatal(err)
	}
	return testApp
}

func GenerateAdminToken(t *testing.T, email string) string {
	app, err := tests.NewTestApp(testDataDir)
	if err != nil {
		t.Fatal(err)
	}
	defer app.Cleanup()

	admin, err := app.Dao().FindAdminByEmail(email)
	if err != nil {
		t.Fatal(err)
	}

	token, err := tokens.NewAdminAuthToken(app, admin)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func GenerateRecordToken(t *testing.T, collectionNameOrID string, email string) string {
	app, err := tests.NewTestApp(testDataDir)
	if err != nil {
		t.Fatal(err)
	}
	defer app.Cleanup()

	record, err := app.Dao().FindAuthRecordByEmail(collectionNameOrID, email)
	if err != nil {
		t.Fatal(err)
	}

	token, err := tokens.NewRecordAuthToken(app, record)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func GetRecord(t *testing.T, collectionNameOrID, filter string) *models.Record {
	app, err := tests.NewTestApp(testDataDir)
	if err != nil {
		t.Fatal(err)
	}
	defer app.Cleanup()

	record, err := app.Dao().FindFirstRecordByFilter(collectionNameOrID, filter, nil)
	if err != nil {
		t.Fatal(err)
	}
	return record
}

func NewRequest(method, url string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(method, url, nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}
