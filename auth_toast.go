package untappd

import (
	"fmt"
	"net/http"
	"net/url"
)

// ToastRequest represents a request to toast an Untappd checkin
type ToastRequest struct {
	// Mandatory parameters
	CheckinID int
}

// Toast toasts a checkin specified by the input ToastRequest struct.
func (a *AuthService) Toast(r ToastRequest) (*http.Response, error) {
	// Add required parameters
	q := url.Values{}

	// Temporary struct to unmarshal checkin JSON
	var v struct {
		Response rawToast `json:"response"`
	}

	// Perform request to toast a checkin
	res, err := a.client.request("POST", fmt.Sprintf("checkin/toast/%d", r.CheckinID), q, nil, &v)
	if err != nil {
		return res, err
	}

	return res, nil
}
