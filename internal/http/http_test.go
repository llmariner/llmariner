package http

import (
	"testing"

	uv1 "github.com/llm-operator/user-manager/api/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestUnmarshal(t *testing.T) {
	b := []byte(`{
  "object": "list",
  "data": [
     {
       "id": "6a1d418e-6186-4865-8117-47aac1de4714",
       "object": "user.api_key",
       "name": "test",
       "created_at": "1714803213"
     }
  ]
}`)

	m := newMarshaler()
	var got uv1.ListAPIKeysResponse
	err := m.Unmarshal(b, &got)
	assert.NoError(t, err)
	want := uv1.ListAPIKeysResponse{
		Object: "list",
		Data: []*uv1.APIKey{
			{
				Id:        "6a1d418e-6186-4865-8117-47aac1de4714",
				Object:    "user.api_key",
				Name:      "test",
				CreatedAt: 1714803213,
			},
		},
	}
	assert.True(t, proto.Equal(&want, &got))
}
