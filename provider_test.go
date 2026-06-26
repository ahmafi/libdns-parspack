package parspack

import (
	"context"
	"testing"
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
