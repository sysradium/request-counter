package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sysradium/request-counter/internal/counter"
	"github.com/sysradium/request-counter/internal/snapshot"
)

var (
	dataFilePath = "file.data"
)

func counterHandler(
	c *counter.SlidingWindowStorage,
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

func main() {
	logger := log.Default()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := counter.New(
		5*time.Second,
		counter.WithContext(ctx),
		counter.WithPeriodicVacuum(10*time.Second),
		counter.WithLogger(logger),
	)

	snap := snapshot.NewPeriodicSnapshotTaker(
		10*time.Second,
		c,
		func(b []byte) {
			tmpFile, _ := os.CreateTemp("", "tempfile")
			if _, err := tmpFile.Write(b); err != nil {
				return
			}

			if err := tmpFile.Close(); err != nil {
				return
			}

			if err := os.Rename(tmpFile.Name(), dataFilePath); err != nil {
				return
			}
		},
	)

	go snap.Run(ctx)

	go func() {
		if err := c.Start(); err != nil {
			logger.Fatal(err)
		}
	}()

	http.HandleFunc("/", counterHandler(c))

	logger.Fatal(http.ListenAndServe(":8080", nil))
}
