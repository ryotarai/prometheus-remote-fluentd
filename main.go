package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/fluent/fluent-logger-golang/fluent"
)

func main() {
	fluentPort := flag.Int("fluent-port", 24224, "Fluentd port")
	fluentHost := flag.String("fluent-host", "", "Fluentd host")
	fluentTag := flag.String("fluent-tag", "", "Fluentd tag")
	listen := flag.String("listen", ":8080", "Address to listen on")
	flag.Parse()

	if *fluentHost == "" {
		log.Fatal("--fluent-host option is required")
	}
	if *fluentTag == "" {
		log.Fatal("--fluent-tag option is required")
	}

	f, err := fluent.New(fluent.Config{
		FluentHost: *fluentHost,
		FluentPort: *fluentPort,
	})
	if err != nil {
		log.Fatal(err)
	}

	w := NewFluentWriter(f, *fluentTag)
	s, err := NewServer(w)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening %s", *listen)
	err = http.ListenAndServe(*listen, s)
	if err != nil {
		log.Fatal(err)
	}
}
