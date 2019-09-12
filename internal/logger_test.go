package internal

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger_Log(t *testing.T) {
	var buf bytes.Buffer
	l := Logger{
		Verbosity: 2,
		Logger:    log.New(&buf, "", 0),
	}
	l.Log(3, "hello")
	require.Equal(t, "", buf.String())
	l.Log(1, "hello")
	require.NotEqual(t, "", buf.String())
}
