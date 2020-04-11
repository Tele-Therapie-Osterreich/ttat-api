package server

import (
	"net/http"
)

// accessControl applies access control rules to user routes. These
// basically implement "session authentication only" + "owner or
// admin" access rules for all modifications of user account
// information. Returns the user ID being operated on and the user ID
// of the user accessing the route.

// TODO: MAKE THIS A MIDDLEWARE.
func accessControl(r *http.Request, open bool) (*int, *int) {
	// Get ID from URL parameters.
	paramUserID := getIntParam(r, "id")

	// Get authentication information from context.
	authInfo := AuthInfoFromContext(r.Context())

	// Logic:

	// For an open route:

	// If there's a user ID parameter, that's the user we're operating
	// on. If there's not, we're operating on the user that's requesting
	// the route.
	if open {
		if paramUserID == nil {
			return &authInfo.UserID, &authInfo.UserID
		}
		return paramUserID, &authInfo.UserID
	}

	// The user ID we are trying to operate on is either from the
	// authenticated user, or from the URL parameter if it's there.
	actionUserID := authInfo.UserID
	if paramUserID != nil {
		actionUserID = *paramUserID
	}

	// If the user ID that we're trying to operate on is different from
	// the requesting authenticated user ID, then the operation is not
	// allowed.
	if actionUserID != authInfo.UserID {
		return nil, nil
	}

	// Operation is allowed: return the user ID we're operating on.
	return &actionUserID, &authInfo.UserID
}
