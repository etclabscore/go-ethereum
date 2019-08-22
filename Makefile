.DEFAULT_GOAL := build

BUILD_VERSION?=snapshot
SOURCE_FILES?=$$(go list ./... | grep -v /vendor/)
TEST_PATTERN?=.
TEST_OPTIONS?=-race

GO_MOD=GO111MODULE=on
LDFLAGS=-ldflags "-X main.Version="`git describe --tags`

BINARY=bin
BUILD_TIME=`date +%FT%T%z`
COMMIT=`git log --pretty=format:'%h' -n 1`

# Provide default value of GOPATH, if it's not set in environment
export GOPATH?=${HOME}/go

build: cmd/abigen cmd/bootnode cmd/disasm cmd/ethtest cmd/evm cmd/gethrpctest cmd/rlpdump cmd/geth ## Build a local snapshot binary version of all commands
	@ls -ld $(BINARY)/*

win/build: win/abigen win/bootnode win/disasm win/ethtest win/evm win/gethrpctest win/rlpdump win/geth ## Build a local snapshot binary version of all commands
	@ls -ld $(BINARY)/*

cmd/geth: chainconfig ## Build a local snapshot binary version of geth.
	mkdir -p ./${BINARY}
	${GO_MOD} go build ${LDFLAGS} -o ${BINARY}/geth -tags="netgo" ./cmd/geth
	@echo "Done building geth."
	@echo "Run \"$(BINARY)/geth\" to launch geth."

cmd/abigen: ## Build a local snapshot binary version of abigen.
	mkdir -p ./${BINARY} && ${GO_MOD} go build ${LDFLAGS} -o ${BINARY}/abigen ./cmd/abigen
	@echo "Done building abigen."
	@echo "Run \"$(BINARY)/abigen\" to launch abigen."

cmd/bootnode: ## Build a local snapshot of bootnode.
	mkdir -p ./${BINARY} && ${GO_MOD} go build ${LDFLAGS} -o ${BINARY}/bootnode ./cmd/bootnode
	@echo "Done building bootnode."
	@echo "Run \"$(BINARY)/bootnode\" to launch bootnode."

cmd/disasm: ## Build a local snapshot of disasm.
	mkdir -p ./${BINARY} && ${GO_MOD} go build ${LDFLAGS} -o ${BINARY}/disasm ./cmd/disasm
	@echo "Done building disasm."
	@echo "Run \"$(BINARY)/disasm\" to launch disasm."

cmd/ethtest: ## Build a local snapshot of ethtest.
	mkdir -p ./${BINARY} && ${GO_MOD} go build ${LDFLAGS} -o ${BINARY}/ethtest ./cmd/ethtest
	@echo "Done building ethtest."
	@echo "Run \"$(BINARY)/ethtest\" to launch ethtest."

cmd/evm: ## Build a local snapshot of evm.
	mkdir -p ./${BINARY} && ${GO_MOD} go build ${LDFLAGS} -o ${BINARY}/evm ./cmd/evm
	@echo "Done building evm."
	@echo "Run \"$(BINARY)/evm\" to launch evm."

cmd/gethrpctest: ## Build a local snapshot of gethrpctest.
	mkdir -p ./${BINARY} && ${GO_MOD} go build ${LDFLAGS} -o ${BINARY}/gethrpctest ./cmd/gethrpctest
	@echo "Done building gethrpctest."
	@echo "Run \"$(BINARY)/gethrpctest\" to launch gethrpctest."

cmd/rlpdump: ## Build a local snapshot of rlpdump.
	mkdir -p ./${BINARY} && ${GO_MOD} go build ${LDFLAGS} -o ${BINARY}/rlpdump ./cmd/rlpdump
	@echo "Done building rlpdump."
	@echo "Run \"$(BINARY)/rlpdump\" to launch rlpdump."

win/geth: chainconfig ## Build a local snapshot binary version of geth.
	mkdir -p ./${BINARY}
	${GO_MOD} GOOS=windows go build ${LDFLAGS} -o ${BINARY}/geth.exe -tags="netgo" ./cmd/geth
	@echo "Done building geth."
	@echo "Run \"$(BINARY)/geth\" to launch geth."

win/abigen: ## Build a local snapshot binary version of abigen.
	mkdir -p ./${BINARY} && ${GO_MOD} GOOS=windows go build ${LDFLAGS} -o ${BINARY}/abigen.exe ./cmd/abigen
	@echo "Done building abigen."
	@echo "Run \"$(BINARY)/abigen\" to launch abigen."

win/bootnode: ## Build a local snapshot of bootnode.
	mkdir -p ./${BINARY} && ${GO_MOD} GOOS=windows go build ${LDFLAGS} -o ${BINARY}/bootnode.exe ./cmd/bootnode
	@echo "Done building bootnode."
	@echo "Run \"$(BINARY)/bootnode\" to launch bootnode."

win/disasm: ## Build a local snapshot of disasm.
	mkdir -p ./${BINARY} && ${GO_MOD} GOOS=windows go build ${LDFLAGS} -o ${BINARY}/disasm.exe ./cmd/disasm
	@echo "Done building disasm."
	@echo "Run \"$(BINARY)/disasm\" to launch disasm."

win/ethtest: ## Build a local snapshot of ethtest.
	mkdir -p ./${BINARY} && ${GO_MOD} GOOS=windows go build ${LDFLAGS} -o ${BINARY}/ethtest.exe ./cmd/ethtest
	@echo "Done building ethtest."
	@echo "Run \"$(BINARY)/ethtest\" to launch ethtest."

win/evm: ## Build a local snapshot of evm.
	mkdir -p ./${BINARY} && ${GO_MOD} GOOS=windows go build ${LDFLAGS} -o ${BINARY}/evm.exe ./cmd/evm
	@echo "Done building evm."
	@echo "Run \"$(BINARY)/evm\" to launch evm."

win/gethrpctest: ## Build a local snapshot of gethrpctest.
	mkdir -p ./${BINARY} && ${GO_MOD} GOOS=windows go build ${LDFLAGS} -o ${BINARY}/gethrpctest.exe ./cmd/gethrpctest
	@echo "Done building gethrpctest."
	@echo "Run \"$(BINARY)/gethrpctest\" to launch gethrpctest."

win/rlpdump: ## Build a local snapshot of rlpdump.
	mkdir -p ./${BINARY} && ${GO_MOD} GOOS=windows go build ${LDFLAGS} -o ${BINARY}/rlpdump.exe ./cmd/rlpdump
	@echo "Done building rlpdump."
	@echo "Run \"$(BINARY)/rlpdump\" to launch rlpdump."

install: ## Install all packages to $GOPATH/bin
	${GO_MOD} go install ./cmd/{abigen,bootnode,disasm,ethtest,evm,gethrpctest,rlpdump}
	$(MAKE) install_geth

install_geth: chainconfig ## Install geth to $GOPATH/bin
	$(info Installing $$GOPATH/bin/geth)
	${GO_MOD} go install ${LDFLAGS} -tags="netgo" ./cmd/geth

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' -not -wholename './_vendor*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

ci: lint test ## Run all code checks and tests

lint: ## Run all the linters
	gometalinter \
		--tests \
		--vendor \
		--vendored-linters \
		--disable=interfacer \
		--disable=maligned \
		--enable=gosimple \
		--enable=staticcheck \
		--enable=gofmt \
		--enable=goimports \
		--enable=lll \
		--enable=misspell \
		--cyclo-over=15 \
		--dupl-threshold=100 \
		--line-length=120 \
		--deadline=360s \
		./...

test: ## Run all the tests
	echo 'mode: atomic' > coverage.txt && \
	${GO_MOD} go list ./... | xargs -n1 -I{} sh -c '${GO_MOD} go test -covermode=atomic -coverprofile=coverage.tmp {} && \
	tail -n +2 coverage.tmp >> coverage.txt' && \
	rm coverage.tmp

cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

chainconfig: core/assets/assets.go ## Rebuild assets if source config files changed.

core/assets/assets.go: ${GOPATH}/bin/resources core/config/*.json core/config/*.csv
	${GOPATH}/bin/resources -fmt -declare -var=DEFAULTS -package=assets -output=core/assets/assets.go core/config/*.json core/config/*.csv

${GOPATH}/bin/resources:
	${GO_MOD} go get -u github.com/omeid/go-resources/cmd/resources

clean: ## Remove local snapshot binary directory
	if [ -d ${BINARY} ] ; then rm -rf ${BINARY} ; fi
	go clean -i ./...
	rm -rf vendor/github.com/ETCDEVTeam/sputnikvm-ffi/c/ffi/target
	rm -f vendor/github.com/ETCDEVTeam/sputnikvm-ffi/c/libsputnikvm.*

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: setup test cover fmt lint ci build cmd/geth cmd/abigen cmd/bootnode cmd/disasm cmd/ethtest cmd/evm cmd/gethrlptest cmd/rlpdump install install_geth clean help static
