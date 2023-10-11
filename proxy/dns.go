package proxy

import (
	"context"
	"fmt"
	"github.com/xjasonlyu/tun2socks/v2/log"
	"net"

	"github.com/xjasonlyu/tun2socks/v2/dialer"
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
	log.Warnf("DNS go dialUdp")
	pc, err := dialer.ListenPacket("udp", "")
	if err != nil {
		fmt.Printf("resolve udp address %s: %w", "127.0.0.1:53", err)
		return nil, err
	}

	var bindAddr = &net.UDPAddr{
		Port: 53,
	}
	udpAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:53")
	if err != nil {
		return nil, fmt.Errorf("resolve udp address %s: %w", "127.0.0.1:53", err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	bindAddr.IP = udpAddr.IP
	fmt.Printf("resolve udp address %s: %s", "127.0.0.1:53", udpAddr.IP)
	return &dnsDirectPacketConn{PacketConn: pc, rAddr: bindAddr, conn: conn}, nil
}

type dnsDirectPacketConn struct {
	net.PacketConn
	rAddr net.Addr
	conn  net.Conn
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

func (pc *dnsDirectPacketConn) ReadFrom(b []byte) (int, net.Addr, error) {
	for {
		n, from, err := pc.PacketConn.ReadFrom(b)
		if from != nil {
			log.Warnf("[UDP DNS] symmetric NAT %s", from.String())
		}
		//if from != nil && from.String() != pc.dst {
		//	log.Warnf("[UDP] symmetric NAT %s->%s: drop packet from %s", pc.src, pc.dst, from)
		//	continue
		//}

		return n, from, err
	}
}
