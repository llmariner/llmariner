package auth

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cli/browser"

	"github.com/llm-operator/cli/internal/accesstoken"
	"github.com/llm-operator/cli/internal/config"
	"github.com/spf13/cobra"
)

type client struct {
	tokenExechanger *accesstoken.TokenExchanger

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

			c, err := config.LoadOrCreate()
			if err != nil {
				return fmt.Errorf("load config: %s", err)
			}

			tokenExchanger, err := accesstoken.NewTokenExchanger(c)
			if err != nil {
				return fmt.Errorf("create token exchanger: %v", err)
			}
			cli.tokenExechanger = tokenExchanger

			loginURL, err := tokenExchanger.LoginURL()
			if err != nil {
				return fmt.Errorf("get login URL: %v", err)
			}

			fmt.Println("Opening browser to login...")
			if err := browser.OpenURL(loginURL); err != nil {
				return err
			}

			ru, err := url.Parse(c.Auth.RedirectURI)
			if err != nil {
				return fmt.Errorf("parse redirect-uri: %v", err)
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

	if err := c.tokenExechanger.ObtainToken(r.Context(), code); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Successfully logged in.")

	c.stop()
}
