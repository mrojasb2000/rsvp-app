package entity

// RSVP represents a guest's response to an invitation.
type RSVP struct {
	Name, Email, Phone string
	WillAttend         bool
}
