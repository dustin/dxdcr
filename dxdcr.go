package main

import (
	"flag"
	"log"
	"os"

	cb "github.com/couchbaselabs/go-couchbase"
)

var url = flag.String("couchbase", "http://127.0.0.1:8091", "Couchbase REST URL")

func main() {
	flag.Parse()
	log.Printf("%s\n", os.Args[0])
	flag.VisitAll(func(f *flag.Flag) { log.Printf("  -%s=%s\n", f.Name, f.Value) })

	c, err := cb.Connect(*url)
	if err != nil {
		log.Fatalf("error: could not connect to url: %s, err: %v", *url, err)
	}
	log.Printf("c: %#v", c)
}
