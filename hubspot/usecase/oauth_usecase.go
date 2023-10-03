package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"tiktok_api/domain"
	"time"

	utilhttp "tiktok_api/app/utils/http"
	redisRepository "tiktok_api/hubspot/repository/redis"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func ListHubspotObjectFieldsUseCase(accessToken string) ([]*domain.SinglePropertyInfo, error) {
	// token, err := redisRepository.GetOneByTenantIdApiKeyType(config.TenantId, config.ApiKey)
	// if err != nil {
	// 	return nil, err
	// }

	pathURL := &utilhttp.PathURL{
		APIDomain: viper.GetString("HUBSPOT.API_URL"),
		APIURI:    fmt.Sprintf("%s/%s/properties", "/properties/v2", "calls"),
	}
	purl := pathURL.New().Build()
	ri := &utilhttp.CustomRequest{
		MethodName:  "GET",
		PathURL:     purl,
		AccessToken: accessToken,
	}

	byteData, err := ri.Exec()
	if err != nil {
		return nil, err
	}

	var cor []*domain.SinglePropertyInfo
	err = json.NewDecoder(bytes.NewReader(byteData)).Decode(&cor)
	if err != nil {
		return nil, err
	}

	return cor, nil
}

func UpdateTokenUseCase(config *domain.OAuth) (string, error) {
	//- Save to database
	token, err := redisRepository.GetOneByTenantIdApiKeyType(config.TenantId, config.ApiKey)
	if err != nil {
		return "", err
	}

	token.ClientId = config.ClientId
	token.ClientSecret = config.ClientSecret
	token.Subdomain = config.Subdomain
	token.RedirectUrlSuccess = config.RedirectUrlSuccess
	token.RedirectUrlError = config.RedirectUrlError
	token.Scopes = config.Scopes

	// - Update data into redis
	isUpdate := redisRepository.UpdateTenantDataBy(config.TenantId, config.ApiKey, token)
	if !isUpdate {
		return "", errors.New("Update token failed")
	}

	fixedScopes := []string{
		"crm.objects.contacts.read",
		"crm.objects.contacts.write",
		"crm.objects.custom.read",
		"crm.objects.custom.write",
		"crm.lists.read",
		"crm.objects.companies.read",
		"crm.objects.companies.write",
		"crm.objects.deals.read",
		"crm.objects.marketing_events.read",
		"crm.objects.quotes.read",
		"crm.objects.deals.write",
		"crm.objects.marketing_events.write",
		"crm.objects.quotes.write",
		"crm.lists.write",
		"crm.objects.line_items.write",
		"crm.objects.line_items.read",
		"crm.schemas.deals.write",
		"crm.schemas.companies.write",
		"crm.schemas.contacts.write",
	}

	queryParams := fmt.Sprintf("client_id=%[1]s&redirect_uri=%[2]s&scope=%[3]s&state=%[4]s", token.ClientId, viper.GetString("HUBSPOT.REDIRECT_URL"), strings.Join(fixedScopes, " "), fmt.Sprintf("%s-%s", config.TenantId, config.ApiKey))
	finalOAuth2URL := fmt.Sprintf("%[1]s?%[2]s", viper.GetString("HUBSPOT.AUTH_URL"), queryParams)

	return finalOAuth2URL, nil
}

var hubspotOAuthConfig = &oauth2.Config{
	RedirectURL:  "",
	ClientID:     "",
	ClientSecret: "",
	Scopes:       nil,
	Endpoint: oauth2.Endpoint{
		AuthURL:   "",
		TokenURL:  "",
		AuthStyle: 1,
	},
}

func OAuthHubspotCallbackUseCase(id string, code string) (string, error) {
	oauthInfo, err := redisRepository.GetOneById(id)
	if err != nil {
		return "", err
	}

	hubspotOAuthConfig.ClientID = oauthInfo.ClientId
	hubspotOAuthConfig.ClientSecret = oauthInfo.ClientSecret
	hubspotOAuthConfig.RedirectURL = viper.GetString("HUBSPOT.REDIRECT_URL")
	hubspotOAuthConfig.Endpoint.AuthURL = viper.GetString("HUBSPOT.AUTH_URL")
	hubspotOAuthConfig.Endpoint.TokenURL = viper.GetString("HUBSPOT.TOKEN_URL")

	//- get tokens from code
	tokens, err := hubspotOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Code exchange is having errors %s", err.Error())
		return oauthInfo.RedirectUrlError, nil
	}

	//- reassign pointers and mutable
	//- set expiration time
	internalOAuth := &domain.OAuth{
		TenantId:           oauthInfo.TenantId,
		ApiKey:             oauthInfo.ApiKey,
		ClientId:           oauthInfo.ClientId,
		ClientSecret:       oauthInfo.ClientSecret,
		Subdomain:          oauthInfo.Subdomain,
		Scopes:             oauthInfo.Scopes,
		AppId:              oauthInfo.AppId,
		RedirectUrlSuccess: oauthInfo.RedirectUrlSuccess,
		RedirectUrlError:   oauthInfo.RedirectUrlError,
		AccessToken:        tokens.AccessToken,
		RefreshToken:       tokens.RefreshToken,
		ExpiresIn:          tokens.Expiry,
	}
	internalOAuth.SetExpiry()
	isUpdated := redisRepository.UpdateTokensById(id, internalOAuth)
	if !isUpdated {
		fmt.Printf("Update token failed errors %s", err.Error())
		return oauthInfo.RedirectUrlError, nil
	}

	return oauthInfo.RedirectUrlSuccess, nil
}

func SetAndUpdateAccessTokenUseCase(oauthInfo *domain.OAuth) (string, error) {
	config := &domain.OAuthConfig{
		GrantType:    "refresh_token",
		ClientId:     oauthInfo.ClientId,
		ClientSecret: oauthInfo.ClientSecret,
		RefreshToken: oauthInfo.RefreshToken,
	}
	//- make request to get new access_token
	tokensModel, err := RefreshTokenUseCase(config)
	if err != nil {
		return "", err
	}

	internalOAuth := &domain.OAuth{
		TenantId:           oauthInfo.TenantId,
		ApiKey:             oauthInfo.ApiKey,
		ClientId:           oauthInfo.ClientId,
		ClientSecret:       oauthInfo.ClientSecret,
		Subdomain:          oauthInfo.Subdomain,
		Scopes:             oauthInfo.Scopes,
		AppId:              oauthInfo.AppId,
		RedirectUrlSuccess: oauthInfo.RedirectUrlSuccess,
		RedirectUrlError:   oauthInfo.RedirectUrlError,
		AccessToken:        tokensModel.AccessToken,
		RefreshToken:       tokensModel.RefreshToken,
		ExpiresIn:          time.Now().Add(30 * time.Minute),
	}
	internalOAuth.SetExpiry()

	isUpdated := redisRepository.UpdateTokensById(fmt.Sprintf("%s-%s", oauthInfo.TenantId, oauthInfo.ApiKey), internalOAuth)
	if !isUpdated {
		return "", nil
	}

	return tokensModel.AccessToken, nil
}

func RefreshTokenUseCase(config *domain.OAuthConfig) (*domain.OAuthToken, error) {
	var customHeaders = make(map[string]string, 1)
	customHeaders["Content-Type"] = "application/x-www-form-urlencoded"

	ri := &utilhttp.CustomRequest{
		MethodName:    "POST",
		PathURL:       viper.GetString("HUBSPOT.TOKEN_URL"),
		CustomHeaders: customHeaders,
		Body:          []byte(fmt.Sprintf("grant_type=%[1]s&client_id=%[2]s&client_secret=%[3]s&refresh_token=%[4]s", config.GrantType, config.ClientId, config.ClientSecret, config.RefreshToken)),
	}

	tokensByteData, err := ri.Exec()
	if err != nil {
		return nil, err
	}

	var tokensModel *domain.OAuthToken
	err = json.NewDecoder(bytes.NewReader(tokensByteData)).Decode(&tokensModel)
	if err != nil {
		return nil, err
	}

	return tokensModel, nil
}
