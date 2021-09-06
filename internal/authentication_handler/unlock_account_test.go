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

func TestUnLockAccount(t *testing.T){
	randomstringLength := 15
	var scenarios = getLockAccountTestScenarios(randomstringLength)
	ctx := context.Background()

	for _, data := range scenarios {
		var accountId uint32 = 0

		if data.shouldCreateAccountFirst {
			var initialAccountLockStatus = data.doubleLockScenario
			id, err := authCmp.CreateAccount(ctx, data.email, data.password, initialAccountLockStatus)
			if err != nil {
				t.Errorf("error should not have occured")
			}

			accountId = id
		}

		err := authCmp.UnLockAccount(ctx, accountId)
		if data.shouldErrorOccur && err == nil {
			t.Errorf("expected error to occur but none did")
		}

		if !data.shouldErrorOccur && err != nil {
			t.Errorf("error was not expected to occur - error %s", err.Error())
		}

		// get the account and ensure it is unlocked if it should be
		if !data.shouldErrorOccur {
			account, err := authCmp.GetAccount(ctx, accountId)
			if err != nil {
				t.Errorf("error should not have occured - error %s", err.Error())
			}

			if account.Locked {
				t.Errorf("account is not locked and should be locked")
			}
		}
	}
}

func TestUnLockAccountHandler(t *testing.T) {
	randomstringLength := 15
	var scenarios = getUnLockAccountTestScenarios(randomstringLength)

	for _, data := range scenarios {
		var accountId uint32 = 0

		if data.shouldCreateAccountFirst {
			// account is by default created unlocked
			createAccResp, err, _ := CreateAccountInAuthSvc(data.email, data.password, authCmp, t)
			if err != nil || createAccResp.ErrorMessage != EMPTY {
				t.Errorf("error should not have occured")
			}

			accountId = createAccResp.Id
		}

		resp, _, rr := UnLockAccountInAuthSvc(accountId, authCmp, t)
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

		// ensure the unlocked account is actually indeed unlocked
		if !data.shouldErrorOccur {
			getAccResp, err, _ := GetAccountInAuthSvc(accountId, authCmp, t)
			if err != nil || getAccResp.ErrorMessage != EMPTY {
				t.Errorf("error should not have occured")
			}

			if getAccResp.Account.Locked {
				t.Errorf("account should be unlocked")
			}
		}
	}
}

type UnLockAccountHandlerTestMetadata struct {
	email                        string
	password                     string
	responseCode                 int
	shouldErrorOccur             bool
	shouldCreateAccountFirst     bool
	doubleUnLockScenario          bool
}

// getUnLockAccountTestScenarios returns a set of test scenarios for the unlock account test case
func getUnLockAccountTestScenarios(randomstringLength int) []UnLockAccountHandlerTestMetadata {
	email := helper.GenerateRandomString(randomstringLength)
	password := helper.GenerateRandomString(randomstringLength)

	return []UnLockAccountHandlerTestMetadata {
		// test success scenario. create an account first then successfully get it
		{
			email,
			password,
			http.StatusOK,
			false,
			true,
			false,
		},
		// test failure scenario. get non existent account (account was not created)
		{
			email,
			password,
			http.StatusBadRequest,
			true,
			false,
			false,
		},
		// testing double lock scenario (lock an already locked account). error should not arise in this case
		{
			helper.GenerateRandomString(randomstringLength),
			helper.GenerateRandomString(randomstringLength),
			http.StatusOK,
			false,
			true,
			true,
		},
	}
}

func UnLockAccountInAuthSvc(accountId uint32, cmp *AuthenticationComponent,
	t *testing.T) (UnLockAccountResponse, error,
	*httptest.ResponseRecorder) {
	var result UnLockAccountResponse

	id := fmt.Sprint(accountId)
	req, err := http.NewRequest("POST", "//v1/auth/account/unlock/" + id, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": id})


	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cmp.UnLockAccountHandler)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &result)

	return result, err, rr
}
