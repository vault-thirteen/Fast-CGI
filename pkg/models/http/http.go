package hm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vault-thirteen/Fast-CGI/pkg/models/NameValuePair"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/common"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/data"
)

const (
	ErrAuthorizationSyntax = "syntax error in authorization header: %v"
)

// ParseAuthorizationHeader parses the 'Authorization' HTTP header.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Authorization.
func ParseAuthorizationHeader(header string) (scheme string, parameters string, err error) {
	if len(header) == 0 {
		return scheme, parameters, nil
	}

	parts := strings.Split(header, cm.Space)
	if len(parts) != 2 {
		return scheme, parameters, fmt.Errorf(ErrAuthorizationSyntax, header)
	}

	return parts[0], parts[1], nil
}

// AddHttpHeadersToParameters adds HTTP headers to FastCGI parameters.
func AddHttpHeadersToParameters(parameters *[]*nvpair.NameValuePair, headers http.Header) {
	if parameters == nil {
		return
	}

	var parameter *nvpair.NameValuePair
	for hdrName, hdrValues := range headers {
		for _, hdrValue := range hdrValues {
			parameter = nvpair.NewNameValuePairWithTextValueU(ComposeCgiParameterNameFromHttpHeader(hdrName), hdrValue) // 4.1.18.
			if parameter != nil {
				*parameters = append(*parameters, parameter)
			}
		}
	}
}

// ComposeCgiParameterNameFromHttpHeader composes CGI parameter name using the
// specified HTTP header name. See the 4.1.18 section of the CGI 1.1 interface
// specification.
func ComposeCgiParameterNameFromHttpHeader(headerName string) (cgiParamName string) {
	return dm.ParameterPrefix_Http + strings.ToUpper(strings.ReplaceAll(headerName, "-", "_"))
}
