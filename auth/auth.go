package auth

import (
	"context"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
)

//encore:authhandler
func AuthHandler(ctx context.Context, token string) (auth.UID, error) {
	if token != "bearer token" {
		return "", &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "invalid token",
		}
	}

	return "1234567", nil
}
