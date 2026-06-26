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

type getServiceListResp struct {
	parspackResponse
	Data []service `json:"data"`
}

type service struct {
	Id           string `json:"id"`
	Uuid         string `json:"uuid"`
	TargetDomain string `json:"target_domain"`
	Status       string `json:"status"`
	Plan         string `json:"plan"`
	ExpireAt     string `json:"expire_at"`
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
	ttl := time.Duration(d.Ttl) * time.Second

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
		default:
			libdnsRecord, err := libdns.RR{
				Name: name,
				TTL:  ttl,
				Type: d.Type,
				Data: r.Content,
			}.Parse()
			if err != nil {
				return nil, err
			}
			libdnsRecords = append(libdnsRecords, libdnsRecord)
		}

	}

	return libdnsRecords, nil
}

type storeDnsData struct {
	Host          string             `json:"host"`
	Type          string             `json:"type"`
	Ttl           int                `json:"ttl"`
	Proxy         string             `json:"proxy"`
	LoadBalanceId *string            `json:"load_balance_id"`
	Record        storeDnsDataRecord `json:"record"`
}

type storeDnsDataRecord struct {
	Content  string  `json:"content"`
	Port     *uint16 `json:"port"`
	Weight   *uint16 `json:"weight"`
	Priority *uint16 `json:"priority"`
	Flags    *uint8  `json:"flags"`
	Tag      *string `json:"tag"`
}

type storeDnsDataResp struct {
	parspackResponse
	Data []struct{} `json:"data"`
}

type deleteDnsData struct {
	Host   string              `json:"host"`
	Type   string              `json:"type"`
	Record deleteDnsDataRecord `json:"record"`
}

type deleteDnsDataRecord struct {
	Content  string  `json:"content"`
	Port     *uint16 `json:"port"`
	Weight   *uint16 `json:"weight"`
	Priority *uint16 `json:"priority"`
	Flags    *uint8  `json:"flags"`
	Tag      *string `json:"tag"`
}

type deleteDnsDataResp struct {
	parspackResponse
	Data []struct{} `json:"data"`
}

func toParsPackDnsData(input libdns.Record) (dnsData, error) {
	rr := input.RR()

	parspackData := dnsData{
		Host:  rr.Name,
		Type:  rr.Type,
		Ttl:   int(rr.TTL.Seconds()),
		Proxy: "direct", // no CDN by default
	}

	record := dnsRecord{
		Content: rr.Data,
	}

	switch r := input.(type) {
	case libdns.Address:
		record.Content = r.IP.String()
	case libdns.CAA:
		record.Content = r.Value
		record.Flags = &r.Flags
		record.Tag = &r.Tag
	case libdns.CNAME:
		record.Content = r.Target
	case libdns.MX:
		record.Content = r.Target
		record.Priority = &r.Preference
	case libdns.NS:
		record.Content = r.Target
	case libdns.SRV:
		record.Content = r.Target
		record.Priority = &r.Priority
		record.Weight = &r.Weight
		record.Port = &r.Port
	case libdns.TXT:
		record.Content = r.Text
	}

	parspackData.Records = []dnsRecord{record}

	return parspackData, nil
}
