# Bloom Filter CoreDNS Plugin

Serves up a `bloom` file as a DNS zone. Useful to make bloom filter queries available via the network


## Compilation

This package will always be compiled as part of CoreDNS and not in a standalone way. It will require you to use `go get` or as a dependency on [plugin.cfg](https://github.com/coredns/coredns/blob/master/plugin.cfg).

The [manual](https://coredns.io/manual/toc/#what-is-coredns) will have more information about how to configure and extend the server with external plugins.

A simple way to consume this plugin, is by adding the following on [plugin.cfg](https://github.com/coredns/coredns/blob/master/plugin.cfg), and recompile it as [detailed on coredns.io](https://coredns.io/2017/07/25/compile-time-enabling-or-disabling-plugins/#build-with-compile-time-configuration-file).

~~~
bloom:github.com/mosajjal/bloom-coredns
~~~

Put this early in the plugin list, so that *bloom* is executed before any of the other plugins.

After this you can compile coredns by:

``` sh
go generate
go build
```

Or you can instead use make:

``` sh
make
```

## Syntax

~~~ txt
bloom
~~~

## Metrics

If monitoring is enabled (via the *prometheus* directive) the following metric is exported:

* `coredns_bloom_request_count_total{server}` - query count to the *bloom* plugin.

The `server` label indicated which server handled the request, see the *metrics* plugin for details.

## Ready

This plugin reports readiness to the ready plugin. It will be immediately ready.

## Examples

In this configuration, we configure a bloom filter consisting SHA1 hashes of known files in the OS

~~~ corefile
. {
    bloom /opt/mybloom/hashes.bloom  
}
~~~

If the bloom is gzipped:

~~~ corefile
. {
    bloom /opt/mybloom/hashes.bloom.gz gzip
}
~~~

## Also See

See the [manual](https://coredns.io/manual).
