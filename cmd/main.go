package main

import (
	"bufio"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sysradium/request-counter/internal/counter/decorators"
	"github.com/sysradium/request-counter/internal/counter/ephemeral"
)

var (
	dataFilePath = "file.data"
)

type storage interface {
	Add(time.Time) error
	Len() int
}

func counterHandler(
	c storage,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := c.Add(time.Now()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		requestsCount := fmt.Sprint(c.Len())
		w.Write([]byte(requestsCount))
	}
}

func load(f *os.File) []time.Time {
	bufferedReader := bufio.NewReader(f)

	dec := gob.NewDecoder(bufferedReader)
	var times []time.Time
	for {
		var t time.Time
		if err := dec.Decode(&t); err != nil {
			break
		}
		times = append(times, t)
	}

	f.Truncate(0)

	return times

}

func main() {
	logger := log.Default()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	file, err := os.OpenFile(dataFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	c := decorators.NewPersisted(
		ephemeral.New(
			30*time.Second,
			ephemeral.WithContext(ctx),
			ephemeral.WithPeriodicVacuum(5*time.Second),
			ephemeral.WithLogger(logger),
			ephemeral.WithData(load(file)),
		),
		file,
	)

	done := make(chan struct{})
	go func() {
		if err := c.Start(ctx, done); err != nil {
			logger.Fatal(err)
		}
	}()

	http.HandleFunc("/", counterHandler(c))

	logger.Fatal(http.ListenAndServe(":8080", nil))

	<-done
}
