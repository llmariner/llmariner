package accesstoken

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/llmariner/llmariner/cli/internal/configs"
	"golang.org/x/oauth2"
)

// TokenExchanger exchanges a code for a token.
type TokenExchanger interface {
	// ObtainToken obtains a token from the issuer.
	ObtainToken(ctx context.Context, code string) error
}

// DexTokenExchanger exchanges a code for a token via Dex.
type DexTokenExchanger struct {
	auth configs.Auth
}

var _ TokenExchanger = &DexTokenExchanger{}

// NewDexTokenExchanger creates a new token exchanger.
func NewDexTokenExchanger(c *configs.C) (*DexTokenExchanger, error) {
	return &DexTokenExchanger{
		auth: c.Auth,
	}, nil
}

// LoginURL returns a URL to login.
func (e *DexTokenExchanger) LoginURL() (string, error) {
	iu, err := url.Parse(e.auth.IssuerURL)
	if err != nil {
		return "", fmt.Errorf("parse issuer-url: %v", err)
	}

	iu.Path = path.Join(iu.Path, "auth")
	q := iu.Query()
	q.Add("client_id", e.auth.ClientID)
	q.Add("redirect_uri", e.auth.RedirectURI)
	q.Add("response_type", "code")
	// TODO(kenji): Remove unnecessary scopes.
	// "offline_access" for refresh token.
	q.Add("scope", "openid profile email offline_access")
	iu.RawQuery = q.Encode()
	return iu.String(), nil
}

// ObtainToken obtains a token from the issuer.
func (e *DexTokenExchanger) ObtainToken(ctx context.Context, code string) error {
	oauth2Config, err := NewOauth2Config(ctx, e.auth)
	if err != nil {
		return fmt.Errorf("new oauth2 config: %v", err)
	}

	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to get token: %v", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return fmt.Errorf("no id_token in token response")
	}

	provider, err := newOIDCProvider(ctx, e.auth.IssuerURL)
	if err != nil {
		return fmt.Errorf("failed to get provider: %v", err)
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: e.auth.ClientID})
	if _, err := verifier.Verify(ctx, rawIDToken); err != nil {
		return fmt.Errorf("failed to verify ID token: %v", err)
	}

	if err := saveToken(token); err != nil {
		return fmt.Errorf("failed to save token token: %v", err)
	}

	return nil
}

func refreshTokenIfExpired(ctx context.Context, token T, auth configs.Auth) (T, error) {
	if token.TokenExpiry.After(time.Now()) {
		// No need to refresh.
		return token, nil
	}

	oauth2Config, err := NewOauth2Config(ctx, auth)
	if err != nil {
		return T{}, fmt.Errorf("new oauth2 config: %v", err)
	}

	savedToken := &oauth2.Token{
		TokenType:    token.TokenType,
		Expiry:       token.TokenExpiry,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	tokenSource := oauth2Config.TokenSource(ctx, savedToken)
	newToken, err := tokenSource.Token()
	if err != nil {
		return T{}, fmt.Errorf("failed to get token: %v", err)
	}
	if err := saveToken(newToken); err != nil {
		return token, err
	}
	return token, nil
}

func newOIDCProvider(ctx context.Context, issuerURL string) (*oidc.Provider, error) {
	ctx = oidc.ClientContext(ctx, http.DefaultClient)
	return oidc.NewProvider(ctx, issuerURL)
}

// NewOauth2Config creates a new oauth2 config.
func NewOauth2Config(ctx context.Context, auth configs.Auth) (*oauth2.Config, error) {
	provider, err := newOIDCProvider(ctx, auth.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	return &oauth2.Config{
		ClientID:     auth.ClientID,
		ClientSecret: auth.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  auth.RedirectURI,
	}, nil
}
