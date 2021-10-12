package authentication_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
)

func TestGetAccount(t *testing.T) {
	randomstringLength := 15
	var scenarios = getDeleteAccountTestScenarios(randomstringLength)
	ctx := context.Background()

	for _, data := range scenarios {
		var accountId uint32 = 0

		if data.shouldCreateAccountFirst {
			createAccResp, err, _ := CreateAccountInAuthSvc(data.email, data.password, authCmp, t)
			if err != nil || createAccResp.ErrorMessage != EMPTY {
				t.Errorf("error should not have occured")
			}

			accountId = createAccResp.Id
		}

		account, err := authCmp.GetAccount(ctx, accountId)
		if data.shouldErrorOccur && err == nil {
			t.Errorf("expected error to occur but none did")
		}

		if !data.shouldErrorOccur && err != nil {
			t.Errorf("error was not expected to occur - error %s", err.Error())
		}

		if account == nil && !data.shouldErrorOccur {
			t.Errorf("error was not expected to occur - error %s", err.Error())
		}
	}
}

func TestGetAccountHandler(t *testing.T) {
	randomstringLength := 15
	var scenarios = getDeleteAccountTestScenarios(randomstringLength)

	for _, data := range scenarios {
		var accountId uint32 = 0

		if data.shouldCreateAccountFirst {
			createAccResp, err, _ := CreateAccountInAuthSvc(data.email, data.password, authCmp, t)
			if err != nil || createAccResp.ErrorMessage != EMPTY {
				t.Errorf("error should not have occured")
			}

			accountId = createAccResp.Id
		}

		resp, _, rr := GetAccountInAuthSvc(accountId, authCmp, t)
		if data.shouldErrorOccur && (resp.ErrorMessage == EMPTY) {
			t.Errorf("expected error to occur but none did")
		}

		if !data.shouldErrorOccur && (resp.ErrorMessage != EMPTY) {
			t.Errorf("error was not expected to occur - error %s", resp.ErrorMessage)
		}

		// Check the status code is what we expect.
		if status := rr.Code; status != data.responseCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, data.responseCode)
		}

		if resp.Account == nil && !data.shouldErrorOccur {
			t.Errorf("error was not expected to occur")
		}
	}
}

type GetAccountHandlerTestMetadata struct {
	email                    string
	password                 string
	responseCode             int
	shouldErrorOccur         bool
	shouldCreateAccountFirst bool
}

// getGetAccountTestScenarios returns a set of test scenarios for the get account test case
func getGetAccountTestScenarios(randomstringLength int) []GetAccountHandlerTestMetadata {
	email := helper.GenerateRandomString(randomstringLength)
	password := helper.GenerateRandomString(randomstringLength)

	return []GetAccountHandlerTestMetadata{
		// test success scenario. create an account first then successfully get it
		{
			email,
			password,
			http.StatusOK,
			false,
			true,
		},
		// test failure scenario. get non existent account (account was not created)
		{
			email,
			password,
			http.StatusBadRequest,
			true,
			false,
		},
	}
}

func GetAccountInAuthSvc(accountId uint32, cmp *AuthenticationComponent,
	t *testing.T) (GetAccountResponse, error,
	*httptest.ResponseRecorder) {
	var result GetAccountResponse

	id := fmt.Sprint(accountId)
	req, err := http.NewRequest("GET", "//v1/auth/account/"+id, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": id})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cmp.GetAccountHandler)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &result)

	return result, err, rr
}
