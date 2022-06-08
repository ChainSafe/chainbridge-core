# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

FROM  golang:1.18-stretch AS builder
ADD . /src
WORKDIR /src
RUN cd /src && echo $(ls -1 /src)
RUN go mod download
RUN go build -o /bridge ./example/

# final stage
FROM debian:stretch-slim
COPY --from=builder /bridge ./
RUN chmod +x ./bridge

ENTRYPOINT ["./bridge"]
