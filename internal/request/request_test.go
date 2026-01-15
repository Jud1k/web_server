package request_test

import (
	"io"
	"testing"
	"request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}

func TestRequestLine_AllValidMethod(t *testing.T){
	methods := []string{
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"OPTIONS",
		"HEAD",
		"TRACE",
		"CONNECT",
	}

	for _,method := range(methods){
		t.Run(method,func(t *testing.T) {
				reader := &chunkReader{
			data: method + " / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 3,
			}
			r,err := request.RequestFromReader(reader)
			
			require.NoError(t,err)
			require.NotNil(t,r)
			assert.Equal(t,method,r.RequestLine.Method)
		})
		 
	}
}

func TestRequestLine_InvalidMethods(t *testing.T) {
	methods := []string{
		"get",
		"Post",
		"FETCH",
		"",
		"123",
	}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			reader := &chunkReader{
			data:  method + " / HTTP/1.1\r\nHost: localhost\r\n\r\n",
			numBytesPerRead: 3,
			}
			_, err := request.RequestFromReader(reader)
			require.Error(t, err)
		})
	}
}


func TestGetRequestLineWithPath(t *testing.T){
	reader := &chunkReader{data: "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",numBytesPerRead: 3}
	r, err := request.RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

}

func TestGetRequestLineWrongMethodError(t *testing.T) {
	reader := &chunkReader{data: "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",numBytesPerRead: 3}

	_, err := request.RequestFromReader(reader)
	require.Error(t, err)
}

func TestGetRequestLineWrongVersionError(t *testing.T) {
	reader := &chunkReader{data: "GET /coffee HTTP/2.0\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",numBytesPerRead: 3}

	_, err := request.RequestFromReader(reader)
	require.Error(t, err)
}



