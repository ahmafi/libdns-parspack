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

	libdnsRecords := make([]libdns.Record, 0, len(dnsData))
	for _, d := range dnsData {
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
		parsPackData, err := toParsPackDnsData(r)
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
	zoneUuid, err := p.zoneToZoneUuid(ctx, zone)
	if err != nil {
		return nil, err
	}

	existingDnsList, err := p.indexDnsRecord(ctx, zoneUuid)
	if err != nil {
		return nil, err
	}

	var toDelete []dnsData
	for _, existingDns := range existingDnsList {
		existingLibdnsList, err := existingDns.libdnsRecord()
		if err != nil {
			return nil, err
		}

		for _, existingLibDns := range existingLibdnsList {
			for _, input := range records {
				if existingLibDns.RR().Name == input.RR().Name && existingLibDns.RR().Type == input.RR().Type {
					toDelete = append(toDelete, existingDns)
				}
			}
		}

	}

	for _, d := range toDelete {
		p.deleteDnsRecord(ctx, zoneUuid, d)
	}

	return p.AppendRecords(ctx, zone, records)
}

// DeleteRecords deletes the specified records from the zone. It returns the records that were deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	zoneUuid, err := p.zoneToZoneUuid(ctx, zone)
	if err != nil {
		return nil, err
	}

	existingDns, err := p.indexDnsRecord(ctx, zoneUuid)
	if err != nil {
		return nil, err
	}

	var toDelete []dnsData
	for _, r := range records {
		inputDns, err := toParsPackDnsData(r)
		if err != nil {
			return nil, err
		}

		for _, existingDns := range existingDns {
			nameMatch := inputDns.Host == existingDns.Host
			typeMatch := inputDns.Type == existingDns.Type || inputDns.Type == ""
			ttlMatch := inputDns.Ttl == existingDns.Ttl || inputDns.Ttl == 0

			if nameMatch && typeMatch && ttlMatch {
				for _, record := range existingDns.Records {
					valueMatch := inputDns.Records[0].Content == record.Content || inputDns.Records[0].Content == ""
					if valueMatch {
						toDelete = append(toDelete, inputDns)
						break
					}
				}
			}
		}

	}

	for _, d := range toDelete {
		p.deleteDnsRecord(ctx, zoneUuid, d)
	}

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
