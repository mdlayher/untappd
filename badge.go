package untappd

import (
	"net/url"
	"time"
)

// Badge represents an Untappd badge, and contains information regarding its name,
// description, when it was earned, and various media associated with the badge.
type Badge struct {
	// Metadata from Untappd.
	ID          int
	CheckinID   int
	Name        string
	Description string
	Hint        string
	Active      bool

	// Links to images of the badge.
	Media BadgeMedia

	// If applicable, time when the specified user earned this badge.
	Earned time.Time

	// If applicable, badge levels which the specified user has obtained.
	// If the slice has zero length, no levels exist for this badge.
	Levels []*Badge
}

// BadgeMedia contains links to media regarding a Badge.  Included are links
// to a small, medium, and large image for a given Badge.
type BadgeMedia struct {
	SmallImage  url.URL
	MediumImage url.URL
	LargeImage  url.URL
}

// rawBadge is the raw JSON representation of an Untappd badge.  Its data is
// unmarshaled from JSON and then exported to a Badge struct.
type rawBadge struct {
	ID          int                 `json:"badge_id"`
	CheckinID   int                 `json:"checkin_id"`
	Name        string              `json:"badge_name"`
	Description string              `json:"badge_description"`
	Hint        string              `json:"badge_hint"`
	Active      responseBool        `json:"badge_active_status"`
	Media       rawBadgeMedia       `json:"media"`
	Earned      responseTime        `json:"created_at"`
	Levels      responseBadgeLevels `json:"levels"`
}

// export creates an exported Badge from a rawBadge struct, allowing for more
// useful structures to be created for client consumption.
func (r *rawBadge) export() *Badge {
	b := &Badge{
		ID:          r.ID,
		CheckinID:   r.CheckinID,
		Name:        r.Name,
		Description: r.Description,
		Hint:        r.Hint,
		Active:      bool(r.Active),
		Media:       r.Media.export(),
		Earned:      time.Time(r.Earned),
	}

	// Export badge levels as a slice of badges belonging to parent badge
	levels := make([]*Badge, r.Levels.Count)
	for i := range r.Levels.Items {
		levels[i] = r.Levels.Items[i].export()
	}
	b.Levels = levels

	return b
}

// rawBadgeMedia is the raw JSON representation of Untappd badge media.  Its data is
// unmarshaled from JSON and then exported to a BadgeMedia struct.
type rawBadgeMedia struct {
	SmallImage  responseURL `json:"badge_image_sm"`
	MediumImage responseURL `json:"badge_image_md"`
	LargeImage  responseURL `json:"badge_image_lg"`
}

// export creates an exported BadgeMedia from a rawBadgeMedia struct, allowing
// for more useful structures to be created for client consumption.
func (r *rawBadgeMedia) export() BadgeMedia {
	return BadgeMedia{
		SmallImage:  url.URL(r.SmallImage),
		MediumImage: url.URL(r.MediumImage),
		LargeImage:  url.URL(r.LargeImage),
	}
}
