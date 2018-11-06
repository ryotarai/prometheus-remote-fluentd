package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/fluent/fluent-logger-golang/fluent"
)

func main() {
	fluentPort := flag.Int("fluent-port", 24224, "Fluentd port")
	fluentHost := flag.String("fluent-host", "", "Fluentd host")
	fluentTag := flag.String("fluent-tag", "", "Fluentd tag")
	listen := flag.String("listen", ":8080", "Address to listen on")
	pprof := flag.String("pprof", "", "To enable pprof, pass address to listen such as 'localhost:6060'")
	flag.Parse()

	if *fluentHost == "" {
		log.Fatal("--fluent-host option is required")
	}
	if *fluentTag == "" {
		log.Fatal("--fluent-tag option is required")
	}

	if *pprof != "" {
		go func() {
			log.Printf("Enabling pprof on %s", *pprof)
			log.Println(http.ListenAndServe(*pprof, nil))
		}()
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
