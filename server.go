package main

import (
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
)

type sample struct {
	Value      float64           `msg:"value"`
	ValueIsNaN bool              `msg:"value_is_nan"`
	MetricName string            `msg:"metric_name"`
	Labels     map[string]string `msg:"labels"`
}

type writer interface {
	Post(time.Time, sample) error
}

type Server struct {
	output writer
}

func NewServer(output writer) (*Server, error) {
	s := &Server{
		output: output,
	}

	return s, nil
}

func (s *Server) handleWrite(w http.ResponseWriter, r *http.Request) {
	compressed, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req prompb.WriteRequest
	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.writeTimeseries(req.Timeseries)
}

func (s *Server) writeTimeseries(tss []*prompb.TimeSeries) error {
	log.Printf("Receiving %d timeseries", len(tss))

	for _, ts := range tss {
		for _, ss := range ts.Samples {
			sample := sample{}
			if math.IsNaN(ss.Value) {
				sample.ValueIsNaN = true
			} else {
				sample.ValueIsNaN = false
				sample.Value = ss.Value
			}
			labels := map[string]string{}
			for _, l := range ts.Labels {
				if l.Name == "__name__" {
					sample.MetricName = l.Value
					continue
				}
				labels[l.Name] = l.Value
			}
			sample.Labels = labels

			t := time.Unix(0, ss.Timestamp*1000000)
			err := s.output.Post(t, sample)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/write" {
		s.handleWrite(w, r)
	} else {
		http.NotFound(w, r)
	}
}
