package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/cli/browser"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/llm-operator/cli/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

type client struct {
	authConfig config.Auth

	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	listener net.Listener
}

func loginCmd() *cobra.Command {
	var (
		cli                client
		issuerResolvedAddr string
	)
	cmd := cobra.Command{
		Use: "login",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
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
				return fmt.Errorf("parse issuer-uri: %v", err)
			}

			dialer := &net.Dialer{}
			http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				if issuerResolvedAddr != "" && addr == fmt.Sprintf("%s:80", iu.Host) {
					addr = fmt.Sprintf("%s:80", issuerResolvedAddr)
				}
				return dialer.DialContext(ctx, network, addr)
			}

			ctx := oidc.ClientContext(context.Background(), http.DefaultClient)
			provider, err := oidc.NewProvider(ctx, c.Auth.IssuerURL)
			if err != nil {
				return err
			}
			cli.provider = provider
			cli.verifier = provider.Verifier(&oidc.Config{ClientID: cli.authConfig.ClientID})

			if issuerResolvedAddr != "" {
				iu.Host = issuerResolvedAddr
			}
			iu.Path = path.Join(iu.Path, "auth")
			q := iu.Query()
			q.Add("client_id", cli.authConfig.ClientID)
			q.Add("redirect_uri", cli.authConfig.RedirectURI)
			q.Add("response_type", "code")
			// TODO(kenji): Remove unnecessary scopes.
			q.Add("scope", "openid profile email")
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
				return err
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&issuerResolvedAddr, "issuer-resolved-addr", "", "Address of the issuer. Set when the issuer URL is not resolvable from the client.")
	return &cmd
}

func (c *client) stop() {
	go func() {
		time.Sleep(time.Second)
		_ = c.listener.Close()
	}()
}

func (c *client) handleCallback(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		token *oauth2.Token
	)

	ctx := oidc.ClientContext(r.Context(), http.DefaultClient)
	oauth2Config := &oauth2.Config{
		ClientID:     c.authConfig.ClientID,
		ClientSecret: c.authConfig.ClientSecret,
		Endpoint:     c.provider.Endpoint(),
		RedirectURL:  c.authConfig.RedirectURI,
	}

	if r.Method == http.MethodGet {
		if errMsg := r.FormValue("error"); errMsg != "" {
			http.Error(w, fmt.Sprintf("%s: %s", errMsg, r.FormValue("error_description")), http.StatusBadRequest)
			return
		}
		code := r.FormValue("code")
		if code == "" {
			http.Error(w, fmt.Sprintf("no code in request: %q", r.Form), http.StatusBadRequest)
			return
		}
		token, err = oauth2Config.Exchange(ctx, code)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get token: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, fmt.Sprintf("method not implemented: %s", r.Method), http.StatusBadRequest)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "no id_token in token response", http.StatusInternalServerError)
		return
	}
	if _, err := c.verifier.Verify(r.Context(), rawIDToken); err != nil {
		http.Error(w, fmt.Sprintf("failed to verify ID token: %v", err), http.StatusInternalServerError)
		return
	}

	accessToken, ok := token.Extra("access_token").(string)
	if !ok {
		http.Error(w, "no access_token in token response", http.StatusInternalServerError)
		return
	}

	fmt.Println("Successfully logged in.")
	fmt.Println("Token:", accessToken)
	c.stop()
}
