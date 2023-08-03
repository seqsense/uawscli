package main

import (
	"strconv"
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
