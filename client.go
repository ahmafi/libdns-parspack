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

func (p *Provider) indexDnsRecord(ctx context.Context, zoneUuid string) ([]dnsData, error) {
	var body dnsDataList
	err := p.doRequest(ctx, http.MethodGet, "/external/api/v2/zones/"+zoneUuid+"/dns-records", nil, &body)
	if err != nil {
		return nil, err
	}

	if !body.Success {
		return nil, errors.New("Get DNS records error. Message:" + body.Message)
	}

	return body.Data, nil
}

func (p *Provider) storeDnsRecord(ctx context.Context, zoneUuid string, data dnsData) error {
	reqBody := storeDnsData{
		Host:  data.Host,
		Type:  data.Type,
		Ttl:   data.Ttl,
		Proxy: data.Proxy,
		Record: storeDnsDataRecord{
			Content:  data.Records[0].Content,
			Port:     data.Records[0].Port,
			Weight:   data.Records[0].Weight,
			Priority: data.Records[0].Priority,
			Flags:    data.Records[0].Flags,
			Tag:      data.Records[0].Tag,
		},
	}

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

func (p *Provider) deleteDnsRecord(ctx context.Context, zoneUuid string, data dnsData) error {
	reqBody := deleteDnsData{
		Host: data.Host,
		Type: data.Type,
		Record: deleteDnsDataRecord{
			Content: data.Records[0].Content,
		},
	}

	var body storeDnsDataResp
	err := p.doRequest(ctx, http.MethodDelete, "/external/api/v2/zones/"+zoneUuid+"/dns-records", reqBody, body)
	if err != nil {
		return err
	}

	if !body.Success {
		return errors.New("Store DNS records error. Message:" + body.Message)
	}

	return nil
}
