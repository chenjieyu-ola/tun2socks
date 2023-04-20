package proxy

import (
	"context"
	"fmt"
	"net"

	"github.com/xjasonlyu/tun2socks/v2/component/dialer"
	M "github.com/xjasonlyu/tun2socks/v2/metadata"
	"github.com/xjasonlyu/tun2socks/v2/proxy/proto"
)

var _ Proxy = (*DnsDirect)(nil)

type DnsDirect struct {
	*Base
}

func NewDnsDirect() *DnsDirect {
	return &DnsDirect{
		Base: &Base{
			proto: proto.Direct,
		},
	}
}

func (d *DnsDirect) DialContext(ctx context.Context, metadata *M.Metadata) (net.Conn, error) {
	c, err := dialer.DialContext(ctx, "tcp", metadata.DestinationAddress())
	if err != nil {
		return nil, err
	}
	setKeepAlive(c)
	return c, nil
}

func (d *DnsDirect) DialUDP(*M.Metadata) (net.PacketConn, error) {
	pc, err := dialer.ListenPacket("udp", "")
	if err != nil {
		return nil, err
	}

	var bindAddr = &net.UDPAddr{
		Port: 53,
	}
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:53")
	if err != nil {
		return nil, fmt.Errorf("resolve udp address %s: %w", "127.0.0.1:53", err)
	}
	bindAddr.IP = udpAddr.IP

	return &dnsDirectPacketConn{PacketConn: pc, rAddr: bindAddr}, nil
}

type dnsDirectPacketConn struct {
	net.PacketConn
	rAddr net.Addr
}

func (pc *dnsDirectPacketConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	if _, ok := addr.(*net.UDPAddr); ok {
		return pc.PacketConn.WriteTo(b, pc.rAddr)
	}

	_, err := net.ResolveUDPAddr("udp", addr.String())
	if err != nil {
		return 0, err
	}
	return pc.PacketConn.WriteTo(b, pc.rAddr)
}
