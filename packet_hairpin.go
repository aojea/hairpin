package hairpin

import (
	"context"
	"fmt"
	"io"
	"net"
)

type packetHairpin struct {
	conn
}

// packetHairpin implements net.PacketConn interface
var _ net.PacketConn = &packetHairpin{}

func (p *packetHairpin) ReadFrom(b []byte) (int, net.Addr, error) {
	n, err := p.conn.Read(b)
	if err != nil && err != io.EOF && err != io.ErrClosedPipe {
		err = &net.OpError{Op: "read", Net: "PacketHairpin", Err: err}
	}
	return n, packetHairpinAddress{}, err
}

func (p *packetHairpin) WriteTo(b []byte, _ net.Addr) (int, error) {
	n, err := p.conn.Write(b)
	if err != nil && err != io.ErrClosedPipe {
		err = &net.OpError{Op: "write", Net: "PacketHairpin", Err: err}
	}
	return n, err
}

type packetHairpinAddress struct{}

func (p packetHairpinAddress) Network() string {
	return "packetHairpin"
}
func (p packetHairpinAddress) String() string {
	return "packetHairpin"
}

func (p *packetHairpin) LocalAddr() net.Addr {
	return packetHairpinAddress{}
}
func (p *packetHairpin) RemoteAddr() net.Addr {
	return packetHairpinAddress{}
}

// packetHairpin creates a half-duplex, in-memory, synchronous packet connection where
// data written on the connection is processed by an optional hook and then read
// back on the same connection. Reads and Write are serialized, Writes are
// blocked by Reads.
func PacketHairpin(fn packetHandlerFunc) net.Conn {
	return &packetHairpin{newConn(fn)}
}

// Dialer
type PacketHairpinDialer struct {
	PacketHandler packetHandlerFunc
}

// Dial creates an in memory connection that is processed by the packet handler
func (p *PacketHairpinDialer) Dial(ctx context.Context, network, address string) (net.Conn, error) {
	return PacketHairpin(p.PacketHandler), nil
}

// Listener
type PacketHairpinListener struct {
	connPool []net.PacketConn

	PacketHandler packetHandlerFunc
}

var _ net.Listener = &PacketHairpinListener{}

func (p *PacketHairpinListener) Accept() (net.Conn, error) {
	return PacketHairpin(p.PacketHandler), nil
}

func (p *PacketHairpinListener) Close() error {
	var aggError error
	for _, c := range p.connPool {
		if err := c.Close(); err != nil {
			aggError = fmt.Errorf("%w", err)
		}
	}
	return aggError
}

func (p *PacketHairpinListener) Addr() net.Addr {
	return packetHairpinAddress{}
}

func (p *PacketHairpinListener) Listen(network, address string) (net.Listener, error) {
	return p, nil
}
