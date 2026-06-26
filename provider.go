// Package libdnstemplate implements a DNS record management client compatible
// with the libdns interfaces for ParsPack.
package parspack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/libdns/libdns"
)

// TODO: Providers must not require additional provisioning steps by the callers; it
// should work simply by populating a struct and calling methods on it. If your DNS
// service requires long-lived state or some extra provisioning step, do it implicitly
// when methods are called; sync.Once can help with this, and/or you can use a
// sync.(RW)Mutex in your Provider struct to synchronize implicit provisioning.

// Provider facilitates DNS record manipulation with ParsPack.
type Provider struct {
	APIToken string `json:"api_token,omitempty"`

	client *http.Client
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	zoneUuid, err := p.zoneToZoneUuid(ctx, zone)
	if err != nil {
		return nil, err
	}

	dnsData, err := p.indexDnsRecord(ctx, zoneUuid)
	if err != nil {
		return nil, err
	}

	libdnsRecords := make([]libdns.Record, 0, len(dnsData.Data))
	for _, d := range dnsData.Data {
		records, err := d.libdnsRecord()
		if err != nil {
			return nil, err
		}
		libdnsRecords = append(libdnsRecords, records...)
	}

	return libdnsRecords, nil
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	zoneUuid, err := p.zoneToZoneUuid(ctx, zone)
	if err != nil {
		return nil, err
	}

	var createdRecords []libdns.Record
	for _, r := range records {
		parsPackData, err := toParsPackStoreDnsData(r)
		if err != nil {
			return nil, err
		}
		err = p.storeDnsRecord(ctx, zoneUuid, parsPackData)
		if err != nil {
			// should we ignore this error or return when the record already exists?
			return createdRecords, err
		}
		createdRecords = append(createdRecords, r)
	}

	return createdRecords, nil
}

// SetRecords sets the records in the zone, either by updating existing records or creating new ones.
// It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	// Make sure to return RR-type-specific structs, not libdns.RR structs.
	return nil, fmt.Errorf("TODO: not implemented")
}

// DeleteRecords deletes the specified records from the zone. It returns the records that were deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	// Make sure to return RR-type-specific structs, not libdns.RR structs.
	return nil, fmt.Errorf("TODO: not implemented")
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
