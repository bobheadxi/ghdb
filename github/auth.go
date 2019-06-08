package github

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

func newClient() *http.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	return oauth2.NewClient(context.Background(), src)
}

// AppAuth is a container for authorization configuration
type AppAuth struct {
	AppID          string
	ClientID       string
	ClientSecret   string
	SigningKeyPath string
}

// NewEnvAuth instantiates authentication from environment variables
func NewEnvAuth() *AppAuth {
	return &AppAuth{
		ClientID:       os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret:   os.Getenv("GITHUB_CLIENT_SECRET"),
		SigningKeyPath: os.Getenv("GITHUB_APP_KEY"),
		AppID:          os.Getenv("GITHUB_APP_ID"),
	}
}

// Token implements oauth2.TokenSource, and is used as an autogenerating token
// source
func (a *AppAuth) Token() (*oauth2.Token, error) {
	t, exp, err := a.generateJWT()
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken: t,
		Expiry:      *exp,
	}, nil
}

// generateJWT signs a new JWT for use with the GitHub API
func (a *AppAuth) generateJWT() (string, *time.Time, error) {
	priv, err := ioutil.ReadFile(a.SigningKeyPath)
	if err != nil {
		return "", nil, fmt.Errorf("could not read singing key: %s", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		return "", nil, fmt.Errorf("could not parse signing key: %s", err)
	}

	var expiry = time.Now().Add(time.Minute)
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: expiry.Unix(),
		Issuer:    a.AppID,
	}).SignedString(key)
	if err != nil {
		return "", nil, fmt.Errorf("could not sign token: %s", err)
	}

	return token, &expiry, nil
}
