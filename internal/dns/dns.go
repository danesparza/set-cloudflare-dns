package dns

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
)

// API defines the methods needed from Cloudflare API.
type API interface {
	ZoneIDByName(zoneName string) (string, error)
	CreateDNSRecord(ctx context.Context, rc *cloudflare.ResourceContainer, params cloudflare.CreateDNSRecordParams) (cloudflare.DNSRecord, error)
	ListDNSRecords(ctx context.Context, rc *cloudflare.ResourceContainer, params cloudflare.ListDNSRecordsParams) ([]cloudflare.DNSRecord, *cloudflare.ResultInfo, error)
	UpdateDNSRecord(ctx context.Context, rc *cloudflare.ResourceContainer, params cloudflare.UpdateDNSRecordParams) (cloudflare.DNSRecord, error)
}

// EnsureDNSRecord creates or updates a DNS record for the given zone.
func EnsureDNSRecord(ctx context.Context, api API, zoneName, recordType, recordName, content string, ttl int) error {
	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return err
	}

	// See if a record already exists with the given name and type.
	records, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{
		Type: recordType,
		Name: recordName,
	})
	if err != nil {
		return err
	}

	if len(records) > 0 {
		// Update the first matched record
		update := cloudflare.UpdateDNSRecordParams{
			ID:      records[0].ID,
			Type:    recordType,
			Name:    recordName,
			Content: content,
			TTL:     ttl,
		}
		_, err = api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), update)
		return err
	}

	// Record does not exist, create it
	create := cloudflare.CreateDNSRecordParams{
		Type:    recordType,
		Name:    recordName,
		Content: content,
		TTL:     ttl,
	}
	_, err = api.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), create)
	return err
}
