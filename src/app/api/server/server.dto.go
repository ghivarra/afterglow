package server

type ServerCreateRequestDto struct {
	Id          int     `json:"id"`
	Alias       string  `json:"alias"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IpAddress   string  `json:"ipAddress"`
	AccountId   string  `json:"accountId"`
}

type ServerUpdateRequestDto struct {
	Alias       *string `json:"alias"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IpAddress   *string `json:"ipAddress"`
}
