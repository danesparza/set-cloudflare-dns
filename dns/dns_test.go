package dns

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudflare/cloudflare-go"
)

// fakeAPI is a simple mock of the Cloudflare API interface.
type fakeAPI struct {
	zoneID       string
	createRecord cloudflare.CreateDNSRecordParams
	updateRecord cloudflare.UpdateDNSRecordParams
	listRecords  []cloudflare.DNSRecord
	returnErr    error
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
	f.createRecord = params
	return cloudflare.DNSRecord{ID: "123"}, nil
}

func (f *fakeAPI) ListDNSRecords(ctx context.Context, rc *cloudflare.ResourceContainer, params cloudflare.ListDNSRecordsParams) ([]cloudflare.DNSRecord, *cloudflare.ResultInfo, error) {
	if f.returnErr != nil {
		return nil, nil, f.returnErr
	}
	return f.listRecords, nil, nil
}

func (f *fakeAPI) UpdateDNSRecord(ctx context.Context, rc *cloudflare.ResourceContainer, params cloudflare.UpdateDNSRecordParams) (cloudflare.DNSRecord, error) {
	if f.returnErr != nil {
		return cloudflare.DNSRecord{}, f.returnErr
	}
	f.updateRecord = params
	return cloudflare.DNSRecord{ID: params.ID}, nil
}

func TestEnsureDNSRecordCreate(t *testing.T) {
	f := &fakeAPI{zoneID: "zone123"}
	err := EnsureDNSRecord(context.Background(), f, "example.com", "A", "test.example.com", "1.2.3.4", 120)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if f.createRecord.Name != "test.example.com" || f.createRecord.Content != "1.2.3.4" || f.createRecord.TTL != 120 {
		t.Fatalf("unexpected record: %+v", f.createRecord)
	}
	if f.updateRecord.ID != "" {
		t.Fatalf("unexpected update: %+v", f.updateRecord)
	}
}

func TestEnsureDNSRecordUpdate(t *testing.T) {
	f := &fakeAPI{
		zoneID:      "zone123",
		listRecords: []cloudflare.DNSRecord{{ID: "abc", Name: "test.example.com", Type: "A"}},
	}
	err := EnsureDNSRecord(context.Background(), f, "example.com", "A", "test.example.com", "2.2.2.2", 60)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if f.updateRecord.ID != "abc" || f.updateRecord.Content != "2.2.2.2" || f.updateRecord.TTL != 60 {
		t.Fatalf("unexpected update: %+v", f.updateRecord)
	}
	if f.createRecord.Name != "" {
		t.Fatalf("unexpected create: %+v", f.createRecord)
	}
}

func TestEnsureDNSRecordError(t *testing.T) {
	f := &fakeAPI{returnErr: errors.New("fail")}
	err := EnsureDNSRecord(context.Background(), f, "example.com", "A", "host", "content", 1)
	if err == nil {
		t.Fatal("expected error")
	}
}
