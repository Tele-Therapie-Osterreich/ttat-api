package server

import (
	"net/http"
)

// accessControl applies access control rules to user routes. These
// basically implement "session authentication only" + "owner or
// admin" access rules for all access to user account information.
// Returns the user ID being operated on, the user ID of the user
// accessing the route and the admin flag of the user accessing the
// route.

// TODO: MAKE THIS A MIDDLEWARE.
func accessControl(r *http.Request) (*int, *int) {
	// Get ID from URL parameters.
	paramUserID := getIntParam(r, "id")

	// Get authentication information from context.
	authInfo := AuthInfoFromContext(r.Context())

	// Logic:

	// The user ID we are trying to operate on is either from the
	// authenticated user, or from the URL parameter if it's there.
	actionUserID := authInfo.UserID
	if paramUserID != nil {
		actionUserID = *paramUserID
	}

	// If the user ID that we're trying to operate on is different from
	// the requesting authenticated user ID, then the user must be an
	// administrator.
	if actionUserID != authInfo.UserID { // && !authInfo.UserIsAdmin {
		return nil, nil // , false
	}

	// Operation is allowed: return the user ID we're operating on.
	return &actionUserID, &authInfo.UserID // , authInfo.UserIsAdmin
}
