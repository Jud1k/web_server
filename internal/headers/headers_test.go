package headers_test

import (
	"testing"

	"github.com/Jud1k/web_server/internal/headers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeadersSuccess(t *testing.T) {
	headers := headers.Headers{}
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}

func TestParseHeadersInvalidSpaces(t *testing.T) {
	headers := headers.Headers{}
	data := []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestParseHeadersWithOneExisting(t *testing.T) {
	headers := headers.Headers{}
	data := []byte("Host: localhost:42069\r\nHost: localhost:12345\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, len(data[:n]), n)
	assert.False(t, done)

	_, _, err = headers.Parse(data[n:])
	assert.Equal(t, "localhost:42069, localhost:12345", headers["host"])
	assert.Equal(t, len(data[n:])-2, n)
	assert.False(t, done)
}

func TestParseHeadersInvalidCharacters(t *testing.T) {
	headers := headers.Headers{}
	data := []byte("H@st: localhost:42069\r\n\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestParseHeadersNotCrlfToBody(t *testing.T) {
	headers := headers.Headers{}
	data := []byte("Host: localhost:42069\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, len(data), n)
	assert.False(t, done)
}
