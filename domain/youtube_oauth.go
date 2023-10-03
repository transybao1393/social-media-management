package domain

import (
	"time"
)

type YoutubeOAuth struct {
	//- compulsory fields
	ClientKey string `json:"client_key,omitempty" bson:"client_key,omitempty"`

	AccessToken  string `json:"access_token,omitempty" bson:"access_token,omitempty"`
	UserToken    string `json:"user_token,omitempty" bson:"user_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`

	//- need these to build OAuth2 link
	ClientId     string   `json:"clientId" bson:"clientId" `
	ClientSecret string   `json:"clientSecret" bson:"clientSecret"`
	Scopes       []string `json:"scopes" bson:"scopes"`

	//- expiry
	ExpiresIn time.Time `json:"-"`
	Expiry    time.Time `json:"expiry,omitempty" bson:"expiry,omitempty"`
}

func (o *OAuth) setYoutubeExpiresIn() {
	// To prevent last minute expirations, the expiration date will be accelerated by 10 minutes.
	o.ExpiresIn = o.Expiry.Add(-5 * time.Minute)
}

func (o *OAuth) YoutubeTokenExpired() bool {
	if o.Expiry.IsZero() {
		return false
	}
	return o.Expiry.Before(timeNow())
}

type YoutubeOAuthKey struct {
	ClientKey string `json:"client_key" bson:"client_key"`
}

type YoutubeOAuthConfig struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	// ProjectId    string `json:"project_id"` //- not used
}

type YoutubeOAuthToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	Expiry       time.Time `json:"expiry"`
}
