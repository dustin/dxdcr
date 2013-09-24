package main

import (
	"flag"
	"log"
	"os"

	"github.com/couchbaselabs/go-couchbase"
	"github.com/dustin/gomemcached/client"
)

func main() {
	sourceUrl := flag.String("sourceUrl", "http://127.0.0.1:8091",
		"Source Couchbase REST URL")
	sourceBucket := flag.String("sourceBucket", "default",
		"Source bucket")
	targetUrl := flag.String("targetUrl", "http://127.0.0.1:8091",
		"Target Couchbase REST URL")
	targetBucket := flag.String("targetBucket", "default",
		"Target bucket")

	flag.Parse()
	log.Printf("%s\n", os.Args[0])
	flag.VisitAll(func(f *flag.Flag) { log.Printf("  -%s=%s\n", f.Name, f.Value) })

	start(*sourceUrl, "default", *sourceBucket, *targetUrl, "default", *targetBucket)
}

func start(srcUrl, srcPool, srcBucket, targetUrl, targetPool, targetBucket string) {
	src, err := couchbase.GetBucket(srcUrl, srcPool, srcBucket)
	if err != nil {
		log.Fatalf("could not connect to src url: %s, bucket: %s, err: %v",
			srcUrl, srcBucket, err)
	}
	defer src.Close()

	target, err := couchbase.GetBucket(targetUrl, targetPool, targetBucket)
	if err != nil {
		log.Fatalf("could not connect to target url: %s, bucket: %s, err: %v",
			targetUrl, targetBucket, err)
	}
	defer target.Close()

	tapArgs := memcached.TapArguments{
		Backfill: 0,     // Timestamp of oldest item to send.
		Dump:     false, // Disconnect after sending existing items?
		KeysOnly: false,
	}
	tap, err := src.StartTapFeed(&tapArgs)
	if err != nil {
		log.Fatalf("could not StartTapFeed, err: %v", err)
	}
	defer tap.Close()

	for e := range tap.C {
		stop, err := processTapEvent(src, target, e)
		if err != nil {
			log.Fatalf("processTapEvent err: %v", err)
		}
		if stop {
			return
		}
	}
}

func processTapEvent(src, target *couchbase.Bucket,
	e memcached.TapEvent) (stop bool, err error) {

	log.Printf(" e: %#v", e)
	return false, nil
}
