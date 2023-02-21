#!/bin/sh

if [ $2 = "windows" ] && [ $3 = "arm64" ]; then
  # Not supported by upx
  exit 0
fi

exec upx $1
