package server

import (
	"net/http"
)

// accessControl applies basic access control rules to user routes.
// These basically implement "authenticated only" for all
// modifications of user account information and views of private
// information. Returns the user ID being operated on.

func accessControl(r *http.Request) *int {
	// Get authentication information from context.
	authInfo := AuthInfoFromContext(r.Context())

	// Not an authenticated user.
	if !authInfo.Authenticated {
		return nil
	}

	// Authenticated user: return the user ID.
	return &authInfo.UserID
}
