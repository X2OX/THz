package cors

import (
	"go.x2ox.com/THz"
)

// Config
//
// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#the_http_response_headers
type Config struct {
	Skip             func(*THz.Context) bool // Optional. Default return false
	AllowOrigins     string                  // Optional. Default value "*"
	AllowMethods     string                  // Optional. Default value "GET,POST,HEAD,PUT,DELETE,PATCH"
	AllowHeaders     string                  // Optional. Default value "".
	AllowCredentials bool                    // Optional. Default value false.
	ExposeHeaders    string                  // Optional. Default value "".
	MaxAge           int                     // Optional. Default value 0.
}
