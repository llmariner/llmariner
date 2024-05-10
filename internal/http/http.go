package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	iruntime "github.com/llm-operator/cli/internal/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

// NewClient creates a new HTTP client.
func NewClient(env *iruntime.Env) *Client {
	return &Client{
		env: env,
	}
}

// Client is an HTTP client.
type Client struct {
	env *iruntime.Env
}

// Send sends a request to the server.
//
// We use this client instead of using gRPC as we don't know if an ingress controller in a customer's
// environment supports gRPC.
func (c *Client) Send(
	method string,
	path string,
	req any,
	resp any,
) error {
	m := newMarshaler()

	reqBody, err := m.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %s", err)
	}

	hreq, err := http.NewRequest(method, c.env.Config.EndpointURL+path, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("create request: %s", err)
	}

	c.addHeaders(hreq)
	hresp, err := http.DefaultClient.Do(hreq)
	if err != nil {
		return fmt.Errorf("send request: %s", err)
	}
	if hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %s", hresp.Status)
	}

	defer func() {
		_ = hresp.Body.Close()
	}()
	respBody, err := io.ReadAll(hresp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %s", err)
	}

	if err := m.Unmarshal(respBody, resp); err != nil {
		return fmt.Errorf("unmarshal response: %s", err)
	}

	return nil
}

// addHeaders adds headers to the request.
func (c *Client) addHeaders(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+c.env.Token.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
}

func newMarshaler() *runtime.JSONPb {
	return &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
}
