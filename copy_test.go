package dd

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncorrectString(t *testing.T) {

	// size = 36 bytes
	src := strings.NewReader("abcdefghijklmnopqrstuvwxyz0123456789")

	{
		var reader io.Reader = src
		var b bytes.Buffer
		dst := bufio.NewWriter(&b)

		var writer io.Writer = dst

		copy(&reader, &writer, 0, 2, src.Size())
		assert.Equal(t, src.Size(), int64(b.Len()))
	}

	{
		src.Seek(0, 0)
		var reader io.Reader = src
		var b bytes.Buffer
		dst := bufio.NewWriter(&b)

		var writer io.Writer = dst

		copy(&reader, &writer, 10, 2, src.Size())
		assert.Equal(t, int64(10), int64(b.Len()))
	}

	{
		src.Seek(0, 0)
		var reader io.Reader = src
		var b bytes.Buffer
		dst := bufio.NewWriter(&b)

		var writer io.Writer = dst

		copy(&reader, &writer, 11, 2, src.Size())
		assert.Equal(t, int64(11), int64(b.Len()))
	}
	//assert.Equal(t, "incorrect string", "11")
}
