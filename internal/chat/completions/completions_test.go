package completions

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	input := `line0

line1

line2

line3

`
	scanner := bufio.NewScanner(bytes.NewReader([]byte(input)))
	scanner.Buffer(make([]byte, 4096), 4096)
	scanner.Split(split)
	var got []string
	for scanner.Scan() {
		got = append(got, scanner.Text())
	}
	err := scanner.Err()
	assert.NoError(t, err)

	want := []string{"line0", "line1", "line2", "line3"}
	assert.ElementsMatch(t, want, got)
}
