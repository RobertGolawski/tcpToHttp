package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersValidSpacing(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", (headers)["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}

func TestHeadersInvalidSpacing(t *testing.T) {
	headers := NewHeaders()
	data := []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069                           \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 57, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = map[string]string{"host": "localhost:42069"}
	data = []byte("User-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	assert.Equal(t, 25, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n a bunch of other stuff")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Empty(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Same header key
	headers = map[string]string{"host": "localhost:8000"}
	data = []byte("Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:8000, localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}

//	func TestMultipleValues(t *testing.T) {
//		headers := NewHeaders()
//		data := []byte("Set-Person: lane-loves-go\r\n")
//		data1 := []byte("Set-Person: prime-loves-zig\r\n")
//		data2 := []byte("Set-Person: tj-loves-ocaml\r\n\r\n")
//		n, done, err := headers.Parse(data)
//		n1, done1, err1 := headers.Parse(data1)
//		n2, done2, err2 := headers.Parse(data2)
//		require.NoError(t, err)
//		require.NoError(t, err1)
//		require.NoError(t, err2)
//		assert.Equal(t, 27, n)
//		assert.Equal(t, 29, n1)
//		assert.Equal(t, 28, n2)
//		assert.False(t, done)
//		assert.False(t, done1)
//		assert.False(t, done2)
//	}
func TestHeadersInvalidCharsParse(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		expectedErr string
	}{
		{
			name:        "Header name with space",
			data:        "Host Name: localhost:42069\r\n\r\n",
			expectedErr: "invalid character",
		},
		{
			name:        "Header name with disallowed special char (comma)",
			data:        "User,Agent: curl/7.81.0\r\n\r\n",
			expectedErr: "invalid character",
		},
		{
			name:        "Header name with disallowed special char (open parenthesis)",
			data:        "(Custom-Header): somevalue\r\n\r\n",
			expectedErr: "invalid character",
		},
		{
			name:        "Header name with disallowed special char (slash)",
			data:        "Content/Type: application/json\r\n\r\n",
			expectedErr: "invalid character",
		},
		{
			name:        "Header name with null byte",
			data:        "Null\x00Header: value\r\n\r\n",
			expectedErr: "invalid character",
		},
		{
			name:        "Header name with non-ASCII char",
			data:        "Héader: value\r\n\r\n",
			expectedErr: "invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := NewHeaders()
			n, done, err := headers.Parse([]byte(tt.data))

			require.Error(t, err, "Expected an error for test case: %s", tt.name)
			assert.Contains(t, err.Error(), tt.expectedErr, "Error message mismatch for test case: %s", tt.name)
			assert.Equal(t, 0, n, "Expected 0 bytes consumed on error for test case: %s", tt.name)
			assert.False(t, done, "Expected not done on error for test case: %s", tt.name)
		})
	}

	// Test case for malformed header (missing value after colon)
	t.Run("Malformed header missing value", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host:\r\n\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, "", headers["host"])
		assert.Equal(t, 7, n)
		assert.False(t, done)
	})
}
