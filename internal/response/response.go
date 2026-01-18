package response

import (
	"bytes"
	"fmt"
	"io"
	"maps"
	"strconv"

	"github.com/Jud1k/web_server/internal/headers"
)

type StatusCode int

const (
	StatusCodeOk            StatusCode = 200
	StatusCodeBadRequest    StatusCode = 400
	StatusCodeInternalError StatusCode = 500
)

type writerState int

const (
	stateInitial writerState = iota
	stateStatusWritten
	stateHeadersWritten
	stateBodyWritten
)

type Writer struct {
	state   writerState
	buff    bytes.Buffer
	headers headers.Headers
}

func NewWriter() *Writer {
	return &Writer{
		state:   stateInitial,
		headers: make(headers.Headers),
	}
}

func (w *Writer) WriteTo(writer io.Writer) (n int64, err error) {
	return w.buff.WriteTo(writer)
}

func (w *Writer) Bytes() []byte {
	return w.buff.Bytes()
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != stateInitial {
		return fmt.Errorf("error: cannot write status line already in state %d", w.state)
	}
	statusMessage := getStatusMessage(statusCode)
	fmt.Fprintf(&w.buff, "HTTP/1.1 %d %s\r\n", statusCode, statusMessage)
	w.state = stateStatusWritten
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != stateStatusWritten {
		return fmt.Errorf("error: cannot write headers already in state %d", w.state)
	}
	w.headers = headers
	w.state = stateHeadersWritten
	return writeHeaders(&w.buff, headers)
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != stateHeadersWritten {
		return 0, fmt.Errorf("error: cannot write body already in state %d", w.state)
	}
	n, err := w.buff.Write(p)
	if err != nil {
		return 0, err
	}
	w.state = stateBodyWritten
	return n, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != stateHeadersWritten {
		return 0, fmt.Errorf("error: cannot write body already in state %d", w.state)
	}
	hexLen := strconv.FormatInt(int64(len(p)), 16)
	chunk := []byte(hexLen + "\r\n")
	chunk = append(chunk, p...)
	chunk = append(chunk, []byte("\r\n")...)
	return w.buff.Write(chunk)
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	w.state = stateBodyWritten
	return w.buff.Write([]byte("0\r\n\r\n"))
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.state != stateBodyWritten {
		return fmt.Errorf("error: cannot write trailers already in state %d", w.state)
	}
	maps.Copy(w.headers, h)
	return writeHeaders(&w.buff, h)
}

func getStatusMessage(statusCode StatusCode) string {
	reason := ""
	switch statusCode {
	case StatusCodeOk:
		reason = "OK"
	case StatusCodeBadRequest:
		reason = "Bad Request"
	case StatusCodeInternalError:
		reason = "Internal Server Error"
	}
	return reason
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.Headers{}
	h["Content-Length"] = fmt.Sprint(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func writeHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", key, val)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(w, "\r\n")
	return err
}
