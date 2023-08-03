package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Bool struct {
	value *bool
}

func (b *Bool) String() string {
	if b.value == nil {
		return "nil"
	}
	return strconv.FormatBool(bool(*b.value))
}

func (b *Bool) IsBoolFlag() bool {
	return true
}

func (b *Bool) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	b.value = &v
	return nil
}

type BoolInv struct {
	*Bool
}

func (b *BoolInv) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	v = !v
	b.value = &v
	return nil
}

func main() {
	var name = flag.String("name", "", "parameter name")
	var withDecription Bool

	flag.Var(&withDecription, "with-decryption", "return decrypted secure string value")
	flag.Var(&BoolInv{Bool: &withDecription}, "no-with-decryption", "not return decrypted secure string value")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *name == "" {
		fmt.Fprintf(os.Stderr, "error: --name is required\n")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load config: %v\n", err)
	}

	svc := ssm.NewFromConfig(cfg)

	res, err := svc.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           name,
		WithDecryption: withDecription.value,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get parameter: %v\n", err)
		return
	}
	out, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to encode response: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", string(out))
}
