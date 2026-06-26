package parspack

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

const baseUrl = "https://my.parspack.com/cdnapi"

func (p *Provider) getClient() *http.Client {
	if p.client == nil {
		p.client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return p.client
}

func (p *Provider) doRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, baseUrl+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+p.APIToken)
	req.Header.Set("Content-Type", "application/json")

	return p.getClient().Do(req)
}

func (p *Provider) zoneToZoneUuid(ctx context.Context, zone string) (string, error) {
	resp, err := p.doRequest(ctx, http.MethodGet, "/external/api/v1/zones", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body serviceList
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", err
	}

	if !body.Success {
		return "", errors.New("Get service list error. Message:" + body.Message)
	}

	domain := strings.TrimSuffix(zone, ".")

	for _, service := range body.Data {
		if service.TargetDomain == domain {
			return service.Uuid, nil
		}
	}

	return "", errors.New("Zone not found")
}

func (p *Provider) indexDnsRecord(ctx context.Context, zoneUuid string) (*dnsDataList, error) {
	resp, err := p.doRequest(ctx, http.MethodGet, "/external/api/v2/zones/"+zoneUuid+"/dns-records", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body dnsDataList
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	if !body.Success {
		return nil, errors.New("Get DNS records error. Message:" + body.Message)
	}

	return &body, nil
}
