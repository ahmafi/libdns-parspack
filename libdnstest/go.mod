module github.com/libdns/parspack/libdnstest

go 1.23

require (
	github.com/libdns/libdns v1.1.1
	github.com/libdns/parspack v1.1.0
)

replace (
	github.com/libdns/libdns => github.com/libdns/libdns v1.2.0-alpha.1.0.20250913035451-da352cac42d0
	github.com/libdns/parspack => ../
)
