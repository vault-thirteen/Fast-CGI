package h

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/vault-thirteen/auxie/number"
)

const (
	LF = '\n'

	// Status is a virtual HTTP header.
	// It looks like PHP returns HTTP status information inside a virtual HTTP
	// header named 'Status'.
	Status                       = "status"
	HttpHeaderNameValueDelimiter = ":"
	Space                        = " "
)

const (
	NoHeaderOnLine     = "no header on line: %v"
	HeaderNameIsEmpty  = "header name is empty: %v"
	HeaderValueIsEmpty = "header value is empty: %v"
)

// SplitHttpHeadersFromStdout splits stdout stream into HTTP headers and HTTP
// body.
func SplitHttpHeadersFromStdout(stdout []byte) (data *Data, err error) {
	data = &Data{
		Headers: make([]*Header, 0),
	}

	rdr := bufio.NewReader(bytes.NewReader(stdout))
	var line string
	var hdr *Header

	// Read HTTP headers.
	for {
		line, err = rdr.ReadString(LF)
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)

		// HTTP headers and body are separated with an empty line.
		if len(line) == 0 {
			break
		}

		hdr, err = ParseHttpHeader(line)
		if err != nil {
			return nil, err
		}

		if strings.ToLower(hdr.Name) == Status {
			data.StatusCode, data.StatusText, err = ParsePhpHttpStatus(hdr.Value)
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

// ParseHttpHeader parses a line of text containing the HTTP header.
func ParseHttpHeader(line string) (hdr *Header, err error) {
	parts := strings.Split(line, HttpHeaderNameValueDelimiter)
	if len(parts) != 2 {
		return nil, fmt.Errorf(NoHeaderOnLine, line)
	}

	hdr = &Header{
		Name:  strings.TrimSpace(parts[0]),
		Value: strings.TrimSpace(parts[1]),
	}

	if len(hdr.Name) == 0 {
		return nil, fmt.Errorf(HeaderNameIsEmpty, line)
	}

	if len(hdr.Value) == 0 {
		return nil, fmt.Errorf(HeaderValueIsEmpty, line)
	}

	return hdr, nil
}

// ParsePhpHttpStatus parses information about HTTP status returned by PHP.
func ParsePhpHttpStatus(statusValue string) (statusCode uint, statusText string, err error) {
	n := strings.Index(statusValue, Space)
	statusCodeStr := statusValue[0:n]
	statusCode, err = number.ParseUint(statusCodeStr)
	if err != nil {
		return statusCode, statusText, err
	}

	n++
	statusText = strings.TrimSpace(statusValue[n:])
	return statusCode, statusText, nil
}
