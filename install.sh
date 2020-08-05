#!/bin/sh

set -eu

install_targets=
install_dir=~/.local/bin/
version=latest
while [ $# -gt 0 ]
do
  case $1 in
    -v)
      shift
      if [ $# -lt 1 ]
      then
        echo "error: -v must be followed by a version" >&2
        exit 1
      fi
      version=$1
      ;;
    -i)
      shift
      if [ $# -lt 1 ]
      then
        echo "error: -i must be followed by a path" >&2
        exit 1
      fi
      install_dir=$1
      ;;
    -*)
      echo "error: unknown option $1" >&2
      exit 1
      ;;
    *)
      install_targets="${install_targets} $1"
      ;;
  esac
  shift
done

if [ -z "${install_targets}" ]
then
  echo "usage: $0 [option...] install_targets" >&2
  echo "option:" >&2
  echo "  -v VERSION    install specific VERSION of uawscli" >&2
  echo "  -i DIR        install binaries under DIR (default: ${install_dir})" >&2
  exit 1
fi

arch=
case $(uname -p) in
  x86_64)
    arch=amd64
    ;;
  *)
    echo "error: unsupported arch" >&2
    exit 1
    ;;
esac

os=
ext=
case $(uname -s) in
  Linux)
    os=linux
    ;;
  *)
    echo "error: unsupported OS" >&2
    exit 1
    ;;
esac

api_auth=
if [ -n "${GITHUB_TOKEN:-${TRAVIS_BOT_GITHUB_TOKEN}}" ]
then
  api_auth="-H \"Authorization: token ${GITHUB_TOKEN:-${TRAVIS_BOT_GITHUB_TOKEN}}\""
fi

gh_api_base=${GITHUB_API_URL_BASE:-https://api.github.com}
slug=seqsense/uawscli

rel=
if [ ${version} = "latest" ]
then
  rel=$(eval curl \
    ${api_auth} \
    -s --retry 4 \
    ${gh_api_base}/repos/${slug}/releases/latest)
else
  rel=$(eval curl \
    ${api_auth} \
    -s --retry 4 \
    ${gh_api_base}/repos/${slug}/releases/tags/${version})
fi

urls=$(echo "${rel}" \
  | sed -n 's/.*"browser_download_url":\s*"\([^"]*\)"/\1/p' \
  | grep "_${os}_${arch}${ext}$")

downloads=
for target in ${install_targets}
do
  url=$(echo "${urls}" | grep "/uaws-${target}[^/]\+$" || true)

  if [ -n "$(echo ${url})" ]
  then
    downloads="${downloads} ${url}"
  else
    echo "error: ${target} not found" >&2
    exit 1
  fi
done

mkdir -p ${install_dir}

for dl in ${downloads}
do
  echo "downloading: ${url}"
  file=$(echo $url | sed -n "s|.*/\([^/]*\)_${os}_${arch}${ext}$|\1|p")
  curl -sL ${url} -o {install_dir}/${file}
done
