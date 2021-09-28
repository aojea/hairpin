package hairpin

import (
	"context"
	"fmt"
	"io"
	"net"
)

type hairpin struct {
	conn
}

// hairpin implements net.Conn interface
var _ net.Conn = &hairpin{}

func (h *hairpin) Read(b []byte) (int, error) {
	n, err := h.conn.Read(b)
	if err != nil && err != io.EOF && err != io.ErrClosedPipe {
		err = &net.OpError{Op: "read", Net: "Hairpin", Err: err}
	}
	return n, err
}

func (h *hairpin) Write(b []byte) (int, error) {
	n, err := h.conn.Write(b)
	if err != nil && err != io.ErrClosedPipe {
		err = &net.OpError{Op: "write", Net: "Hairpin", Err: err}
	}
	return n, err
}

type hairpinAddress struct{}

func (h hairpinAddress) Network() string {
	return "Hairpin"
}
func (h hairpinAddress) String() string {
	return "Hairpin"
}

func (h *hairpin) LocalAddr() net.Addr {
	return hairpinAddress{}
}
func (h *hairpin) RemoteAddr() net.Addr {
	return hairpinAddress{}
}

// Hairpin creates a half-duplex, in-memory, synchronous stream connection where
// data written on the connection is processed by an optional hook and then read
// back on the same connection. Reads and Write are serialized, Writes are
// blocked by Reads.
func Hairpin(fn packetHandlerFunc) *hairpin {
	return &hairpin{newConn(fn)}
}

// Dialer
type HairpinDialer struct {
	PacketHandler packetHandlerFunc
}

// Dial creates an in memory connection that is processed by the packet handler
func (h *HairpinDialer) Dial(ctx context.Context, network, address string) (net.Conn, error) {
	return Hairpin(h.PacketHandler), nil
}

// Listener
type HairpinListener struct {
	connPool []net.Conn

	PacketHandler packetHandlerFunc
}

var _ net.Listener = &HairpinListener{}

func (h *HairpinListener) Accept() (net.Conn, error) {
	conn := Hairpin(h.PacketHandler)
	return conn, nil
}

func (h *HairpinListener) Close() error {
	var aggError error
	for _, c := range h.connPool {
		if err := c.Close(); err != nil {
			aggError = fmt.Errorf("%w", err)
		}
	}
	return aggError
}

func (h *HairpinListener) Addr() net.Addr {
	return hairpinAddress{}
}

func (h *HairpinListener) Listen(network, address string) (net.Listener, error) {
	return h, nil
}
