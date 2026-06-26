package parspack

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/libdns/libdns"
)

func TestGetRecords(t *testing.T) {
	ctx := context.Background()
	p := getProvider()
	zone := getZone()

	got, err := p.GetRecords(ctx, zone)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	printJSON(got)
}

func TestAppendRecords(t *testing.T) {
	ctx := context.Background()
	p := getProvider()
	zone := getZone()

	records := []libdns.Record{
		libdns.Address{
			Name: "bag",
			TTL:  time.Second * 3600,
			IP:   netip.MustParseAddr("109.122.247.242"),
		},
		libdns.TXT{
			Name: "ame",
			TTL:  time.Second * 60,
			Text: "test-acme-me-2",
		},
	}

	got, err := p.AppendRecords(ctx, zone, records)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	printJSON(got)
}
