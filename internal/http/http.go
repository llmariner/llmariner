package http

import (
	"bytes"
	"encoding/json"
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
	body, err := c.SendRequest(method, path, req)
	if err != nil {
		return err
	}

	defer func() {
		_ = body.Close()
	}()
	respBody, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("read response body: %s", err)
	}

	m := newMarshaler()
	if err := m.Unmarshal(respBody, resp); err != nil {
		return fmt.Errorf("unmarshal response: %s", err)
	}

	return nil
}

// SendRequest sends a request to the server and returns the response body.
func (c *Client) SendRequest(
	method string,
	path string,
	req any,
) (io.ReadCloser, error) {
	m := newMarshaler()

	reqBody, err := m.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %s", err)
	}

	var params map[string]interface{}
	if method == http.MethodGet {
		// Convert the body data to params as GET requests don't have a body.
		//
		// TODO(kenji): Support nested params.
		err := json.Unmarshal(reqBody, &params)
		if err != nil {
			return nil, fmt.Errorf("unmarshal request: %s", err)
		}
		reqBody = []byte{}
	}

	hreq, err := http.NewRequest(method, c.env.Config.EndpointURL+path, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %s", err)
	}

	query := hreq.URL.Query()
	for key, value := range params {
		query.Add(key, fmt.Sprintf("%v", value))
	}
	hreq.URL.RawQuery = query.Encode()

	c.addHeaders(hreq)
	hresp, err := http.DefaultClient.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("send request: %s", err)
	}

	if hresp.StatusCode != http.StatusOK {
		defer func() {
			_ = hresp.Body.Close()
		}()
		s := extractErrorMessage(hresp.Body)
		return nil, fmt.Errorf("unexpected status code: %s (message: %q)", hresp.Status, s)
	}

	return hresp.Body, nil
}

func extractErrorMessage(body io.ReadCloser) string {
	b, err := io.ReadAll(body)
	if err != nil {
		return ""
	}
	type resp struct {
		Message string `json:"message"`
	}
	var r resp
	if err := json.Unmarshal(b, &r); err != nil {
		// Return the original body if it's not JSON (error response from an non-gRPC HTTP handler).
		return string(b)
	}
	return r.Message
}

// addHeaders adds headers to the request.
func (c *Client) addHeaders(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+c.env.Token.AccessToken)
	if id := c.env.Config.Context.OrganizationID; id != "" {
		req.Header.Add("Openai-Organization", id)
	}
	if id := c.env.Config.Context.ProjectID; id != "" {
		req.Header.Add("Openai-Project", id)
	}

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
