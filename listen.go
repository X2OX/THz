package THz

import (
	"net"

	"go.uber.org/zap/zapcore"
)

func (thz *THz) init() *THz {
	thz.route.Sort()
	thz.IRouter = nil // not allowed updates routing

	if thz.log == nil {
		thz.SetZapLog(zapcore.DebugLevel)
	}

	return thz
}

func (thz *THz) Stop() error               { return thz.srv.Shutdown() }
func (thz *THz) Run(ln net.Listener) error { return thz.init().srv.Serve(ln) }
func (thz *THz) RunTLS(ln net.Listener, cert, key string) error {
	return thz.init().srv.ServeTLS(ln, cert, key)
}
func (thz *THz) ListenAndServe(addr string) error {
	return thz.init().srv.ListenAndServe(addr)
}
