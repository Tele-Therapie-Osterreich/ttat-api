package model

import (
	"time"

	"github.com/Tele-Therapie-Osterreich/ttat-api/model/types"
)

// TherapistSummaryView is a JSON view of the summary public profile
// for a therapist.
type TherapistSummaryView struct {
	ID            int                 `json:"id"`
	Email         string              `json:"email"`
	Type          types.TherapistType `json:"type"`
	Name          *string             `json:"name"`
	StreetAddress *string             `json:"street_address",omitempty`
	City          *string             `json:"city",omitempty`
	Postcode      *string             `json:"postcode",omitempty`
	Country       *string             `json:"country",omitempty`
	Phone         *string             `json:"phone",omitempty`
	Website       *string             `json:"website",omitempty`
	Languages     []string            `json:"languages",omitempty`
	ShortProfile  *string             `json:"short_profile",omitempty`
	Photo         *string             `json:"photo",omitempty`
	// TODO: ADD SUB-SPECIALITIES
}

// TherapistFullView is a JSON view of the full public profile for a
// therapist.
type TherapistFullView struct {
	TherapistSummaryView
	FullProfile *string `json:"full_profile",omitempty`
}

// TherapistLoginView is a JSON view of all the information about a
// therapist that is sent on login. This is basically the full profile
// plus information about the therapist's status.
type TherapistLoginView struct {
	TherapistFullView
	Status      types.ApprovalState `json:"status"`
	LastLoginAt time.Time           `json:"last_login_at"`
	EditedAt    time.Time           `json:"edited_at"`
	Public      bool                `json:"public"`
	HasPublic   *bool               `json:"has_public",omitempty`
	HasPending  *bool               `json:"has_pending",omitempty`
}

// Fill fills in a summary profile view from a therapist and their
// public profile.
func (v *TherapistSummaryView) Fill(info *TherapistInfo) {
	var photoLink *string
	if info.Image != nil {
		photoLink = info.Image.MakeLink()
	}
	v.ID = info.Base.ID
	v.Email = info.Base.Email
	v.Name = info.Profile.Name
	v.StreetAddress = info.Profile.StreetAddress
	v.City = info.Profile.City
	v.Postcode = info.Profile.Postcode
	v.Country = info.Profile.Country
	v.Phone = info.Profile.Phone
	v.Website = info.Profile.Website
	v.Languages = info.Profile.Languages
	v.ShortProfile = info.Profile.ShortProfile
	v.Photo = photoLink

	// TODO: ADD SPECIALITIES
}

// NewTherapistSummaryView creates a summary profile view from a
// therapist and their public profile.
func NewTherapistSummaryView(info *TherapistInfo) *TherapistSummaryView {
	v := &TherapistSummaryView{}
	v.Fill(info)
	return v
}

// NewTherapistFullView creates a full profile view from a therapist
// and their public profile.
func NewTherapistFullView(info *TherapistInfo) *TherapistFullView {
	v := &TherapistFullView{}
	v.Fill(info)
	v.FullProfile = info.Profile.FullProfile
	return v
}

// NewTherapistLoginView creates a private login view from a therapist
// and their public profile.
func NewTherapistLoginView(info *TherapistInfo) *TherapistLoginView {
	v := &TherapistLoginView{}
	v.Fill(info)
	v.FullProfile = info.Profile.FullProfile
	v.Status = info.Base.Status
	v.LastLoginAt = info.Base.LastLoginAt
	v.EditedAt = info.Profile.EditedAt
	v.Public = info.Profile.Public
	if info.Profile.Public && info.HasPendingProfile {
		hasPending := true
		v.HasPending = &hasPending
	}
	if !info.Profile.Public && info.HasPublicProfile {
		hasPublic := true
		v.HasPublic = &hasPublic
	}
	return v
}
