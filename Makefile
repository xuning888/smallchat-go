.PHONY: build
build:
	@rm -rf ./out
	@mkdir -p ./out/bin
	@go build -o smallchat .
	@mv ./smallchat ./out/bin/smallchat