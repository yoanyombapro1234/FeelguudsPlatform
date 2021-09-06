package authentication_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
)

var authCmp = InitializeMockAuthenticationComponent()

func TestAuthenticateAccount(t *testing.T){
	randomstringLength := 15
	var scenarios  = getAuthenticationTestScenarios(randomstringLength)

	for _, data := range scenarios {
		ctx := context.Background()
		if data.shouldCreateUserAccountFirst {
			_, err := authCmp.CreateAccount(ctx, data.email, data.password, false)
			if err != nil && !data.shouldErrorOccur {
				t.Errorf("obtained error but not expected - %s", err.Error())
			}
		}

		token, err := authCmp.AuthenticateAccount(ctx, data.email, data.password)
		if data.shouldErrorOccur && err == nil {
			t.Errorf("expected error to occur but none did")
		}

		if !data.shouldErrorOccur && err != nil {
			t.Errorf("error was not expected to occur - error %s", err.Error())
		}

		if token == "" && !data.shouldErrorOccur {
			t.Errorf("null or empty jwt token not expected")
		} else if token != "" && data.shouldErrorOccur {
			t.Errorf("expected empty json web token but got a token instead")
		}
	}
}

func TestAuthenticateAccountHandler(t *testing.T) {
	randomstringLength := 15
	var scenarios  = getAuthenticationTestScenarios(randomstringLength)

	for _, data := range scenarios {
		if data.shouldCreateUserAccountFirst {
			resp, err, _ := CreateAccountInAuthSvc(data.email, data.password, authCmp ,t)
			if (err != nil || resp.ErrorMessage != EMPTY) && !data.shouldErrorOccur {
				t.Errorf("obtained error but not expected - %s", err.Error())
			}
		}

		resp, err, rr := AuthenticateAccountInAuthSvc(data.email, data.password, authCmp, t)
		if data.shouldErrorOccur && resp.ErrorMessage == EMPTY {
			t.Errorf("expected error to occur but none did")
		}

		if !data.shouldErrorOccur && resp.ErrorMessage != EMPTY {
			t.Errorf("error was not expected to occur - error %s", err.Error())
		}

		// Check the status code is what we expect.
		if status := rr.Code; status != data.responseCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, data.responseCode)
		}

		if resp.Token == "" && !data.shouldErrorOccur {
			t.Errorf("null or empty jwt token not expected")
		} else if resp.Token != "" && data.shouldErrorOccur {
			t.Errorf("expected empty json web token but got a token instead")
		}
	}
}

type AuthenticationHandlerTestMetadata struct {
	email                        string
	password                     string
	responseCode                 int
	shouldErrorOccur             bool
	shouldCreateUserAccountFirst bool
}

// getAuthenticationTestScenarios returns a set of test scenarios for the authentication test case
func getAuthenticationTestScenarios(randomstringLength int) []AuthenticationHandlerTestMetadata {
	return []AuthenticationHandlerTestMetadata{
		// test success scenario. authenticating a valid account
		{
			email:                        helper.GenerateRandomString(randomstringLength),
			password:                     helper.GenerateRandomString(randomstringLength),
			responseCode:                 http.StatusOK,
			shouldErrorOccur:             false,
			shouldCreateUserAccountFirst: true,
		},
		// test failure scenario. authenticating an account that doesn't exist
		{
			email:                        helper.GenerateRandomString(randomstringLength),
			password:                     helper.GenerateRandomString(randomstringLength),
			responseCode:                 http.StatusBadRequest,
			shouldErrorOccur:             true,
			shouldCreateUserAccountFirst: false,
		},
		// test failure scenario. attempting authentication with invalid parameters (empty email and password)
		{
			email:                        EMPTY,
			password:                     EMPTY,
			responseCode:                 http.StatusBadRequest,
			shouldErrorOccur:             true,
			shouldCreateUserAccountFirst: false,
		},
		// test failure scenario. attempting authentication with invalid parameters (empty email)
		{
			email:                        helper.GenerateRandomString(randomstringLength),
			password:                     EMPTY,
			responseCode:                 http.StatusBadRequest,
			shouldErrorOccur:             true,
			shouldCreateUserAccountFirst: false,
		},
		// test failure scenario. attempting authentication with invalid parameters (empty password)
		{
			email:                        EMPTY,
			password:                     helper.GenerateRandomString(randomstringLength),
			responseCode:                 http.StatusBadRequest,
			shouldErrorOccur:             true,
			shouldCreateUserAccountFirst: false,
		},
	}
}

func AuthenticateAccountInAuthSvc(email, password string, cmp *AuthenticationComponent,
	t *testing.T) (AuthenticateAccountResponse, error,
	*httptest.ResponseRecorder) {
	var result AuthenticateAccountResponse

	reqBody := AuthenticateAccountRequest{
		Email:    email,
		Password: password,
	}

	body, err := helper.CreateRequestBody(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "//v1/auth/account/login", body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cmp.AuthenticateAccountHandler)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &result)

	return result, err, rr
}
