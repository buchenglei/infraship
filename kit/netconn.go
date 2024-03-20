package kit

import (
	"net"
	"strings"
)

func ReadNetConnIP(addr net.Addr) string {
	items := strings.Split(addr.String(), ":")
	if len(items) != 2 {
		return addr.String()
	}

	return items[0]
}

func ReadNetConnPort(addr net.Addr) string {
	items := strings.Split(addr.Network(), ":")
	if len(items) != 2 {
		return addr.Network()
	}

	return items[1]
}
