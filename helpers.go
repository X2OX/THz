package THz

import (
	"bytes"
	"net"
	"time"
)

type testConn struct {
	r bytes.Buffer
	w bytes.Buffer
}

func (c *testConn) Read(b []byte) (int, error)  { return c.r.Read(b) }
func (c *testConn) Write(b []byte) (int, error) { return c.w.Write(b) }
func (c *testConn) Close() error                { return nil }

func (c *testConn) LocalAddr() net.Addr                { return &net.TCPAddr{Port: 0, Zone: "", IP: net.IPv4zero} }
func (c *testConn) RemoteAddr() net.Addr               { return &net.TCPAddr{Port: 0, Zone: "", IP: net.IPv4zero} }
func (c *testConn) SetDeadline(_ time.Time) error      { return nil }
func (c *testConn) SetReadDeadline(_ time.Time) error  { return nil }
func (c *testConn) SetWriteDeadline(_ time.Time) error { return nil }
