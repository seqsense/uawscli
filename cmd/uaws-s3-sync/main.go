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
		acl    = flag.String("acl", "", "set Access Control List")
		noSign = flag.Bool("no-sign-request", false, "do not sign the request")
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

	var cfg []*aws.Config
	if *noSign {
		cfg = []*aws.Config{{Credentials: credentials.AnonymousCredentials}}
	}
	sess, err := session.NewSession(cfg...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create session: %v\n", err)
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
		opts = append(opts, s3sync.WithACL(*acl))
	}

	err = s3sync.New(sess, opts...).Sync(flag.Arg(0), flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
