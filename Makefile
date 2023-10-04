
.PHONY: help run build install license example
all: help

## license: Adds license header to missing files.
license:
	@echo "  >  \033[32mAdding license headers...\033[0m "
	GO111MODULE=off go get -u github.com/google/addlicense
	addlicense -c "ChainSafe Systems" -f ./scripts/header.txt -y 2021 .

## license-check: Checks for missing license headers
license-check:
	@echo "  >  \033[Checking for license headers...\033[0m "
	GO111MODULE=off go get -u github.com/google/addlicense
	addlicense -check -c "ChainSafe Systems" -f ./scripts/header.txt -y 2021 .


coverage:
	go tool cover -func cover.out | grep total | awk '{print $3}'

test:
	./scripts/tests.sh

e2e-test:
	./scripts/e2e_tests.sh

example:
	docker-compose --file=./example/docker-compose.yml up

## Install dependency subkey
install-subkey:
	curl https://getsubstrate.io -sSf | bash -s -- --fast
	cargo install --force --git https://github.com/paritytech/substrate subkey

genmocks:
	mockgen -destination=./mock/client.go -source=./chains/evm/client/client.go -package mock
	mockgen -destination=./mock/gas.go -source=./chains/evm/transactor/gas/gas-pricer.go -package mock
	mockgen -destination=./mock/relayer.go -source=./relayer/relayer.go -package mock
	mockgen -source=chains/evm/transactor/transact.go -destination=./mock/transact.go -package mock
	mockgen -source=chains/evm/transactor/signAndSend/signAndSend.go -destination=./mock/signAndSend.go -package mock
	mockgen -source=./store/store.go -destination=./mock/store.go -package mock
