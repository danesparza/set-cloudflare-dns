package dns

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudflare/cloudflare-go"
	"github.com/danesparza/set-cloudflare-dns/internal/dns/mocks"
	"github.com/stretchr/testify/mock"
)

func TestEnsureDNSRecordCreate(t *testing.T) {
	m := mocks.NewAPI(t)
	m.On("ZoneIDByName", "example.com").Return("zone123", nil)
	m.On("ListDNSRecords", mock.Anything, cloudflare.ZoneIdentifier("zone123"), mock.Anything).Return([]cloudflare.DNSRecord{}, nil, nil)

	var create cloudflare.CreateDNSRecordParams
	m.On("CreateDNSRecord", mock.Anything, cloudflare.ZoneIdentifier("zone123"), mock.Anything).Run(func(args mock.Arguments) {
		create = args.Get(2).(cloudflare.CreateDNSRecordParams)
	}).Return(cloudflare.DNSRecord{ID: "123"}, nil)

	err := EnsureDNSRecord(context.Background(), m, "example.com", "A", "test.example.com", "1.2.3.4", 120)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if create.Name != "test.example.com" || create.Content != "1.2.3.4" || create.TTL != 120 {
		t.Fatalf("unexpected record: %+v", create)
	}
}

func TestEnsureDNSRecordUpdate(t *testing.T) {
	m := mocks.NewAPI(t)
	m.On("ZoneIDByName", "example.com").Return("zone123", nil)
	m.On("ListDNSRecords", mock.Anything, cloudflare.ZoneIdentifier("zone123"), mock.Anything).Return([]cloudflare.DNSRecord{{ID: "abc", Name: "test.example.com", Type: "A"}}, nil, nil)

	var update cloudflare.UpdateDNSRecordParams
	m.On("UpdateDNSRecord", mock.Anything, cloudflare.ZoneIdentifier("zone123"), mock.Anything).Run(func(args mock.Arguments) {
		update = args.Get(2).(cloudflare.UpdateDNSRecordParams)
	}).Return(cloudflare.DNSRecord{ID: "abc"}, nil)

	err := EnsureDNSRecord(context.Background(), m, "example.com", "A", "test.example.com", "2.2.2.2", 60)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if update.ID != "abc" || update.Content != "2.2.2.2" || update.TTL != 60 {
		t.Fatalf("unexpected update: %+v", update)
	}
}

func TestEnsureDNSRecordError(t *testing.T) {
	m := mocks.NewAPI(t)
	m.On("ZoneIDByName", "example.com").Return("", errors.New("fail"))

	err := EnsureDNSRecord(context.Background(), m, "example.com", "A", "host", "content", 1)
	if err == nil {
		t.Fatal("expected error")
	}
}
