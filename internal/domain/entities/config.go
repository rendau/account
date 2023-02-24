package entities

type ConfigSt struct {
	RefreshTokenDurSeconds int64 `json:"refresh_token_dur_seconds"`
	AccessTokenDurSeconds  int64 `json:"access_token_dur_seconds"`
}
