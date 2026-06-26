package parspack

import (
	"bytes"
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

func (p *Provider) doRequest(ctx context.Context, method, path string, reqBody any, resBody any) error {
	var body io.Reader

	if reqBody != nil {
		var reqBodyBuf bytes.Buffer
		if err := json.NewEncoder(&reqBodyBuf).Encode(reqBody); err != nil {
			return err
		}
		body = &reqBodyBuf
	}

	req, err := http.NewRequestWithContext(ctx, method, baseUrl+path, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+p.APIToken)
	req.Header.Set("Accept", "application/json")

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")

	}

	resp, err := p.getClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(body))

	err = json.NewDecoder(resp.Body).Decode(&resBody)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) zoneToZoneUuid(ctx context.Context, zone string) (string, error) {
	var body serviceList

	err := p.doRequest(ctx, http.MethodGet, "/external/api/v1/zones", nil, &body)
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
	var body dnsDataList
	err := p.doRequest(ctx, http.MethodGet, "/external/api/v2/zones/"+zoneUuid+"/dns-records", nil, &body)
	if err != nil {
		return nil, err
	}

	if !body.Success {
		return nil, errors.New("Get DNS records error. Message:" + body.Message)
	}

	return &body, nil
}

func (p *Provider) storeDnsRecord(ctx context.Context, zoneUuid string, reqBody storeDnsData) error {
	var body storeDnsDataResp
	err := p.doRequest(ctx, http.MethodPost, "/external/api/v2/zones/"+zoneUuid+"/dns-records", reqBody, body)
	if err != nil {
		return err
	}

	if !body.Success {
		return errors.New("Store DNS records error. Message:" + body.Message)
	}

	return nil
}
