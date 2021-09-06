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

func TestUpdateAccount(t *testing.T){
	randomstringLength := 15
	var scenarios = getUpdateAccountTestScenarios(randomstringLength)
	ctx := context.Background()

	for _, data := range scenarios {
		var accountId uint32 = 0
		if data.shouldCreateAccount {
			password := helper.GenerateRandomString(15)
			id, err := authCmp.CreateAccount(ctx, data.email, password, false)
			if err != nil {
				t.Errorf("expected error to occur but none did")
			}

			accountId = id
		}

		err := authCmp.UpdateAccount(ctx, accountId, data.updatedEmail)
		if data.shouldErrorOccur && err == nil {
			t.Errorf("expected error to occur but none did")
		}

		if !data.shouldErrorOccur && err != nil {
			t.Errorf("error was not expected to occur - error %s", err.Error())
		}
	}
}

func TestUpdateAccountHandler(t *testing.T) {
	randomstringLength := 15
	var scenarios = getUpdateAccountTestScenarios(randomstringLength)

	for _, data := range scenarios {
		var accountId uint32 = 0
		if data.shouldCreateAccount {
			password := helper.GenerateRandomString(15)
			resp, err, _ := CreateAccountInAuthSvc(data.email, password, authCmp, t)
			if resp.ErrorMessage != EMPTY || err != nil {
				t.Errorf("expected error to occur but none did")
			}

			accountId = resp.Id
		}

		resp, _, rr := UpdateEmailInAuthSvc(accountId, data.updatedEmail, authCmp, t)
		if data.shouldErrorOccur && (resp.ErrorMessage == EMPTY) {
			t.Errorf("expected error to occur but none did")
		}

		if !data.shouldErrorOccur && (resp.ErrorMessage != EMPTY) {
			t.Errorf("error was not expected to occur")
		}

		// Check the status code is what we expect.
		if status := rr.Code; status != data.responseCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, data.responseCode)
		}
	}
}

type UpdateAccountHandlerTestMetadata struct {
	email                        string
	updatedEmail string
	responseCode                 int
	shouldErrorOccur             bool
	shouldCreateAccount bool
}

// getUpdateAccountTestScenarios returns a set of test scenarios for the update account test case
func getUpdateAccountTestScenarios(randomstringLength int) []UpdateAccountHandlerTestMetadata {
	return []UpdateAccountHandlerTestMetadata {
		// test success scenario. update an existing account
		{
			helper.GenerateRandomString(randomstringLength)+"@gmail.com",
			helper.GenerateRandomString(randomstringLength)+"@gmail.com",
			http.StatusOK,
			false,
			true,
		},
		// test failure scenario. empty input parameters
		{
			helper.GenerateRandomString(randomstringLength)+"@gmail.com",
			EMPTY,
			http.StatusBadRequest,
			true,
			true,
		},
	}
}

func UpdateEmailInAuthSvc(accountId uint32, email string, cmp *AuthenticationComponent,
	t *testing.T) (UpdateEmailResponse, error,
	*httptest.ResponseRecorder) {
	var result UpdateEmailResponse

	reqBody := UpdateEmailRequest{
		Email:    email,
	}

	body, err := helper.CreateRequestBody(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "//v1/auth/account/update", body)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(accountId)})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cmp.UpdateEmailHandler)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &result)

	return result, err, rr
}
