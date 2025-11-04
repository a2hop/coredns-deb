package nat64

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type NAT64 struct {
	Next   plugin.Handler
	Prefix string // NAT64 prefix, default: "64:ff9b::" (RFC 6052 Well-Known Prefix)
}

func (n NAT64) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	// Block A queries
	if state.QType() == dns.TypeA {
		m := new(dns.Msg)
		m.SetRcode(r, dns.RcodeNameError)
		w.WriteMsg(m)
		return dns.RcodeNameError, nil
	}

	// For AAAA queries, synthesize from A records
	if state.QType() == dns.TypeAAAA {
		// Create A query
		aReq := r.Copy()
		aReq.Question[0].Qtype = dns.TypeA

		// Query upstream for A record
		aResp := &dns.Msg{}
		rec := &ResponseWriter{ResponseWriter: w, msg: aResp}

		n.Next.ServeDNS(ctx, rec, aReq)

		// Synthesize AAAA from A
		m := new(dns.Msg)
		m.SetReply(r)

		for _, ans := range aResp.Answer {
			if a, ok := ans.(*dns.A); ok {
				aaaa := &dns.AAAA{
					Hdr: dns.RR_Header{
						Name:   a.Hdr.Name,
						Rrtype: dns.TypeAAAA,
						Class:  dns.ClassINET,
						Ttl:    a.Hdr.Ttl,
					},
					AAAA: n.synthesizeIPv6(a.A),
				}
				m.Answer = append(m.Answer, aaaa)
			}
		}

		w.WriteMsg(m)
		return dns.RcodeSuccess, nil
	}

	// Pass through other queries
	return plugin.NextOrFailure(n.Name(), n.Next, ctx, w, r)
}

func (n NAT64) synthesizeIPv6(ipv4 net.IP) net.IP {
	// RFC 6052 Well-Known Prefix: 64:ff9b::/96
	// Prefix: 64:ff9b:: (96 bits)
	// IPv4: a.b.c.d (32 bits)
	// Result: 64:ff9b::a.b.c.d

	prefix := net.ParseIP(n.Prefix)
	ipv6 := make(net.IP, 16)
	copy(ipv6, prefix)

	// Embed IPv4 in last 32 bits
	copy(ipv6[12:], ipv4.To4())

	return ipv6
}

func (n NAT64) Name() string { return "nat64" }

type ResponseWriter struct {
	dns.ResponseWriter
	msg *dns.Msg
}

func (r *ResponseWriter) WriteMsg(m *dns.Msg) error {
	r.msg = m
	return nil
}
