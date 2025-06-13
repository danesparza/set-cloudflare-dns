package dns

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
)

// API defines the methods needed from Cloudflare API.
type API interface {
	ZoneIDByName(zoneName string) (string, error)
	CreateDNSRecord(ctx context.Context, rc *cloudflare.ResourceContainer, params cloudflare.CreateDNSRecordParams) (cloudflare.DNSRecord, error)
}

// EnsureDNSRecord creates a DNS record for the given zone.
func EnsureDNSRecord(ctx context.Context, api API, zoneName, recordType, recordName, content string, ttl int) error {
	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return err
	}
	record := cloudflare.CreateDNSRecordParams{
		Type:    recordType,
		Name:    recordName,
		Content: content,
		TTL:     ttl,
	}
	_, err = api.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), record)
	return err
}
