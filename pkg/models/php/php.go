package pm

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/vault-thirteen/Fast-CGI/pkg/models/common"
	"github.com/vault-thirteen/auxie/number"
	"github.com/vault-thirteen/header"
)

const (
	HeaderNameValueDelimiter = ":"

	// Status is a virtual HTTP header. It looks like PHP returns HTTP status
	// information inside a virtual HTTP header named 'Status'.
	Status = "status"

	OldRelativeUrlMarker = `./`
	ForwardSlash         = `/`
)

const (
	ErrNoHeaderOnLine            = "no header on line: %v"
	ErrHeaderNameIsEmpty         = "header name is empty: %v"
	ErrHeaderValueIsEmpty        = "header value is empty: %v"
	ErrTooManyRelativeUrlMarkers = "too many relative URL markers: %v"
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
	sepIdx := strings.Index(line, HeaderNameValueDelimiter)
	if sepIdx < 0 {
		return nil, fmt.Errorf(ErrNoHeaderOnLine, line)
	}

	hdr = &Header{
		Name:  strings.TrimSpace(line[:sepIdx]),
		Value: strings.TrimSpace(line[sepIdx+1:]),
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

// FixLocationHeader fixes relative URLs in 'Location' HTTP headers.
// 'currentUrlPath' is the value of 'URL.Path' of the current request.
func (dta *Data) FixLocationHeader(currentUrlPath string) (err error) {
	currentUrlPathStripped := strings.TrimSuffix(strings.TrimPrefix(currentUrlPath, ForwardSlash), ForwardSlash)

	for _, hdr := range dta.Headers {
		if hdr.Name == header.HttpHeaderLocation {
			hdr.Value, err = fixRelativeUrl(hdr.Value, currentUrlPathStripped)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// fixRelativeUrl fixes the relative URL specified in 'url' parameter.
func fixRelativeUrl(url string, currentUrlPathStripped string) (fixedUrl string, err error) {
	markersCount := strings.Count(url, OldRelativeUrlMarker)
	if markersCount == 0 {
		return url, nil
	}
	if markersCount > 1 {
		return "", fmt.Errorf(ErrTooManyRelativeUrlMarkers, url)
	}

	return ForwardSlash + path.Join(currentUrlPathStripped, url[strings.Index(url, OldRelativeUrlMarker)+2:]), nil
}
