# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

FROM  golang:1.17-stretch AS builder

ADD . /src
WORKDIR /src
RUN cd /src && echo $(ls -1 /src)
RUN go mod download
RUN go build -o /bridge ./e2e/evm-optimism/example/.

# final stage
FROM debian:stretch-slim

RUN apt-get update -y && apt-get install -y curl

COPY --from=builder /bridge ./
RUN chmod +x ./bridge

COPY ./e2e/evm-optimism/scripts/relayer.sh .
RUN chmod +x relayer.sh

ENTRYPOINT ["./bridge"]