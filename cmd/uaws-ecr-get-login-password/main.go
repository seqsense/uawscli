package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

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
	userPass, err := base64.StdEncoding.DecodeString(*res.AuthorizationData[0].AuthorizationToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to decode authorization data: %v\n", err)
		return
	}

	pass := strings.Split(string(userPass), ":")
	if len(pass) != 2 {
		fmt.Fprint(os.Stderr, "error: invalid number of fields of user:pass is returned")
		return
	}
	fmt.Println(pass[1])
}
