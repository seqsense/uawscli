package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/seqsense/s3sync"
)

func main() {
	var (
		del    = flag.Bool("delete", false, "delete files which are not existing in the source location")
		dryrun = flag.Bool("dryrun", false, "display the operations that would be performed")
		noSign = flag.Bool("no-sign-request", false, "do not sign the request")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] src dst\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *del {
		fmt.Fprintf(os.Stderr, "error: -delete is not yet implemented\n")
		os.Exit(1)
	}
	if *dryrun {
		fmt.Fprintf(os.Stderr, "error: -dryrun is not yet implemented\n")
		os.Exit(1)
	}

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	var cfg []*aws.Config
	if *noSign {
		cfg = []*aws.Config{{Credentials: credentials.AnonymousCredentials}}
	}
	sess, err := session.NewSession(cfg...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create session: %v\n", err)
		os.Exit(1)
	}

	err = s3sync.New(sess).Sync(flag.Arg(0), flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
