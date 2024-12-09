package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToAge(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tcs := []struct {
		t    time.Time
		want string
	}{
		{
			t:    now.Add(-1 * time.Second),
			want: "1s",
		},
		{
			t:    now.Add(-60 * time.Second),
			want: "1m",
		},
		{
			t:    now.Add(-65 * time.Second),
			want: "1m",
		},
		{
			t:    now.Add(-4*time.Hour - 3*time.Minute),
			want: "4h3m",
		},
		{
			t:    now.Add(-6 * time.Hour),
			want: "6h",
		},
		{
			t:    now.Add(-26 * time.Hour),
			want: "1d",
		},
		{
			t:    now.Add(-24 * 4 * time.Hour),
			want: "4d",
		},
		{
			t:    now.Add(-24 * 20 * time.Hour),
			want: "20d",
		},
		{
			t:    now.Add(-24 * 400 * time.Hour),
			want: "1y",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.t.String(), func(t *testing.T) {
			got := toAge(tc.t, now)
			assert.Equal(t, tc.want, got)
		})
	}
}
