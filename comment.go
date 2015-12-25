package untappd

import (
	"time"
)

// Comment represents an Untappd comment from a User to a Checkin, and
// contains metadata regarding the comment, and the User who submitted
// the comment.
type Comment struct {
	// Metadata from Untappd.
	ID        int
	CheckinID int

	// The actual comment about a Checkin.
	Comment string

	// Time when this comment was submitted to Untappd.
	Created time.Time

	// The user who submitted the Comment.
	User *User
}

// rawComment is the raw JSON representation of an Untappd toast.  Its data is
// unmarshaled from JSON and then exported to a Comment struct.
type rawComment struct {
	ID        int          `json:"comment_id"`
	CheckinID int          `json:"checkin_id"`
	Comment   string       `json:"comment"`
	Created   responseTime `json:"created_at"`
	User      *rawUser     `json:"user"`
}

// export creates an exported Comment from a rawComment struct, allowing for more
// useful structures to be created for client consumption.
func (r *rawComment) export() *Comment {
	return &Comment{
		ID:        r.ID,
		CheckinID: r.CheckinID,
		Comment:   r.Comment,
		Created:   time.Time(r.Created),
		User:      r.User.export(),
	}
}
