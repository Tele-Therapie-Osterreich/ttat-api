package model

import "github.com/Tele-Therapie-Osterreich/ttat-api/model/types"

// UserFullProfile is a JSON view of the full profile for a user.
type UserFullProfile struct {
	ID            int                 `json:"id"`
	Email         string              `json:"email"`
	Type          types.TherapistType `json:"type"`
	Name          *string             `json:"name"`
	StreetAddress *string             `json:"street_address"`
	City          *string             `json:"city"`
	Postcode      *string             `json:"postcode"`
	Country       *string             `json:"country"`
	Phone         *string             `json:"phone"`
	ShortProfile  *string             `json:"short_profile"`
	FullProfile   *string             `json:"full_profile"`
	Photo         *string             `json:"photo"`
	// TODO: ADD SUB-SPECIALITIES
}

// UserFullProfileFromUser creates a full profile view from a user.
func UserFullProfileFromUser(user *User, photo *Image) *UserFullProfile {
	var photoLink *string
	if photo != nil {
		photoLink = photo.MakeLink()
	}
	return &UserFullProfile{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		StreetAddress: user.StreetAddress,
		City:          user.City,
		Postcode:      user.Postcode,
		Country:       user.Country,
		Phone:         user.Phone,
		ShortProfile:  user.ShortProfile,
		FullProfile:   user.FullProfile,
		Photo:         photoLink,
	}
}
