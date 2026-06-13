package dto

import "time"

type ContaboAuthResult struct {
	ResultStatus bool
	AccountId    *string
	Token        *string
	Message      string
	RawData      *string
	Result       *ContaboAuthResponse
	Error        error
}

type ContaboCreateSnapshotResult struct {
	ResultStatus bool
	Message      string
	RawData      *string
	Result       *ContaboCreateSnapshotResponse
	Error        error
}

type ContaboDeleteSnapshotResult struct {
	ResultStatus bool
	Message      string
	Error        error
}

type ContaboCreateSnapshotPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ContaboSnapshot struct {
	TenantID       string    `json:"tenantId"`
	CustomerID     string    `json:"customerId"`
	SnapshotID     string    `json:"snapshotId"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InstanceID     int64     `json:"instanceId"`
	CreatedDate    time.Time `json:"createdDate"`
	AutoDeleteDate time.Time `json:"autoDeleteDate"`
	ImageID        string    `json:"imageId"`
	ImageName      string    `json:"imageName"`
}

type ContaboMetadata struct {
	Self string `json:"self"`
}

type ContaboAuthResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type ContaboCreateSnapshotResponse struct {
	Data     []ContaboSnapshot `json:"data"`
	Metadata ContaboMetadata   `json:"_links"`
}
