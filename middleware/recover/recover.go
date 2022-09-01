package recover

import (
	"go.uber.org/zap"
	thz "go.x2ox.com/THz"
	"go.x2ox.com/sorbifolia/coarsetime"
)

type Config struct {
	TraceStack bool
}

func New() *Config {
	return &Config{TraceStack: false}
}

func (c *Config) Middleware() thz.Handler {
	return func(c *thz.Context) {
		defer func() {
			if e := recover(); e != nil {
				c.L().Error("Recover Panic",
					zap.Time("time", coarsetime.FloorTime()),
					zap.Any("error", e),
					zap.StackSkip("stack", 1),
				)
				c.Abort().Status(500)
			}
		}()

		c.Next()
	}
}
