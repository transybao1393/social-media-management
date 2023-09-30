package domain

// - General response
type Response struct {
	StatusCode interface{} `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// - Authentication Response
type AuthParams struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}
