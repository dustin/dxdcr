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

	start(*sourceUrl, "default", *sourceBucket, *targetUrl, "default", *targetBucket)
}

func start(sourceUrl, sourcePool, sourceBucket, targetUrl, targetPool, targetBucket string) {
	source, err := couchbase.GetBucket(sourceUrl, sourcePool, sourceBucket)
	if err != nil {
		log.Fatalf("error: could not connect to sourceUrl: %s, sourceBucket: %s, err: %v",
			sourceUrl, sourceBucket, err)
	}
	target, err := couchbase.GetBucket(targetUrl, targetPool, targetBucket)
	if err != nil {
		log.Fatalf("error: could not connect to targetUrl: %s, targetBucket: %s, err: %v",
			targetUrl, targetBucket, err)
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

	for e := range tap.C {
		stop, err := processTapEvent(source, target, e);
		if err != nil {
			log.Fatalf("error: processTapEvent err: %v", err)
		}
		if stop {
			return
		}
	}

	tap.Close()
}

func processTapEvent(source, target *couchbase.Bucket, e memcached.TapEvent) (stop bool, err error) {
	log.Printf(" e: %#v", e)
	return false, nil
}
