package proxy

import (
	"context"
	"fmt"
	"net"

	M "github.com/xjasonlyu/tun2socks/v2/metadata"
	"github.com/xjasonlyu/tun2socks/v2/proxy/proto"
)

var _ Proxy = (*Base)(nil)

type Base struct {
	addr  string
	proto proto.Proto
}

func (b *Base) Addr() string {
	return b.addr
}

func (b *Base) Proto() proto.Proto {
	return b.proto
}

func (b *Base) DialContext(context.Context, *M.Metadata) (net.Conn, error) {
	return nil, fmt.Errorf("un support")
}

func (b *Base) DialUDP(*M.Metadata) (net.PacketConn, error) {
	return nil, fmt.Errorf("un support")
}
