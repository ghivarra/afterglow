package account

type AccountCreateRequestDto struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ApiClientId  string `json:"apiClientId"`
	ApiClientKey string `json:"apiClientKey"`
}

type AccountCreateResponseDto struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}
