# µawscli - collection of tiny AWS console client utilities

## Install helper script
```shell
# Install helper script
$ curl -sL https://raw.githubusercontent.com/seqsense/uawscli/master/uaws \
  -o ~/.local/bin/uaws
$ chmod +x ~/.local/bin/uaws
```

```shell
$ uaws ecr get-login-password
```
will automatically download uaws-ecr-get-login-password under `~/.local/lib/uaws` and execute.

## Install specific utility (without installing helper script)
```shell
# Install latest version of uaws-ecr-get-login-password under ~/.local/bin
$ curl -sL https://raw.githubusercontent.com/seqsense/uawscli/master/uaws \
  | sh -s install ecr-get-login-password

# Install latest version of uaws-ecr-get-login-password
# under /path/prefix/bin and /path/prefix/lib/uawscli
$ curl -sL https://raw.githubusercontent.com/seqsense/uawscli/master/uaws \
  | sh -s install -i /path/prefix ecr-get-login-password

# Install specific version of uaws-ecr-get-login-password
# under /path/prefix/bin and /path/prefix/lib/uawscli
$ curl -sL https://raw.githubusercontent.com/seqsense/uawscli/master/uaws \
  | sh -s install -v v0.0.0 -i /path/prefix ecr-get-login-password
```

```shell
$ uaws-ecr-get-login-password
```
will be available.
