package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"tiktok_api/app/logger"
	"tiktok_api/youtube/repository/redis"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

var log = logger.NewLogrusLogger()
var ctx = context.Background()
var config *oauth2.Config

func init() {
	path, err := os.Getwd()
	if err != nil {
		handleError(err, "Getwd is error", "error")
	}
	clientSecretPath := filepath.Join(path, "/app/config/client_secret_1.json")
	b, err := ioutil.ReadFile(clientSecretPath) //- current directory json file
	if err != nil {
		handleError(err, "Unable to read Youtube client secret file", "error")
	}

	//- If modifying these scopes, delete your previously saved credentials
	//- at ~/.credentials/youtube-go-quickstart.json
	config, err = google.ConfigFromJSON(b,
		youtube.YoutubeReadonlyScope,
		youtube.YoutubeUploadScope,
		youtube.YoutubeScope,
	)
	if err != nil {
		handleError(err, "Google cannot load config from json file", "error")
	}
	log.Printf("Load Youtube Client Secret file successfully at path: %s", clientSecretPath)
}

func YoutubeOAuthCodeExchange(clientKey string, code string) string {
	successURL := "https://www.meritize.com/wp-content/uploads/2018/06/keys-to-success.jpg"
	failURL := "https://www.freecodecamp.org/news/content/images/2023/02/image-126.png"
	tokens, err := config.Exchange(ctx, code)
	if err != nil {
		handleError(err, "Unable to retrieve token from web", "error")
		return failURL
	}

	//- check if clientKey is exist
	//- if exists, update value with clientKey
	//- if not exist, CSRF => panic
	isClientKeyExist := redis.IsExist(clientKey)
	if !isClientKeyExist {
		handleError(errors.New("redis.IsExist"), "Error when call redis.IsExist", "error")
		return failURL
	}

	//- if exist
	//- update tokens for clientKey
	youtubeOAuth := redis.GetClientByClientKey(clientKey)
	youtubeOAuth.AccessToken = tokens.AccessToken
	youtubeOAuth.RefreshToken = tokens.RefreshToken
	youtubeOAuth.Expiry = tokens.Expiry

	isUpdate := redis.UpdateYoutubeByClientKey(clientKey, youtubeOAuth)
	if !isUpdate {
		handleError(errors.New("redis.UpdateYoutubeByClientKey"), "Error when call UpdateYoutubeByClientKey", "error")
		return failURL
	}

	return successURL
}

func BuildClientFromTokens(clientKey string) *http.Client {
	//- FIXME: check if clientKey is exist

	//- get client tokens base on clientKey
	youtubeOAuth := redis.GetClientByClientKey(clientKey)
	tokens := &oauth2.Token{
		AccessToken:  youtubeOAuth.AccessToken,
		RefreshToken: youtubeOAuth.RefreshToken,
		Expiry:       youtubeOAuth.Expiry,
		TokenType:    "Bearer",
	}
	return config.Client(ctx, tokens)
}

func BuildServiceFromToken(clientKey string) *youtube.Service {
	//- FIXME: check if clientKey is exist

	//- get client tokens base on clientKey
	youtubeOAuth := redis.GetClientByClientKey(clientKey)
	tokens := &oauth2.Token{
		AccessToken:  youtubeOAuth.AccessToken,
		RefreshToken: youtubeOAuth.RefreshToken,
		Expiry:       youtubeOAuth.Expiry,
		TokenType:    "Bearer",
	}
	client := config.Client(ctx, tokens)
	service, err := youtube.New(client)
	if err != nil {
		handleError(err, "Unable to create Youtube service", "error")
	}
	return service
}

func GetAuthURL() (string, string) {
	clientKey, _ := redis.CreateNewYoutubeClient(config.ClientID, config.ClientSecret)
	state := clientKey
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return authURL, clientKey
}

// - tokenCacheFile generates credential file path/filename.
// - It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("client_credential.json")), err
}

// - tokenFromFile retrieves a Token from a given file path.
// - It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// - saveToken uses a file path to create a file and store the
// - token in it.
func saveToken(file string, token *oauth2.Token) {
	log.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		handleError(err, "Unable to cache oauth token", "error")
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
