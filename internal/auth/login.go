package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/cli/browser"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/llm-operator/cli/internal/accesstoken"
	"github.com/llm-operator/cli/internal/config"
	"github.com/llm-operator/cli/internal/ui"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

type client struct {
	authConfig         config.Auth
	issuerResolvedAddr string

	listener net.Listener
}

func loginCmd() *cobra.Command {
	var (
		cli client
	)
	cmd := cobra.Command{
		Use:   "login",
		Short: "Login to LLM service",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := ui.NewPrompter()

			c, err := config.LoadOrCreate()
			if err != nil {
				return fmt.Errorf("load config: %s", err)
			}

			cli.authConfig = c.Auth

			ru, err := url.Parse(cli.authConfig.RedirectURI)
			if err != nil {
				return fmt.Errorf("parse redirect-uri: %v", err)
			}
			iu, err := url.Parse(c.Auth.IssuerURL)
			if err != nil {
				return fmt.Errorf("parse issuer-url: %v", err)
			}

			// Check if the issuer URL is resolvable. If not, fall back to the endpoint URL.
			if _, err := net.LookupIP(iu.Host); err != nil {
				ep, err := url.Parse(c.EndpointURL)
				if err != nil {
					return err
				}
				p.Warn(fmt.Sprintf("Unable to resolve the issuer address (%q). Fallling back to the endpoint address (%q)", iu.Host, ep.Host))
				cli.issuerResolvedAddr = ep.Host
			}

			if cli.issuerResolvedAddr != "" {
				iu.Host = cli.issuerResolvedAddr
			}

			iu.Path = path.Join(iu.Path, "auth")
			q := iu.Query()
			q.Add("client_id", cli.authConfig.ClientID)
			q.Add("redirect_uri", cli.authConfig.RedirectURI)
			q.Add("response_type", "code")
			// TODO(kenji): Remove unnecessary scopes.
			// "offline_access" for refresh token.
			q.Add("scope", "openid profile email offline_access")
			iu.RawQuery = q.Encode()
			fmt.Println("Opening browser to login...")
			if err := browser.OpenURL(iu.String()); err != nil {
				return err
			}

			l, err := net.Listen("tcp", ru.Host)
			if err != nil {
				return err
			}
			cli.listener = l
			http.HandleFunc(ru.Path, cli.handleCallback)
			if err := http.Serve(l, nil); err != nil {
				// Ignore an error if that is caused by closing the listener.
				if !strings.Contains(err.Error(), "use of closed network connection") {
					return err
				}
			}

			return nil
		},
	}
	return &cmd
}

func (c *client) stop() {
	go func() {
		time.Sleep(time.Second)
		_ = c.listener.Close()
	}()
}

func (c *client) handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("method not implemented: %s", r.Method), http.StatusNotImplemented)
		return
	}

	if errMsg := r.FormValue("error"); errMsg != "" {
		http.Error(w, fmt.Sprintf("%s: %s", errMsg, r.FormValue("error_description")), http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")
	if code == "" {
		http.Error(w, fmt.Sprintf("no code in request: %q", r.Form), http.StatusBadRequest)
		return
	}

	if err := c.obtainToken(r.Context(), code); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Successfully logged in.")

	c.stop()
}

func (c *client) obtainToken(ctx context.Context, code string) error {
	iu, err := url.Parse(c.authConfig.IssuerURL)
	if err != nil {
		return fmt.Errorf("parse issuer-url: %v", err)
	}

	dialer := &net.Dialer{}
	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		if c.issuerResolvedAddr != "" && addr == fmt.Sprintf("%s:80", iu.Host) {
			if strings.Contains(c.issuerResolvedAddr, ":") {
				addr = c.issuerResolvedAddr
			} else {
				addr = fmt.Sprintf("%s:80", c.issuerResolvedAddr)
			}
		}

		return dialer.DialContext(ctx, network, addr)
	}

	ctx = oidc.ClientContext(ctx, http.DefaultClient)
	provider, err := oidc.NewProvider(ctx, c.authConfig.IssuerURL)
	if err != nil {
		return err
	}

	oauth2Config := &oauth2.Config{
		ClientID:     c.authConfig.ClientID,
		ClientSecret: c.authConfig.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  c.authConfig.RedirectURI,
	}

	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to get token: %v", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return fmt.Errorf("no id_token in token response")
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: c.authConfig.ClientID})
	if _, err := verifier.Verify(ctx, rawIDToken); err != nil {
		return fmt.Errorf("failed to verify ID token: %v", err)
	}

	accessToken, ok := token.Extra("access_token").(string)
	if !ok {
		return fmt.Errorf("no access_token in token response")
	}

	refreshToken, ok := token.Extra("refresh_token").(string)
	if !ok {
		return fmt.Errorf("no refresh_token in token response")
	}

	if err := accesstoken.SaveToken(&accesstoken.T{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}); err != nil {
		return fmt.Errorf("failed to save token token: %v", err)
	}

	return nil

}
