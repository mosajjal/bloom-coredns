// Package bloom is a CoreDNS plugin that gets a bloom file as an input and serves it
// as a DNS zone. NXDOMAIN means no match and NOERROR means a match
package bloom

import (
	"context"
	"strings"

	"github.com/DCSO/bloom"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("bloom")

// Bloom plugin
type Bloom struct {
	Next        plugin.Handler
	bloomFilter *bloom.BloomFilter
	isReady     bool
}

// ServeDNS implements the plugin.Handler interface. This method gets called when example is used
// in a Server.
func (e Bloom) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	reply := new(dns.Msg)
	reply.SetReply(r)

	for _, q := range r.Question {

		key := strings.Split(q.Name, ".")[0]
		isPresent := e.bloomFilter.Check([]byte(key))

		//write the response
		if isPresent {
			rr, _ := dns.NewRR(key + " CNAME " + key)
			reply.Answer = append(r.Answer, rr)
			reply.SetRcode(r, dns.RcodeSuccess)
		} else {
			// rr, err := dns.NewRR("example.com. CNAME .")
			// reply.Answer = append(r.Answer, rr)
			reply.SetRcode(r, dns.RcodeNameError)
		}

		w.WriteMsg(reply)
	}

	pw := NewResponsePrinter(w)
	// Export metric with the server label set to the current server handling the request.
	requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	// Call next plugin (if any).
	return plugin.NextOrFailure(e.Name(), e.Next, ctx, pw, r)
}

// Name implements the Handler interface.
func (e Bloom) Name() string { return "bloom" }

// ResponsePrinter wrap a dns.ResponseWriter and will write example to standard output when WriteMsg is called.
type ResponsePrinter struct {
	dns.ResponseWriter
}

// NewResponsePrinter returns ResponseWriter.
func NewResponsePrinter(w dns.ResponseWriter) *ResponsePrinter {
	return &ResponsePrinter{ResponseWriter: w}
}

// WriteMsg calls the underlying ResponseWriter's WriteMsg method and prints "example" to standard output.
func (r *ResponsePrinter) WriteMsg(res *dns.Msg) error {
	return r.ResponseWriter.WriteMsg(res)
}
