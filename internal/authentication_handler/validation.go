package authentication_handler

import (
	"fmt"

	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-auth-sdk"
	"go.uber.org/zap"
)

// IsPasswordOrEmailInValid checks request parameters for validity
func (c *AuthenticationComponent) IsPasswordOrEmailInValid(email string, password string) (bool, error) {
	err, ok := c.isValidEmail(email)
	if !ok {
		return true, err
	}

	err, ok = c.isValidPassword(password)
	if !ok {
		return true, err
	}

	return false, nil
}

// isValidEmail checks if an email is valid.
func (c *AuthenticationComponent) isValidEmail(email string) (error, bool) {
	if email == "" {
		err := ErrInvalidInputArguments
		c.Logger.Error("invalid input parameters. please specify a valid email", zap.Error(err))
		return err, false
	}

	return nil, true
}

// IsValidPassword checks if a password is valid;
func (c *AuthenticationComponent) isValidPassword(password string) (error, bool) {
	if password == "" {
		err := ErrInvalidInputArguments
		c.Logger.Error("invalid input parameters. please specify a valid password", zap.Error(err))
		return err, false
	}

	return nil, true
}

// checkJwtTokenForInValidity checks jwt token and asserts the token is a valid one.
func (c *AuthenticationComponent) checkJwtTokenForInValidity(result interface{}) (bool, error) {
	token := fmt.Sprintf("%v", result)
	if token == "" {
		err := ErrJWTCastingError
		c.Logger.Error("casting error", zap.Error(err))
		return true, err
	}

	return false, nil
}

// getIdFromResponseObject attempts to cast a generic response to a type int and returns the proper value if no errors occurred.
func (c *AuthenticationComponent) getIdFromResponseObject(response interface{}) (int, error) {
	// TODO: change to to int64 in order to limit overflow from happening if tons of customer accounts are created
	id, ok := response.(int)
	if !ok {
		err := ErrTypeConversionError
		c.Logger.Error("casting error", zap.Error(err))
		return 0, err
	}
	return id, nil
}

// isValidID checks that the ID passed as part of the request parameters is indeed valid.
func (c *AuthenticationComponent) isValidID(Id uint32) (error, bool) {
	if Id == 0 {
		err := ErrInvalidInputArguments
		c.Logger.Error("invalid input parameters. please specify a valid user id", zap.Error(err))
		return err, false
	}

	return nil, true
}

// getAccountFromResponseObject obtains an account object from the response object; this account is obtained via an attempted casting operation
func (c *AuthenticationComponent) getAccountFromResponseObject(result interface{}) (*core_auth_sdk.Account, error) {
	account, ok := result.(*core_auth_sdk.Account)
	if !ok {
		err := ErrFailedToCastAccount
		c.Logger.Error(err.Error())
		return nil, err
	}
	return account, nil
}
