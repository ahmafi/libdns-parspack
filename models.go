package parspack

import (
	"fmt"
	"net/netip"
	"time"

	"github.com/libdns/libdns"
)

type parspackResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type serviceList struct {
	parspackResponse
	Data []struct {
		Id           string `json:"id"`
		Uuid         string `json:"uuid"`
		TargetDomain string `json:"target_domain"`
		Status       string `json:"status"`
		Plan         string `json:"plan"`
		ExpireAt     string `json:"expire_at"`
	} `json:"data"`
}

type dnsDataList struct {
	parspackResponse
	Data []dnsData `json:"data"`
}

type dnsData struct {
	Zone    string      `json:"zone"`
	Ttl     int         `json:"ttl"`
	Type    string      `json:"type"`
	Host    string      `json:"host"`
	Proxy   string      `json:"proxy"`
	Records []dnsRecord `json:"records"`
}

type dnsRecord struct {
	Content  string  `json:"content"`
	Disabled bool    `json:"disabled"`
	Port     *uint16 `json:"port"`
	Weight   *uint16 `json:"weight"`
	Priority *uint16 `json:"priority"`
	Flags    *uint8  `json:"flags"`
	Tag      *string `json:"tag"`
	Serial   *string `json:"serial"`
	Refresh  *string `json:"refresh"`
	Minimum  *string `json:"minimum"`
}

func (d dnsData) libdnsRecord() ([]libdns.Record, error) {
	name := d.Host
	ttl := time.Duration(d.Ttl)

	libdnsRecords := make([]libdns.Record, 0, len(d.Records))

	for _, r := range d.Records {
		switch d.Type {
		case "A", "AAAA":
			ip, err := netip.ParseAddr(r.Content)
			if err != nil {
				return nil, fmt.Errorf("unexpected type for A/AAAA value: %T", r.Content)
			}
			libdnsRecords = append(libdnsRecords, libdns.Address{
				Name: name,
				TTL:  ttl,
				IP:   ip,
			})
		case "CAA":
			libdnsRecords = append(libdnsRecords, libdns.CAA{
				Name:  name,
				TTL:   ttl,
				Value: r.Content,
				Flags: *r.Flags,
				Tag:   *r.Tag,
			})
		case "CNAME":
			libdnsRecords = append(libdnsRecords, libdns.CNAME{
				Name:   name,
				TTL:    ttl,
				Target: r.Content,
			})
		case "MX":
			libdnsRecords = append(libdnsRecords, libdns.MX{
				Name:       name,
				TTL:        ttl,
				Target:     r.Content,
				Preference: *r.Priority,
			})
		case "NS":
			libdnsRecords = append(libdnsRecords, libdns.NS{
				Name:   name,
				TTL:    ttl,
				Target: r.Content,
			})
		case "SRV":
			libdnsRecords = append(libdnsRecords, libdns.SRV{
				Name:     name,
				TTL:      ttl,
				Target:   r.Content,
				Priority: *r.Priority,
				Weight:   *r.Weight,
				Port:     *r.Port,
			})
		case "TXT":
			libdnsRecords = append(libdnsRecords, libdns.TXT{
				Name: name,
				TTL:  ttl,
				Text: r.Content,
			})
		}

	}

	return libdnsRecords, nil
}
