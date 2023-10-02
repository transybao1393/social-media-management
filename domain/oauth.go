package domain

import (
	"time"
)

type OAuth struct {
	//- compulsory fields
	TenantId  string `json:"tenantId" bson:"tenantId"`
	AppId     string `json:"appId" bson:"appId"`
	ApiKey    string `json:"apiKey" bson:"apiKey"`
	Subdomain string `json:"subdomain" bson:"subdomain"`

	AccessToken  string `json:"access_token,omitempty" bson:"access_token,omitempty"`
	UserToken    string `json:"user_token,omitempty" bson:"user_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`

	//- need these to build OAuth2 link
	ClientId     string   `json:"clientId" bson:"clientId" `
	ClientSecret string   `json:"clientSecret" bson:"clientSecret"`
	Scopes       []string `json:"scopes" bson:"scopes"`

	//- expiry
	ExpiresIn time.Time `json:"expiresIn,omitempty" bson:"expiresIn,omitempty"`
	Expiry    time.Time `json:"-"`

	//- Redirect Url
	RedirectUrlSuccess string `json:"redirectUrlSuccess" bson:"redirectUrlSuccess"`
	RedirectUrlError   string `json:"redirectUrlError" bson:"redirectUrlError"`
}

func (o *OAuth) SetExpiry() {
	// To prevent last minute expirations, the expiration date will be accelerated by 10 minutes.
	o.Expiry = o.ExpiresIn.Add(-5 * time.Minute)
}
