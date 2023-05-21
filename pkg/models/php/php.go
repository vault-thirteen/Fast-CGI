package pm

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/vault-thirteen/Fast-CGI/pkg/models/common"
	"github.com/vault-thirteen/auxie/number"
)

const (
	HeaderNameValueDelimiter = ":"

	// Status is a virtual HTTP header. It looks like PHP returns HTTP status
	// information inside a virtual HTTP header named 'Status'.
	Status = "status"
)

const (
	ErrNoHeaderOnLine     = "no header on line: %v"
	ErrHeaderNameIsEmpty  = "header name is empty: %v"
	ErrHeaderValueIsEmpty = "header value is empty: %v"
)

// Data is data returned by a PHP script.
type Data struct {
	StatusCode uint
	StatusText string
	Headers    []*Header
	Body       []byte
}

// Header is an HTTP header returned by a PHP script.
type Header struct {
	Name  string
	Value string
}

// SplitHeadersFromStdout splits PHP stdout stream into HTTP headers and HTTP
// body.
func SplitHeadersFromStdout(stdout []byte) (data *Data, err error) {
	data = &Data{
		Headers: make([]*Header, 0),
	}

	rdr := bufio.NewReader(bytes.NewReader(stdout))
	var line string
	var hdr *Header

	// Read HTTP headers.
	for {
		line, err = rdr.ReadString(cm.LF)
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)

		// HTTP headers and body are separated with an empty line.
		if len(line) == 0 {
			break
		}

		hdr, err = ParseHeader(line)
		if err != nil {
			return nil, err
		}

		if strings.ToLower(hdr.Name) == Status {
			data.StatusCode, data.StatusText, err = ParseStatus(hdr.Value)
			if err != nil {
				return nil, err
			}
		} else {
			data.Headers = append(data.Headers, hdr)
		}
	}

	// It looks like PHP does not set the HTTP status when by default.
	// If the script have not set the status code, it will be zero.
	// We do not change this behaviour while some PHP scripts may be dependent
	// on it.

	// Read HTTP body.
	data.Body, err = io.ReadAll(rdr)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ParseHeader parses a PHP output line of text containing the HTTP header.
func ParseHeader(line string) (hdr *Header, err error) {
	parts := strings.Split(line, HeaderNameValueDelimiter)
	if len(parts) != 2 {
		return nil, fmt.Errorf(ErrNoHeaderOnLine, line)
	}

	hdr = &Header{
		Name:  strings.TrimSpace(parts[0]),
		Value: strings.TrimSpace(parts[1]),
	}

	if len(hdr.Name) == 0 {
		return nil, fmt.Errorf(ErrHeaderNameIsEmpty, line)
	}

	if len(hdr.Value) == 0 {
		return nil, fmt.Errorf(ErrHeaderValueIsEmpty, line)
	}

	return hdr, nil
}

// ParseStatus parses information about HTTP status returned by PHP.
func ParseStatus(statusValue string) (statusCode uint, statusText string, err error) {
	n := strings.Index(statusValue, cm.Space)
	statusCodeStr := statusValue[0:n]
	statusCode, err = number.ParseUint(statusCodeStr)
	if err != nil {
		return statusCode, statusText, err
	}

	n++
	statusText = strings.TrimSpace(statusValue[n:])
	return statusCode, statusText, nil
}
