package parspack

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func printJSON(obj any) {
	bytes, _ := json.MarshalIndent(obj, "  ", "  ")
	fmt.Println(string(bytes))
}

func getProvider() Provider {
	apiToken := os.Getenv("PARSPACK_API_TOKEN")

	return Provider{
		APIToken: apiToken,
	}
}

func getZone() string {
	return os.Getenv("PARSPACK_ZONE")
}

func TestZoneToZoneUuid(t *testing.T) {
	ctx := context.Background()
	p := getProvider()
	zone := getZone()

	got, err := p.zoneToZoneUuid(ctx, zone)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fmt.Println("got: " + got)
}

func TestIndexDnsRecord(t *testing.T) {
	ctx := context.Background()
	p := getProvider()
	zone := getZone()

	zoneUuid, err := p.zoneToZoneUuid(ctx, zone)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fmt.Println("got: " + zoneUuid)

	dnsRecords, err := p.indexDnsRecord(ctx, zoneUuid)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	printJSON(dnsRecords)
}
