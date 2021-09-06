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

func TestDeleteAccount(t *testing.T){
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

		err := authCmp.DeleteAccount(ctx, accountId)
		if data.shouldErrorOccur && err == nil {
			t.Errorf("expected error to occur but none did")
		}

		if !data.shouldErrorOccur && err != nil {
			t.Errorf("error was not expected to occur - error %s", err.Error())
		}
	}
}

func TestDeleteAccountHandler(t *testing.T) {
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

		resp, _, rr := DeleteAccountInAuthSvc(accountId, authCmp, t)
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

type DeleteAccountHandlerTestMetadata struct {
	email                        string
	password                     string
	responseCode                 int
	shouldErrorOccur             bool
	shouldCreateAccountFirst bool
}

// getDeleteAccountTestScenarios returns a set of test scenarios for the delete account test case
func getDeleteAccountTestScenarios(randomstringLength int) []DeleteAccountHandlerTestMetadata {
	email := helper.GenerateRandomString(randomstringLength)
	password := helper.GenerateRandomString(randomstringLength)

	return []DeleteAccountHandlerTestMetadata {
		// test success scenario. create an account first then successfully delete it
		{
			email,
			password,
			http.StatusOK,
			false,
			true,
		},
		// test failure scenario. delete non existent account (account was not created)
		{
			email,
			password,
			http.StatusBadRequest,
			true,
			false,
		},
	}
}

func DeleteAccountInAuthSvc(accountId uint32, cmp *AuthenticationComponent,
	t *testing.T) (DeleteAccountResponse, error,
	*httptest.ResponseRecorder) {
	var result DeleteAccountResponse

	id := fmt.Sprint(accountId)
	req, err := http.NewRequest("DELETE", "//v1/auth/account/login/" + id, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": id})


	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cmp.DeleteAccountHandler)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &result)

	return result, err, rr
}
