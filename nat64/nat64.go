package nat64

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
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
		log.Infof("[NAT64] Blocking A query for %s", state.QName())
		m := new(dns.Msg)
		m.SetRcode(r, dns.RcodeNameError)
		w.WriteMsg(m)
		return dns.RcodeNameError, nil
	}

	// For AAAA queries, synthesize from A records
	if state.QType() == dns.TypeAAAA {
		log.Infof("[NAT64] Processing AAAA query for %s", state.QName())

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
		// Note: We do NOT filter AAAA records here as they may be needed for glue records
		for _, extra := range aResp.Extra {
			if !isDNSSECRecord(extra) {
				m.Extra = append(m.Extra, extra)
			}
		}

		// Synthesize AAAA records from A records
		// Also preserve CNAME records in the chain
		hasAnswers := false
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
				hasAnswers = true
			case *dns.CNAME:
				// Preserve CNAME records
				m.Answer = append(m.Answer, rr)
				hasAnswers = true
			case *dns.AAAA:
				// IMPORTANT: Never include real AAAA records in NAT64 responses
				// This would bypass NAT64 and break IPv6-only networks
				continue
			default:
				// For any other record types in the answer, skip them
				continue
			}
		}

		// If no valid answers were synthesized, return NXDOMAIN or empty response
		if !hasAnswers {
			m.Rcode = aResp.Rcode
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
	if prefix == nil {
		log.Errorf("[NAT64] Failed to parse prefix '%s', using default 64:ff9b::", n.Prefix)
		prefix = net.ParseIP("64:ff9b::")
	}

	ipv6 := make(net.IP, 16)
	copy(ipv6, prefix)

	// Embed IPv4 in last 32 bits
	ipv4Bytes := ipv4.To4()
	if ipv4Bytes == nil {
		log.Errorf("[NAT64] Invalid IPv4 address: %s", ipv4)
		// If ipv4 is already IPv6 or invalid, return as-is
		return ipv4
	}
	copy(ipv6[12:], ipv4Bytes)

	log.Debugf("[NAT64] Prefix=%s + IPv4=%s = IPv6=%s", n.Prefix, ipv4, ipv6)

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
