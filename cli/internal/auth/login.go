package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cli/browser"

	"github.com/llmariner/llmariner/cli/internal/accesstoken"
	"github.com/llmariner/llmariner/cli/internal/configs"
	llmcontext "github.com/llmariner/llmariner/cli/internal/context"
	"github.com/spf13/cobra"
)

func loginCmd() *cobra.Command {
	var noOpen bool
	cmd := cobra.Command{
		Use:   "login",
		Short: "Login to LLM service",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := configs.LoadOrCreate()
			if err != nil {
				return fmt.Errorf("load config: %s", err)
			}

			if c.EnableOkta {
				return oktaLogin(cmd.Context(), c, noOpen)
			}

			return dexLogin(cmd.Context(), c, noOpen)
		},
	}
	cmd.Flags().BoolVar(&noOpen, "no-open", false, "Do not open the browser")
	return &cmd
}

func dexLogin(ctx context.Context, c *configs.C, noOpen bool) error {
	cli := newDexClient()

	tokenExchanger, err := accesstoken.NewDexTokenExchanger(c)
	if err != nil {
		return fmt.Errorf("create token exchanger: %s", err)
	}
	cli.tokenExechanger = tokenExchanger

	loginURL, err := tokenExchanger.LoginURL()
	if err != nil {
		return fmt.Errorf("get login URL: %s", err)
	}

	if noOpen {
		fmt.Printf("Please open the following URL from your browser:\n%s\n", loginURL)
	} else {
		fmt.Println("Opening browser to login...")
		if err := browser.OpenURL(loginURL); err != nil {
			return err
		}
	}

	ru, err := url.Parse(c.Auth.RedirectURI)
	if err != nil {
		return fmt.Errorf("parse redirect-uri: %s", err)
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

	fmt.Println("\nSetting the context...")
	if err := llmcontext.Set(ctx); err != nil {
		return err
	}

	return nil
}

type dexClient struct {
	tokenExechanger *accesstoken.DexTokenExchanger

	listener net.Listener
}

func newDexClient() *dexClient {
	return &dexClient{}
}

func (c *dexClient) stop() {
	go func() {
		time.Sleep(time.Second)
		_ = c.listener.Close()
	}()
}

func (c *dexClient) handleCallback(w http.ResponseWriter, r *http.Request) {
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
