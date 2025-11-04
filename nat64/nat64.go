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
		rec := &ResponseWriter{ResponseWriter: w}

		rcode, err := n.Next.ServeDNS(ctx, rec, aReq)
		if err != nil {
			return rcode, err
		}

		// Get the A record response
		aResp := rec.msg
		if aResp == nil {
			// No response received, return empty AAAA response
			m := new(dns.Msg)
			m.SetReply(r)
			w.WriteMsg(m)
			return dns.RcodeSuccess, nil
		}

		// Synthesize AAAA from A
		m := new(dns.Msg)
		m.SetReply(r)
		m.Authoritative = aResp.Authoritative
		m.RecursionAvailable = aResp.RecursionAvailable
		m.Rcode = aResp.Rcode

		// Copy NS section, but filter out DNSSEC records
		for _, ns := range aResp.Ns {
			if !isDNSSECRecord(ns) {
				m.Ns = append(m.Ns, ns)
			}
		}

		// Copy Extra section, but filter out DNSSEC records
		for _, extra := range aResp.Extra {
			if !isDNSSECRecord(extra) {
				m.Extra = append(m.Extra, extra)
			}
		}

		// Synthesize AAAA records from A records
		// Also preserve CNAME records in the chain
		for _, ans := range aResp.Answer {
			switch rr := ans.(type) {
			case *dns.A:
				// Convert A to AAAA
				aaaa := &dns.AAAA{
					Hdr: dns.RR_Header{
						Name:   rr.Hdr.Name,
						Rrtype: dns.TypeAAAA,
						Class:  dns.ClassINET,
						Ttl:    rr.Hdr.Ttl,
					},
					AAAA: n.synthesizeIPv6(rr.A),
				}
				m.Answer = append(m.Answer, aaaa)
			case *dns.CNAME:
				// Preserve CNAME records
				m.Answer = append(m.Answer, rr)
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

// isDNSSECRecord checks if a DNS record is a DNSSEC-related record
func isDNSSECRecord(rr dns.RR) bool {
	switch rr.Header().Rrtype {
	case dns.TypeRRSIG, dns.TypeNSEC, dns.TypeNSEC3, dns.TypeDS, dns.TypeDNSKEY:
		return true
	}
	return false
}

func (n NAT64) Name() string { return "nat64" }

type ResponseWriter struct {
	dns.ResponseWriter
	msg *dns.Msg
}

func (r *ResponseWriter) WriteMsg(m *dns.Msg) error {
	// Deep copy the message to preserve it
	r.msg = m.Copy()
	return nil
}
