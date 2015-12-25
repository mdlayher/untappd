package untappd

import (
	"time"
)

// Toast represents an Untappd toast to a Checkin, and contains metadata
// regarding the toast, and the User who performed the toast.
type Toast struct {
	// Metadata from Untappd.
	ID     int
	UserID int

	// Time when this toast was submitted to Untappd.
	Created time.Time

	// The user who performed the Toast.
	User *User
}

// rawToast is the raw JSON representation of an Untappd toast.  Its data is
// unmarshaled from JSON and then exported to a Toast struct.
type rawToast struct {
	ID      int          `json:"like_id"`
	UserID  int          `json:"uid"`
	Created responseTime `json:"created_at"`
	User    *rawUser     `json:"user"`
}

// export creates an exported Toast from a rawToast struct, allowing for more
// useful structures to be created for client consumption.
func (r *rawToast) export() *Toast {
	return &Toast{
		ID:      r.ID,
		UserID:  r.UserID,
		Created: time.Time(r.Created),
		User:    r.User.export(),
	}
}
