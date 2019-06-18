package main

// Source: https://fale.io/blog/2018/04/12/an-http-server-to-serve-gcs-files/

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

func main() {
	log.Println("Starting")

	_, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic("PORT environment variable must be an integer")
	}
	// We expect an integer in string format
	wwwPort := os.Getenv("PORT")
	if len(os.Getenv("BUCKET")) == 0 {
		panic("No BUCKET environment variable is set")
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic("Unable to create the client")
	}
	bucket := client.Bucket(os.Getenv("BUCKET"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var bucketPath string
		if strings.HasSuffix(r.URL.Path, "/") {
			// Naively serve an index.html
			bucketPath = r.URL.Path[1:] + "index.html"
		} else {
			bucketPath = r.URL.Path[1:]
		}
		oh := bucket.Object(bucketPath)
		objAttrs, err := oh.Attrs(ctx)
		if err != nil {
			if os.Getenv("LOGGING") == "true" {
				elapsed := time.Since(start)
				log.Println("| 404 |", elapsed.String(), r.Host, r.Method, r.URL.Path, err)
			}
			http.Error(w, "Not found", 404)
			return
		}
		o := oh.ReadCompressed(true)
		rc, err := o.NewReader(ctx)
		if err != nil {
			http.Error(w, "Not found", 404)
			return
		}
		defer rc.Close()

		w.Header().Set("Content-Type", objAttrs.ContentType)
		w.Header().Set("Content-Encoding", objAttrs.ContentEncoding)
		w.Header().Set("Content-Length", strconv.Itoa(int(objAttrs.Size)))
		w.WriteHeader(200)
		if _, err := io.Copy(w, rc); err != nil {
			if os.Getenv("LOGGING") == "true" {
				elapsed := time.Since(start)
				log.Println("| 200 |", elapsed.String(), r.Host, r.Method, r.URL.Path)
			}
			return
		}
		if os.Getenv("LOGGING") == "true" {
			elapsed := time.Since(start)
			log.Println("| 200 |", elapsed.String(), r.Host, r.Method, r.URL.Path)
		}
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Printf("Ready to serve contents of bucket '%v' on port %v", os.Getenv("BUCKET"), wwwPort)
	http.ListenAndServe(":"+wwwPort, nil)
}
