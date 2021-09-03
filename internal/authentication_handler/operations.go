package authentication_handler

import (
	"context"
	"fmt"
	"strconv"

	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-auth-sdk"
	"go.uber.org/zap"
)

// AuthenticateAccount attempts to authenticate a user based on provided credentials
func (c *AuthenticationComponent) AuthenticateAccount(ctx context.Context, email, password string) (string, error) {
	if _, err := c.IsPasswordOrEmailInValid(email, password); err != nil {
		return "", err
	}

	token, err := c.Client.LoginAccount(email, password)
	if err != nil {
		c.Logger.Error(err.Error())
		return "", err
	}

	if _, err := c.checkJwtTokenForInValidity(token); err != nil {
		c.Logger.Error(err.Error())
	}

	c.Logger.Info("successfully authenticated user account")
	return token, nil
}

// CreateAccount attempts to create a user account against the authentication service
func (c *AuthenticationComponent) CreateAccount(ctx context.Context, email, password string, accountLocked bool) (uint32, error) {
	if _, err := c.IsPasswordOrEmailInValid(email, password); err != nil {
		return 0, err
	}

	accountId, err := c.Client.ImportAccount(email, password, accountLocked)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("failed to create account against authentication service for user %s. error: %s", email, err.Error()))
		return 0, err
	}

	c.Logger.Info("Successfully created user account", zap.Int("Id", int(accountId)))
	return uint32(accountId), nil
}

// DeleteAccount attempts to archive an account from the context of the authentication service (authn)
func (c *AuthenticationComponent) DeleteAccount(ctx context.Context, Id uint32) error {
	if err, _ := c.isValidID(Id); err != nil {
		c.Logger.Error(err.Error())
		return err
	}

	accountId := strconv.Itoa(int(Id))
	if err := c.Client.ArchiveAccount(accountId); err != nil {
		c.Logger.Error(err.Error())
		return err
	}

	c.Logger.Info("Successfully deleted user account", zap.Int("Id", int(Id)))
	return nil
}

// GetAccount obtains a user account from the context of the authentications service (authn) based on a provided user id
func (c *AuthenticationComponent) GetAccount(ctx context.Context, Id uint32) (*core_auth_sdk.Account, error) {
	if err, _ := c.isValidID(Id); err != nil {
		c.Logger.Error(err.Error())
		return nil, err
	}

	accountId := strconv.Itoa(int(Id))
	account, err := c.Client.GetAccount(accountId)
	if err != nil {
		c.Logger.Error(err.Error())
		return nil, err
	}

	c.Logger.Info("Successfully obtained user account", zap.Int("Id", int(Id)))
	return account, nil
}

// LockAccount locks a user account
func (c *AuthenticationComponent) LockAccount(ctx context.Context, Id uint32) error {
	if err, _ := c.isValidID(Id); err != nil {
		c.Logger.Error(err.Error())
		return err
	}

	accountId := strconv.Itoa(int(Id))
	if err := c.Client.LockAccount(accountId); err != nil {
		return err
	}

	c.Logger.Info("Successfully locked user account", zap.Int("Id", int(Id)))
	return nil
}

// UnLockAccount unlocks a user account
func (c *AuthenticationComponent) UnLockAccount(ctx context.Context, Id uint32) error {
	if err, _ := c.isValidID(Id); err != nil {
		c.Logger.Error(err.Error())
		return err
	}

	accountId := strconv.Itoa(int(Id))
	if err := c.Client.UnlockAccount(accountId); err != nil {
		return err
	}

	c.Logger.Info("Successfully unlocked user account", zap.Int("Id", int(Id)))
	return nil
}

// UpdateAccount updates a user account's credentials
func (c *AuthenticationComponent) UpdateAccount(ctx context.Context, Id uint32, email string) error {
	if err, _ := c.isValidID(Id); err != nil {
		c.Logger.Error(err.Error())
		return nil
	}

	if err, _ := c.isValidEmail(email); err != nil {
		c.Logger.Error(err.Error())
		return nil
	}

	accountId := strconv.Itoa(int(Id))
	if err := c.Client.Update(accountId, email); err != nil {
		c.Logger.Error(err.Error())
		return err
	}

	c.Logger.Info("Successfully updated user account", zap.Int("Id", int(Id)))
	return nil
}

func (c *AuthenticationComponent) LogoutAccount(ctx context.Context, Id uint32) error {
	if err, _ := c.isValidID(Id); err != nil {
		c.Logger.Error(err.Error())
		return nil
	}

	// TODO: think about how to handle this failed call
	if err := c.Client.LogOutAccount(); err != nil {
		c.Logger.Error(err.Error())
		return err
	}

	return nil
}
