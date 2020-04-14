package model

// TherapistInfo is an aggregate used for convenience in database
// functions.
type TherapistInfo struct {
	Base              *Therapist
	Profile           *TherapistProfile
	Image             *Image
	HasPublicProfile  bool
	HasPendingProfile bool
}
