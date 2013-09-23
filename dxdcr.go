package main

import (
	"flag"
	"log"
	"os"

	"github.com/couchbaselabs/go-couchbase"
	"github.com/dustin/gomemcached/client"
)

var sourceUrl = flag.String("sourceUrl", "http://127.0.0.1:8091", "Source Couchbase REST URL")
var sourceBucket = flag.String("sourceBucket", "default", "Source bucket")
var targetUrl = flag.String("targetUrl", "http://127.0.0.1:8091", "Target Couchbase REST URL")
var targetBucket = flag.String("targetBucket", "default", "Target bucket")

func main() {
	flag.Parse()
	log.Printf("%s\n", os.Args[0])
	flag.VisitAll(func(f *flag.Flag) { log.Printf("  -%s=%s\n", f.Name, f.Value) })

	source, err := couchbase.GetBucket(*sourceUrl, "default", *sourceBucket)
	if err != nil {
		log.Fatalf("error: could not connect to sourceUrl: %s, sourceBucket: %s, err: %v",
			*sourceUrl, sourceBucket, err)
	}
	target, err := couchbase.GetBucket(*targetUrl, "default", *targetBucket)
	if err != nil {
		log.Fatalf("error: could not connect to targetUrl: %s, targetBucket: %s, err: %v",
			*targetUrl, *targetBucket, err)
	}
	tapArgs := memcached.TapArguments{
		Backfill: 0,     // Timestamp of oldest item to send.
		Dump:     false, // If true, source will disconnect after sending existing items.
		KeysOnly: false,
	}
	tap, err := source.StartTapFeed(&tapArgs)
	if err != nil {
		log.Fatalf("error: could not StartTapFeed, err: %v", err)
	}

	log.Printf("source: %#v", source)
	log.Printf("target: %#v", target)
	log.Printf("tap: %#v", tap)
}
