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

func TestLogoutAccount(t *testing.T){
	randomstringLength := 15
	var scenarios = getLogoutAccountTestScenarios(randomstringLength)
	ctx := context.Background()

	for _, data := range scenarios {
		var accountId uint32 = 0

		if data.shouldCreateAndAuthenticateAccountFirst {
			id, err := authCmp.CreateAccount(ctx, data.email, data.password, false)
			if err != nil {
				t.Errorf("error should not have occured")
			}

			accountId = id

			_, err = authCmp.AuthenticateAccount(ctx, data.email, data.password)
			if err != nil {
				t.Errorf("error should not have occured")
			}
		}

		err := authCmp.LogoutAccount(ctx, accountId)
		if data.shouldErrorOccur && err == nil {
			t.Errorf("expected error to occur but none did")
		}

		if !data.shouldErrorOccur && err != nil {
			t.Errorf("error was not expected to occur - error %s", err.Error())
		}
	}
}

func TestLogoutAccountHandler(t *testing.T) {
	randomstringLength := 15
	var scenarios = getLogoutAccountTestScenarios(randomstringLength)

	for _, data := range scenarios {
		var accountId uint32 = 0

		if data.shouldCreateAndAuthenticateAccountFirst {
			createAccResp, err, _ := CreateAccountInAuthSvc(data.email, data.password, authCmp, t)
			if err != nil || createAccResp.ErrorMessage != EMPTY {
				t.Errorf("error should not have occured")
			}

			accountId = createAccResp.Id
			// authenticate account
			authAccResp, err, _ := AuthenticateAccountInAuthSvc(data.email, data.password, authCmp, t)
			if err != nil || authAccResp.ErrorMessage != EMPTY {
				t.Errorf("error should not have occured")
			}
		}

		resp, _, rr := LogoutAccountInAuthSvc(accountId, authCmp, t)
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
	}
}

type LogoutAccountHandlerTestMetadata struct {
	email                        string
	password                     string
	responseCode                 int
	shouldErrorOccur             bool
	shouldCreateAndAuthenticateAccountFirst bool
}

// getLogoutAccountTestScenarios returns a set of test scenarios for the logout account test case
func getLogoutAccountTestScenarios(randomstringLength int) []LogoutAccountHandlerTestMetadata {
	email := helper.GenerateRandomString(randomstringLength)
	password := helper.GenerateRandomString(randomstringLength)

	return []LogoutAccountHandlerTestMetadata {
		// test success scenario. create an account first then authenticate and then successfully log out
		{
			email,
			password,
			http.StatusOK,
			false,
			true,
		},
		// test failure scenario. logout non-existent account (account was not created)
		{
			email,
			password,
			http.StatusBadRequest,
			true,
			false,
		},
	}
}

func LogoutAccountInAuthSvc(accountId uint32, cmp *AuthenticationComponent,
	t *testing.T) (LogoutAccountResponse, error,
	*httptest.ResponseRecorder) {
	var result LogoutAccountResponse

	id := fmt.Sprint(accountId)
	req, err := http.NewRequest("POST", "//v1/auth/account/logout/" + id, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": id})


	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cmp.LogoutAccountHandler)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &result)

	return result, err, rr
}
