package THz

import (
	"net"

	"go.uber.org/zap"
)

func (thz *THz) init() *THz {
	thz.route.Sort()
	thz.IRouter = nil // not allowed updates routing

	if thz.log == nil {
		thz.log, _ = zap.NewProduction()
	}

	return thz
}

func (thz *THz) Stop() error               { return thz.init().srv.Shutdown() }
func (thz *THz) Run(ln net.Listener) error { return thz.init().srv.Serve(ln) }
func (thz *THz) RunTLS(ln net.Listener, cert, key string) error {
	return thz.init().srv.ServeTLS(ln, cert, key)
}
func (thz *THz) ListenAndServe(addr string) error {
	return thz.init().srv.ListenAndServe(addr)
}
