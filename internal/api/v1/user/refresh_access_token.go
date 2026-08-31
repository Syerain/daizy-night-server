package v1

//"github.com/atomreforge/daizy-night-server/internal/model"

type RefrehAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
