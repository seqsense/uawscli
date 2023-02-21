.PHONY: build
build:
	goreleaser build --snapshot

.PHONY: clean
clean:
	rm -rf dist

.PHONY: install-upx
install-upx:
	gh release download --repo=upx/upx --pattern="upx-*-amd64_linux.tar.xz" --output=- \
		| tar -C $(HOME)/.local/bin/ -xJf - --strip-components=1 --wildcards '*/upx'
	chmod +x $(HOME)/.local/bin/upx
