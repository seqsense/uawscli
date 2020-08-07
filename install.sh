#!/bin/sh

curl -sL https://raw.githubusercontent.com/seqsense/uawscli/master/uaws \
  | sh -s install $@
