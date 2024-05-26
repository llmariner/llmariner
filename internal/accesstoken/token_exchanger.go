package accesstoken

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/llm-operator/cli/internal/configs"
	"github.com/llm-operator/cli/internal/ui"
	"golang.org/x/oauth2"
)

// NewTokenExchanger creates a new token exchanger.
func NewTokenExchanger(c *configs.C) (*TokenExchanger, error) {
	var issuerResolvedAddr string

	// Check if the issuer URL is resolvable. If not, fall back to the endpoint URL.
	iu, err := url.Parse(c.Auth.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("parse issuer-url: %v", err)
	}

	if _, err := net.LookupIP(iu.Host); err != nil {
		ep, err := url.Parse(c.EndpointURL)
		if err != nil {
			return nil, err
		}
		issuerResolvedAddr = ep.Host
	}

	return &TokenExchanger{
		auth: c.Auth,

		issuerHost:         iu.Host,
		issuerResolvedAddr: issuerResolvedAddr,
	}, nil
}

// TokenExchanger exchanges a code for a token.
type TokenExchanger struct {
	auth configs.Auth

	issuerHost         string
	issuerResolvedAddr string
	warningPrinted     bool
}

// LoginURL returns a URL to login.
func (e *TokenExchanger) LoginURL() (string, error) {
	iu, err := url.Parse(e.auth.IssuerURL)
	if err != nil {
		return "", fmt.Errorf("parse issuer-url: %v", err)
	}

	if addr := e.issuerResolvedAddr; addr != "" {
		e.printWarning()
		iu.Host = addr
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
func (e *TokenExchanger) ObtainToken(ctx context.Context, code string) error {
	provider, err := e.newOIDCProvider(ctx)
	if err != nil {
		return fmt.Errorf("failed to get provider: %v", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     e.auth.ClientID,
		ClientSecret: e.auth.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  e.auth.RedirectURI,
	}

	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to get token: %v", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return fmt.Errorf("no id_token in token response")
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: e.auth.ClientID})
	if _, err := verifier.Verify(ctx, rawIDToken); err != nil {
		return fmt.Errorf("failed to verify ID token: %v", err)
	}

	tokenType, ok := token.Extra("token_type").(string)
	if !ok {
		return fmt.Errorf("no token_type in token response")
	}

	accessToken, ok := token.Extra("access_token").(string)
	if !ok {
		return fmt.Errorf("no access_token in token response")
	}

	refreshToken, ok := token.Extra("refresh_token").(string)
	if !ok {
		return fmt.Errorf("no refresh_token in token response")
	}

	if err := saveToken(&T{
		TokenType:    tokenType,
		TokenExpiry:  token.Expiry,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}); err != nil {
		return fmt.Errorf("failed to save token token: %v", err)
	}

	return nil
}

func (e *TokenExchanger) refreshTokenIfExpired(ctx context.Context, token T) (T, error) {
	if token.TokenExpiry.After(time.Now()) {
		// No need to refresh.
		return token, nil
	}

	provider, err := e.newOIDCProvider(ctx)
	if err != nil {
		return T{}, fmt.Errorf("failed to get provider: %v", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     e.auth.ClientID,
		ClientSecret: e.auth.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  e.auth.RedirectURI,
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
	token = T{
		TokenType:    newToken.TokenType,
		TokenExpiry:  newToken.Expiry,
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
	}
	if err := saveToken(&token); err != nil {
		return token, err
	}
	return token, nil
}

func (e *TokenExchanger) newOIDCProvider(ctx context.Context) (*oidc.Provider, error) {
	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		// Replace the addr with the resolved address if it is set.
		if e.issuerResolvedAddr != "" && addr == fmt.Sprintf("%s:80", e.issuerHost) {
			e.printWarning()

			if strings.Contains(e.issuerResolvedAddr, ":") {
				addr = e.issuerResolvedAddr
			} else {
				addr = fmt.Sprintf("%s:80", e.issuerResolvedAddr)
			}
		}
		d := net.Dialer{}
		return d.DialContext(ctx, network, addr)
	}
	ctx = oidc.ClientContext(ctx, http.DefaultClient)
	return oidc.NewProvider(ctx, e.auth.IssuerURL)
}

func (e *TokenExchanger) printWarning() {
	if e.warningPrinted {
		return
	}
	p := ui.NewPrompter()
	p.Warn(fmt.Sprintf("Unable to resolve the issuer address (%q) while obtaining an OIDC access token. Fallling back to the endpoint address (%q)", e.issuerHost, e.issuerResolvedAddr))
	e.warningPrinted = true
}
