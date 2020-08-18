package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

func main() {
	var (
		distributionID = flag.String("distribution-id", "", "distribution ID")
		paths          = flag.String("paths", "", "space-separated paths to be invalidated")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *distributionID == "" {
		fmt.Fprintf(os.Stderr, "error: --distribution-id is required\n")
		os.Exit(1)
	}

	pathsSlice := strings.Fields(*paths)
	if len(pathsSlice) == 0 {
		fmt.Fprintf(os.Stderr, "error: --paths is required\n")
		os.Exit(1)
	}

	if flag.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "error: unused arguments: %v\n", flag.Args())
		fmt.Fprintf(os.Stderr, "       uawscli requires that value of --paths argument be quated\n")
		os.Exit(1)
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load config: %v\n", err)
	}

	svc := cloudfront.New(cfg)

	req := svc.CreateInvalidationRequest(&cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(*distributionID),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(time.Now().UTC().Format("20060102150405")),
			Paths: &cloudfront.Paths{
				Items:    pathsSlice,
				Quantity: aws.Int64(int64(len(pathsSlice))),
			},
		},
	})
	res, err := req.Send(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create invalidation: %v\n", err)
		return
	}
	fmt.Printf("%v\n", *res)
}
