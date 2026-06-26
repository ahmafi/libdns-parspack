package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/libdns/libdns/libdnstest"
	"github.com/libdns/parspack"
)

func TestParsPackProvider(t *testing.T) {
	apiToken := os.Getenv("PARSPACK_API_TOKEN")
	testZone := os.Getenv("PARSPACK_TEST_ZONE")

	if apiToken == "" {
		t.Skip("Skipping ParsPack provider tests: PARSPACK_API_TOKEN environment variables must be set")
	}

	if !strings.HasSuffix(testZone, ".") {
		t.Fatal("We expect the test zone to to have trailing dot")
	}

	provider := &parspack.Provider{
		APIToken: apiToken,
	}

	suite := libdnstest.NewTestSuite(provider, testZone)
	suite.Timeout = 80 * time.Second
	suite.SkipRRTypes = map[string]bool{
		"AAAA":  true, // Skip MX record tests
		"SVCB":  true, // Skip SVCB record tests
		"HTTPS": true, // Skip HTTPS record tests
	}
	suite.RunTests(t)
}
