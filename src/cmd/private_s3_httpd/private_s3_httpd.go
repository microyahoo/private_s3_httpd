package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awscred "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/handlers"
)

func main() {
	showVersion := flag.Bool("version", false, "print version string")
	listen := flag.String("listen", ":8444", "address:port to listen on.")
	bucket := flag.String("bucket", "", "S3 bucket name")
	logRequests := flag.Bool("log-requests", true, "log HTTP requests")
	region := flag.String("region", "us-east-1", "S3 Region")
	ak := flag.String("ak", "", "S3 access key")
	sk := flag.String("sk", "", "S3 secret key")
	s3Endpoint := flag.String("s3-endpoint", "", "S3 endpoint")
	flag.Parse()

	if *showVersion {
		fmt.Printf("private_s3_httpd v%s (built w/%s)\n", VERSION, runtime.Version())
		return
	}

	if *bucket == "" {
		log.Fatalf("bucket name required")
	}
	if *s3Endpoint == "" {
		log.Fatal("s3 endpoint required")
	}
	if *ak == "" || *sk == "" {
		log.Fatalf("ak and sk required")
	}

	log.Printf("Using alternate S3 Endpoint diwht DisableSSL:true, S3ForcePathStyle:true %q", *s3Endpoint)
	creds := awscred.NewStaticCredentials(*ak, *sk, "")
	sess, err := session.NewSession(&aws.Config{
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		Region:           region,
		Endpoint:         s3Endpoint,
		Credentials:      creds,
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("Failed to create s3 session: %s", err)
	}
	svc := s3.New(sess)

	var h http.Handler
	h = &Proxy{
		Bucket: *bucket,
		Svc:    svc,
	}
	if *logRequests {
		h = handlers.CombinedLoggingHandler(os.Stdout, h)
	}

	s := &http.Server{
		Addr:           *listen,
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("listening on %s", *listen)
	log.Fatal(s.ListenAndServe())

}
