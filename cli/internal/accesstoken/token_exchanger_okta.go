package accesstoken

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// OktaTokenExchanger exchanges a code for a token with Okta directly.
type OktaTokenExchanger struct {
	config       *oauth2.Config
	codeVerifier string
}

var _ TokenExchanger = &OktaTokenExchanger{}

// NewOktaTokenExchanger creates a new token exchanger.
func NewOktaTokenExchanger(c *oauth2.Config, cv string) (*OktaTokenExchanger, error) {
	return &OktaTokenExchanger{
		config:       c,
		codeVerifier: cv,
	}, nil
}

// ObtainToken obtains a token from the issuer.
func (e *OktaTokenExchanger) ObtainToken(ctx context.Context, code string) error {
	token, err := e.config.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", e.codeVerifier))
	if err != nil {
		return fmt.Errorf("failed to get token: %v", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return fmt.Errorf("no id_token in token response")
	}

	provider, err := newOIDCProvider(ctx, strings.TrimSuffix(e.config.Endpoint.TokenURL, "/v1/token"))
	if err != nil {
		return fmt.Errorf("failed to get provider: %v", err)
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: e.config.ClientID})
	if _, err := verifier.Verify(ctx, rawIDToken); err != nil {
		return fmt.Errorf("failed to verify ID token: %v", err)
	}

	if err := saveToken(token); err != nil {
		return fmt.Errorf("failed to save token token: %v", err)
	}

	return nil
}
