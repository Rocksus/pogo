package gauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/nickylogan/go-log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/Rocksus/pogo/configs"
)

// TODO: need to revamp google tokens to be on user level, not service level

type GoogleAuth interface {
	GetClient() *http.Client
}

type gAuth struct {
	client *http.Client
	scopes []string
}

type ScopeOption func() []string

func New(cfg configs.GoogleConfig, scopes ...ScopeOption) (GoogleAuth, error) {
	authScopes := make([]string, 0)
	for _, s := range scopes {
		authScopes = append(authScopes, s()...)
	}
	config, err := google.ConfigFromJSON([]byte(cfg.Credentials),
		authScopes...)
	if err != nil {
		return nil, err
	}
	client, err := getClient(config, cfg.Token)
	if err != nil {
		return nil, err
	}

	return &gAuth{
		scopes: authScopes,
		client: client,
	}, nil
}

func (g *gAuth) GetClient() *http.Client {
	return g.client
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, token string) (*http.Client, error) {
	tok, err := tokenFromString(token)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		if err := saveToken("token.json", tok); err != nil {
			log.WithError(err).Errorln("Unable to save token")
		}
	}
	return config.Client(context.Background(), tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.WithError(err).Errorln("Unable to read authorization code")
		return nil, err
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.WithError(err).Errorln("Unable to retrieve token from web")
		return nil, err
	}
	return tok, nil
}

// Retrieves a token from a local file.
func tokenFromString(tokenString string) (*oauth2.Token, error) {
	if tokenString == "" {
		return nil, errors.New("no token")
	}
	tok := &oauth2.Token{}
	err := json.Unmarshal([]byte(tokenString), tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	log.Infof("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.WithError(err).Errorln("Unable to cache oauth token")
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}
