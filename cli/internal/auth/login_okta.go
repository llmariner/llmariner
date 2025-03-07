package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/cli/browser"
	"github.com/llmariner/llmariner/cli/internal/accesstoken"
	"github.com/llmariner/llmariner/cli/internal/configs"
	llmcontext "github.com/llmariner/llmariner/cli/internal/context"
	"golang.org/x/oauth2"
)

func oktaLogin(ctx context.Context, c *configs.C, noOpen bool) error {
	// port within 8080 - 8084 is allowed. port should match with the one used in redirect url specified in the config.
	port := int(8084)
	listerner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("listen: %s", err)
	}

	cli, err := newOktaClient(ctx, c.Auth)
	if err != nil {
		return fmt.Errorf("new okta client: %s", err)
	}

	tokenExchanger, err := accesstoken.NewOktaTokenExchanger(cli.config, cli.codeVerifier)
	if err != nil {
		return fmt.Errorf("create token exchanger: %s", err)
	}
	cli.tokenExechanger = tokenExchanger

	cli.wg.Add(1)
	go cli.start(listerner)

	address := fmt.Sprintf("localhost:%d", port)
	loginURL := fmt.Sprintf("http://%s/login", address)

	if noOpen {
		fmt.Printf("Please open the following URL from your browser:\n%s\n", loginURL)
	} else {
		go func() {
			if err := awaitServerReady(address); err != nil {
				printError(err, loginURL)
				return
			}

			fmt.Println("Opening browser to login...")
			if err := browser.OpenURL(loginURL); err != nil {
				printError(err, loginURL)
				return
			}
		}()
	}

	cli.wg.Wait()

	fmt.Println("\nSetting the context...")
	if err := llmcontext.Set(ctx); err != nil {
		return err
	}

	return nil
}

type oktaClient struct {
	tokenExechanger *accesstoken.OktaTokenExchanger

	listener net.Listener

	state string
	// The code verifier is a cryptographically random
	// string using the characters A-Z, a-z, 0-9, and the
	// punctuation characters -._~ between 43 and 128 characters long
	// (https://www.oauth.com/oauth2-servers/pkce/authorization-request/).
	codeVerifier string

	config *oauth2.Config

	wg sync.WaitGroup
}

func newOktaClient(ctx context.Context, auth configs.Auth) (*oktaClient, error) {
	oauth2Config, err := accesstoken.NewOauth2Config(ctx, auth)
	if err != nil {
		return nil, fmt.Errorf("new oauth2 config: %s", err)
	}
	oauth2Config.Scopes = []string{"openid", "profile", "email", "offline_access"}

	return &oktaClient{
		config:       oauth2Config,
		state:        getRandomString(64),
		codeVerifier: getRandomString(64),
	}, nil
}

func (c *oktaClient) start(l net.Listener) {
	http.HandleFunc("/login", c.handleLogin)
	http.HandleFunc("/callback", c.handleCallback)
	if err := http.Serve(l, nil); err != nil {
		fmt.Printf("HTTP server for login finished with an error: %s", err)
	}
}

func (c *oktaClient) stop() {
	go func() {
		time.Sleep(time.Second)
		_ = c.listener.Close()
	}()
}

func (c *oktaClient) handleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != c.state {
		fmt.Printf("invalid oauth state\n")
		return
	}

	code := r.FormValue("code")
	if err := c.tokenExechanger.ObtainToken(r.Context(), code); err != nil {
		fmt.Printf("code exchange failed: %s\n", err)
		return
	}

	success := `<html>
<body>
  Successfully logged into CloudNatix.
</body>
</html>`
	if _, err := fmt.Fprint(w, success); err != nil {
		fmt.Printf("failed to print: %s\n", err)
	}

	fmt.Println("Successfully logged in.")

	c.stop()
	c.wg.Done()
}

func (c *oktaClient) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Generate code challenge.
	sha := sha256.Sum256([]byte(c.codeVerifier))
	codeChallenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(sha[:])
	url := c.config.AuthCodeURL(
		c.state,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// printError prints a warning to the console and a prompt to open the browser
// manually.
func printError(err error, url string) {
	fmt.Printf("Unable to open browser: %s\n. Please access %s manually\n", err, url)
}

// awaitServerReady waits for the HTTP server to become ready. Any error
// encountered while waiting, or upon timeout is returned.
func awaitServerReady(address string) error {
	const maxAttempts = 5
	var err error
	var count int
	for {
		if count == maxAttempts {
			return fmt.Errorf("login: timed out waiting for server: %s", err)
		}

		conn, err := net.Dial("tcp", address)
		if err == nil {
			return conn.Close()
		}

		count++
		time.Sleep(time.Second)
	}
}

func getRandomString(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
