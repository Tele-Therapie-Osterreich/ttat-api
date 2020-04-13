package model

import "github.com/Tele-Therapie-Osterreich/ttat-api/model/types"

// TherapistFullProfile is a JSON view of the full profile for a therapist.
type TherapistFullProfile struct {
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
	Status        types.ApprovalState `json:"status"`
	Photo         *string             `json:"photo"`
	// TODO: ADD SUB-SPECIALITIES
}

// TherapistFullProfileFromTherapist creates a full profile view from a therapist.
func TherapistFullProfileFromTherapist(th *Therapist, photo *Image) *TherapistFullProfile {
	var photoLink *string
	if photo != nil {
		photoLink = photo.MakeLink()
	}
	return &TherapistFullProfile{
		ID:            th.ID,
		Email:         th.Email,
		Name:          th.Name,
		StreetAddress: th.StreetAddress,
		City:          th.City,
		Postcode:      th.Postcode,
		Country:       th.Country,
		Phone:         th.Phone,
		ShortProfile:  th.ShortProfile,
		FullProfile:   th.FullProfile,
		Status:        th.Status,
		Photo:         photoLink,
	}
}
