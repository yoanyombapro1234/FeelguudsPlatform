package authentication_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
)

func TestCreateAccount(t *testing.T) {
	randomstringLength := 15
	var scenarios = getCreateAccountTestScenarios(randomstringLength)
	ctx := context.Background()

	for _, data := range scenarios {
		id, err := authCmp.CreateAccount(ctx, data.email, data.password, false)
		if data.shouldErrorOccur && err == nil {
			t.Errorf("expected error to occur but none did")
		}

		if !data.shouldErrorOccur && err != nil {
			t.Errorf("error was not expected to occur - error %s", err.Error())
		}

		if id == 0 && !data.shouldErrorOccur {
			t.Errorf("null or invalid user acocunt id not expected")
		} else if id != 0 && data.shouldErrorOccur {
			t.Errorf("expected invalid user account id but got a valid instead")
		}
	}
}

func TestCreateAccountHandler(t *testing.T) {
	randomstringLength := 15
	var scenarios = getCreateAccountTestScenarios(randomstringLength)

	for _, data := range scenarios {
		resp, err, rr := CreateAccountInAuthSvc(data.email, data.password, authCmp, t)
		if data.shouldErrorOccur && (resp.ErrorMessage == EMPTY) {
			t.Errorf("expected error to occur but none did")
		}

		if !data.shouldErrorOccur && (resp.ErrorMessage != EMPTY) {
			t.Errorf("error was not expected to occur - error %s", err.Error())
		}

		// Check the status code is what we expect.
		if status := rr.Code; status != data.responseCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, data.responseCode)
		}

		if resp.Id == 0 && !data.shouldErrorOccur {
			t.Errorf("null or invalid user acocunt id not expected")
		} else if resp.Id != 0 && data.shouldErrorOccur {
			t.Errorf("expected invalid user account id but got a valid instead")
		}
	}
}

type CreateAccountHandlerTestMetadata struct {
	email            string
	password         string
	responseCode     int
	shouldErrorOccur bool
}

// getCreateAccountTestScenarios returns a set of test scenarios for the create account test case
func getCreateAccountTestScenarios(randomstringLength int) []CreateAccountHandlerTestMetadata {
	email := helper.GenerateRandomString(randomstringLength)
	password := helper.GenerateRandomString(randomstringLength)

	return []CreateAccountHandlerTestMetadata{
		// test success scenario. create an account from scratch
		{
			email,
			password,
			http.StatusOK,
			false,
		},
		// test failure scenario. create duplicate account
		{
			email,
			password,
			http.StatusBadRequest,
			true,
		},
		// test failure scenario. empty input parameters
		{
			EMPTY,
			EMPTY,
			http.StatusBadRequest,
			true,
		},
		// test failure scenario. empty password parameter
		{
			helper.GenerateRandomString(randomstringLength),
			EMPTY,
			http.StatusBadRequest,
			true,
		},
		// test failure scenario. empty email parameter
		{
			EMPTY,
			helper.GenerateRandomString(randomstringLength),
			http.StatusBadRequest,
			true,
		},
	}
}

func CreateAccountInAuthSvc(email, password string, cmp *AuthenticationComponent,
	t *testing.T) (CreateAccountResponse, error,
	*httptest.ResponseRecorder) {
	var result CreateAccountResponse

	reqBody := CreateAccountRequest{
		Email:    email,
		Password: password,
	}

	body, err := helper.CreateRequestBody(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "//v1/auth/account/create", body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cmp.CreateAccountHandler)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &result)

	return result, err, rr
}
