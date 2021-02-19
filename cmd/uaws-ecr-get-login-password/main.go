package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load config: %v\n", err)
	}

	svc := ecr.NewFromConfig(cfg)

	res, err := svc.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get ecr login password: %v\n", err)
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
