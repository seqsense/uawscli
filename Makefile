.PHONY: build
build:
	goreleaser build --snapshot

.PHONY: clean
clean:
	rm -rf dist
