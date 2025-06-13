package dns

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudflare/cloudflare-go"
)

// fakeAPI is a simple mock of the Cloudflare API interface.
type fakeAPI struct {
	zoneID    string
	record    cloudflare.CreateDNSRecordParams
	returnErr error
}

func (f *fakeAPI) ZoneIDByName(name string) (string, error) {
	if f.returnErr != nil {
		return "", f.returnErr
	}
	return f.zoneID, nil
}

func (f *fakeAPI) CreateDNSRecord(ctx context.Context, rc *cloudflare.ResourceContainer, params cloudflare.CreateDNSRecordParams) (cloudflare.DNSRecord, error) {
	if f.returnErr != nil {
		return cloudflare.DNSRecord{}, f.returnErr
	}
	f.record = params
	return cloudflare.DNSRecord{ID: "123"}, nil
}

func TestEnsureDNSRecord(t *testing.T) {
	f := &fakeAPI{zoneID: "zone123"}
	err := EnsureDNSRecord(context.Background(), f, "example.com", "A", "test.example.com", "1.2.3.4", 120)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if f.record.Name != "test.example.com" || f.record.Content != "1.2.3.4" || f.record.TTL != 120 {
		t.Fatalf("unexpected record: %+v", f.record)
	}
}

func TestEnsureDNSRecordError(t *testing.T) {
	f := &fakeAPI{returnErr: errors.New("fail")}
	err := EnsureDNSRecord(context.Background(), f, "example.com", "A", "host", "content", 1)
	if err == nil {
		t.Fatal("expected error")
	}
}
