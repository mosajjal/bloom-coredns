package bloom

import (
	"os"

	"github.com/DCSO/bloom"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

// init registers this plugin.
func init() { plugin.Register("bloom", setup) }

// setup is the function that gets called when the config parser see the token "bloom". Setup is responsible
// for parsing any extra options the bloom plugin may have. The first token this function sees is "bloom".
// example config
/*
hibp.:53 {
	bloom /tmp/hibp.bloom.gz gzip
}
hashlookup.:53 {
	bloom /tmp/hashes.bloom
}
*/
func setup(c *caddy.Controller) error {
	var path, t string
	c.Next() // Ignore "bloom" and give us the next token.
	if c.NextArg() {
		// this is the full path to the bloom filter
		path = c.Val()
		if c.NextArg() {
			// it's either empty (plain) or gzip
			t = c.Val()
		} else {
			// nothing was provided, assuming plain
			t = "plain"
		}
	} else {
		return plugin.Error("bloom", c.ArgErr())
	}

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		b := Bloom{Next: next}
		if err := b.loadFilter(path, t); err != nil {
			// todo: return error here
		}
		return b
	})

	// All OK, return a nil error.
	return nil
}

// read the bloom filter
func (e *Bloom) loadFilter(path, t string) (err error) {
	var isGzip = false
	if t == "gzip" {
		isGzip = true
	}
	if f, err := os.OpenFile(path, os.O_RDONLY, 0755); err == nil {
		if b, err := bloom.LoadFromReader(f, isGzip); err == nil {
			e.bloomFilter = b
			e.isReady = true
		}
	}
	return err
}

// Ready implements the ready.Readiness interface, once this flips to true CoreDNS
// assumes this plugin is ready for queries; it is not checked again.
func (e Bloom) Ready() bool { return e.isReady }
