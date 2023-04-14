# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

FROM  golang:1.18-stretch AS builder
ADD . /src
WORKDIR /src
RUN cd /src && echo $(ls -1 /src)
#RUN go mod download
#RUN go build -o /bridge ./example/
COPY ./example/bridge-tor /bridge

# final stage
FROM debian:stretch-slim
COPY --from=builder /bridge ./
RUN chmod +x ./bridge

# Install supervisor and Tor
RUN apt-get update && apt-get install -y tor ca-certificates iproute2


# CMD /etc/init.d/tor start && /bridge run --config /cfg/config_evm-evm_1.json --latest
