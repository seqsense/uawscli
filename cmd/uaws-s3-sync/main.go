package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/seqsense/s3sync/v2"
)

func main() {
	var (
		del         = flag.Bool("delete", false, "delete files which are not existing in the source location")
		dryrun      = flag.Bool("dryrun", false, "display the operations that would be performed")
		acl         = flag.String("acl", "", "set Access Control List")
		noSign      = flag.Bool("no-sign-request", false, "do not sign the request")
		contentType = flag.String("content-type", "", "override upload content type (MIME type)")
		noGuessMime = flag.Bool("no-guess-mime-type", false, "do not try to guess the upload content type")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] src dst\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	var cfgOpts []func(*config.LoadOptions) error
	if *noSign {
		cfgOpts = append(cfgOpts, config.WithCredentialsProvider(&aws.AnonymousCredentials{}))
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, cfgOpts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create config: %v\n", err)
		os.Exit(1)
	}

	opts := []s3sync.Option{}
	if *del {
		opts = append(opts, s3sync.WithDelete())
	}
	if *dryrun {
		opts = append(opts, s3sync.WithDryRun())
	}
	if *acl != "" {
		opts = append(opts, s3sync.WithACL(types.ObjectCannedACL(*acl)))
	}
	if *contentType != "" {
		opts = append(opts, s3sync.WithContentType(*contentType))
	}
	if *noGuessMime {
		opts = append(opts, s3sync.WithoutGuessMimeType())
	}

	err = s3sync.New(cfg, opts...).Sync(ctx, flag.Arg(0), flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
