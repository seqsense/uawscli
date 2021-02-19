package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cf_types "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load config: %v\n", err)
	}

	svc := cloudfront.NewFromConfig(cfg)

	res, err := svc.CreateInvalidation(ctx, &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(*distributionID),
		InvalidationBatch: &cf_types.InvalidationBatch{
			CallerReference: aws.String(time.Now().UTC().Format("20060102150405")),
			Paths: &cf_types.Paths{
				Items:    pathsSlice,
				Quantity: aws.Int32(int32(len(pathsSlice))),
			},
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create invalidation: %v\n", err)
		return
	}
	fmt.Printf("%v\n", *res)
}
