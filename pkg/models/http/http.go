package http

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

const LF = '\n'

const (
	NoHeaderOnLine     = "no header on line: %v"
	HeaderNameIsEmpty  = "header name is empty: %v"
	HeaderValueIsEmpty = "header value is empty: %v"
)

// SplitHttpHeadersFromStdout splits stdout stream into HTTP headers and HTTP
// body.
func SplitHttpHeadersFromStdout(stdout []byte) (headers []*Header, body []byte, err error) {
	rdr := bufio.NewReader(bytes.NewReader(stdout))
	var line string
	var hdr *Header

	// Read HTTP headers.
	for {
		line, err = rdr.ReadString(LF)
		if err != nil {
			return nil, nil, err
		}

		line = strings.TrimSpace(line)

		// HTTP headers and body are separated with an empty line.
		if len(line) == 0 {
			break
		}

		hdr, err = ParseHttpHeader(line)
		if err != nil {
			return nil, nil, err
		}

		headers = append(headers, hdr)
	}

	// Read HTTP body.
	body, err = io.ReadAll(rdr)
	if err != nil {
		return nil, nil, err
	}

	return headers, body, nil
}

// ParseHttpHeader parses a line of text containing the HTTP header.
func ParseHttpHeader(line string) (hdr *Header, err error) {
	parts := strings.Split(line, ":")
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
