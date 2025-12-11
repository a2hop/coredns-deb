package nat664

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("nat664", setup) }

func setup(c *caddy.Controller) error {
	// Use RFC 6052 Well-Known Prefix as default
	n := NAT664{Prefix: "64:ff9b::"}

	for c.Next() {
		if c.NextArg() {
			n.Prefix = c.Val()
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		n.Next = next
		return n
	})

	return nil
}
