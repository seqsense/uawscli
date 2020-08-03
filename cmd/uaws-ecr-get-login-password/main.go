package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load config: %v\n", err)
	}

	svc := ecr.New(cfg)

	req := svc.GetAuthorizationTokenRequest(&ecr.GetAuthorizationTokenInput{})
	res, err := req.Send(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get ect login password: %v\n", err)
		return
	}

	if len(res.AuthorizationData) == 0 {
		fmt.Fprint(os.Stderr, "error: no authorization data is returned\n")
		return
	}
	fmt.Print(*res.AuthorizationData[0].AuthorizationToken)
}
