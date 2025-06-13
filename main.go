package main

import (
	"context"
	"flag"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/example/set-cloudflare-dns/dns"
)

func main() {
	zone := flag.String("zone", "", "zone name (example.com)")
	recordName := flag.String("name", "", "record name (host.example.com)")
	recordType := flag.String("type", "A", "record type (A, CNAME, etc)")
	content := flag.String("content", "", "record content")
	ttl := flag.Int("ttl", 1, "record TTL")
	token := flag.String("token", os.Getenv("CF_API_TOKEN"), "cloudflare API token")
	verbose := flag.Bool("v", false, "enable verbose logging")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if *verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if *zone == "" || *recordName == "" || *content == "" || *token == "" {
		log.Fatal().Msg("zone, name, content, and token are required")
	}

	api, err := cloudflare.NewWithAPIToken(*token)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create cloudflare api client")
	}

	err = dns.EnsureDNSRecord(context.Background(), api, *zone, *recordType, *recordName, *content, *ttl)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to set dns record")
	}

	log.Info().Msg("dns record created")
}
