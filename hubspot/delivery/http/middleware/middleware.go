package middleware

import (
	"context"
	"fmt"
	"io"

	"encoding/json"
	"net/http"
	"tiktok_api/domain"

	"github.com/go-chi/render"

	redisRepository "tiktok_api/hubspot/repository/redis"
	usecase "tiktok_api/hubspot/usecase"
)

func IsTokensValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config := &domain.OAuthKey{}
		//- data validation

		defer r.Body.Close()
		ByteBody, err := io.ReadAll(r.Body)
		if err != nil {
			render.JSON(w, r, domain.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    http.StatusText(http.StatusUnprocessableEntity),
				Data:       "Invalid identity management",
			})
			return
		}
		if err := json.Unmarshal(ByteBody, config); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			render.JSON(w, r, domain.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    http.StatusText(http.StatusUnprocessableEntity),
				Data:       "Invalid identity management",
			})
			return
		}
		w.Header().Set("content-type", "application/json")

		oauthInfo, err := redisRepository.GetOneByTenantIdApiKeyType(config.TenantId, config.ApiKey)
		if err != nil {
			fmt.Printf("When get one record by id %v\n", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			render.JSON(w, r, domain.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    http.StatusText(http.StatusUnprocessableEntity),
				Data:       "Invalid Tenant with APIKey",
			})
			return
		}

		if len(oauthInfo.AccessToken) <= 0 && len(oauthInfo.RefreshToken) <= 0 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			render.JSON(w, r, domain.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    http.StatusText(http.StatusUnprocessableEntity),
				Data:       fmt.Sprintf("This account has not been setup"),
			})
			return
		}

		accessToken := oauthInfo.AccessToken
		appId := oauthInfo.AppId
		subdomain := oauthInfo.Subdomain
		userToken := oauthInfo.UserToken

		if oauthInfo.Expired() && oauthInfo.RefreshToken != "" {
			fmt.Printf("Token expired, start to get new access_token from refresh_token \n")
			//- set and update new access token to id
			newAccessToken, err := usecase.SetAndUpdateAccessTokenUseCase(oauthInfo)
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				render.JSON(w, r, domain.Response{
					StatusCode: http.StatusUnprocessableEntity,
					Message:    http.StatusText(http.StatusUnprocessableEntity),
					Data:       "Set and update new token failed",
				})
				return
			}
			accessToken = newAccessToken
		}

		ctx := context.WithValue(r.Context(), "access_token", accessToken)
		ctx = context.WithValue(ctx, "user_token", userToken)
		ctx = context.WithValue(ctx, "appId", appId)
		ctx = context.WithValue(ctx, "subdomain", subdomain)
		ctx = context.WithValue(ctx, "bodyRaw", string(ByteBody))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
