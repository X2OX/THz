package cors

import (
	"strings"

	"go.x2ox.com/THz"
	"go.x2ox.com/sorbifolia/pyrokinesis"
	"go.x2ox.com/sorbifolia/strong"
)

const (
	headerAllowOrigin      = "Access-Control-Allow-Origin"
	headerAllowCredentials = "Access-Control-Allow-Credentials"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerAllowMethods     = "Access-Control-Allow-Methods"
	headerExposeHeaders    = "Access-Control-Expose-Headers"
	headerMaxAge           = "Access-Control-Max-Age"

	headerOrigin = "Origin"
)

func New(allowOrigin ...string) *Config { return &Config{AllowOrigins: normalize(allowOrigin)} }

// Config
//
// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#the_http_response_headers
type Config struct {
	Skip             func(*THz.Context) bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int
}

func (cfg *Config) Middleware() THz.Handler {
	var (
		allowOrigins     = "*"
		allowMethods     = "GET,POST,HEAD,PUT,DELETE,PATCH"
		allowHeaders     = "Origin,Content-Length,Content-Type,Authorization"
		allowCredentials = "ture"
		exposeHeaders    = ""
	)
	if len(cfg.AllowOrigins) != 0 {
		allowOrigins = strings.Join(cfg.AllowOrigins, ",")
	}
	if len(cfg.AllowMethods) != 0 {
		allowOrigins = strings.Join(cfg.AllowMethods, ",")
	}
	if len(cfg.AllowHeaders) != 0 {
		allowHeaders = strings.Join(cfg.AllowHeaders, ",")
	}
	if len(cfg.ExposeHeaders) != 0 {
		exposeHeaders = strings.Join(cfg.ExposeHeaders, ",")
	}
	maxAge := strong.Format(cfg.MaxAge)

	return func(c *THz.Context) {
		if cfg.Skip != nil && cfg.Skip(c) {
			return
		}
		if cfg.AllowCredentials {
			c.SetHeader(headerAllowCredentials, allowCredentials)
		}
		if exposeHeaders != "" {
			c.SetHeader(headerExposeHeaders, exposeHeaders)
		}

		origin := "*"
		if allowOrigins != "*" {
			has := false
			origin = pyrokinesis.Bytes.ToString(c.Header(headerOrigin))
			for _, v := range cfg.AllowOrigins {
				if v == origin {
					has = true
					break
				}
			}
			if !has {
				c.SetHeader("Vary", "Origin")
				return
			}
		}

		if !c.IsOptions() {
			c.SetHeader(headerAllowOrigin, origin)
			return // is not CORS request
		}

		c.SetHeader("Vary", "Origin,Access-Control-Allow-Headers,Access-Control-Allow-Methods")

		c.SetHeader(headerAllowOrigin, origin)
		c.SetHeader(headerAllowMethods, allowMethods)
		c.SetHeader(headerAllowHeaders, allowHeaders)
		if cfg.MaxAge != 0 {
			c.SetHeader(headerMaxAge, maxAge)
		}

		c.Status(204).Abort()
	}
}

// normalize is used to format input
func normalize(values []string) []string {
	if values == nil {
		return nil
	}

	distinct := make(map[string]bool, len(values))
	normalized := make([]string, 0, len(values))

	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if _, ok := distinct[value]; !ok {
			normalized = append(normalized, value)
			distinct[value] = true
		}
	}

	return normalized
}
