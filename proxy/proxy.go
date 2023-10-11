// Package proxy provides implementations of proxy protocols.
package proxy

import (
	"context"
	"github.com/xjasonlyu/tun2socks/v2/log"
	M "github.com/xjasonlyu/tun2socks/v2/metadata"
	"github.com/xjasonlyu/tun2socks/v2/proxy/proto"
	"net"
	"time"
)

const (
	tcpConnectTimeout = 5 * time.Second
)

var _defaultDialer Dialer = &Base{}
var _dnsDialer Dialer = NewDnsDirect()

type Dialer interface {
	DialContext(context.Context, *M.Metadata) (net.Conn, error)
	DialUDP(*M.Metadata) (net.PacketConn, error)
}

type Proxy interface {
	Dialer
	Addr() string
	Proto() proto.Proto
}

// SetDialer sets default Dialer.
func SetDialer(d Dialer) {
	_defaultDialer = d
}

// Dial uses default Dialer to dial TCP.
func Dial(metadata *M.Metadata) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), tcpConnectTimeout)
	defer cancel()
	return _defaultDialer.DialContext(ctx, metadata)
}

// DialContext uses default Dialer to dial TCP with context.
func DialContext(ctx context.Context, metadata *M.Metadata) (net.Conn, error) {
	return _defaultDialer.DialContext(ctx, metadata)
}

// DialUDP uses default Dialer to dial UDP.
func DialUDP(metadata *M.Metadata) (net.PacketConn, error) {
	log.Warnf("proxy go dialUDP %d", metadata.DstPort)
	//if metadata.DstPort == 53 {
	//	return _dnsDialer.DialUDP(metadata)
	//}
	return _defaultDialer.DialUDP(metadata)
}
