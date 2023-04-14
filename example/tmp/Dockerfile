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

# RUN echo "deb http://deb.debian.org/debian stretch main" > /etc/apt/sources.list && \
#  apt-get update
#RUN apt-get -y update
#RUN apt-get update && apt-get install -y --no-install-recommends apt-utils
#RUN apt-get update && apt-get install -y apt-transport-https
#RUN apt-get -y update && apt install iproute2
#RUN apt-get update && apt-get install -y ca-certificates
RUN apt-get update && apt-get install -y tor


# Start Tor service
# CMD ["tor"]
# CMD tor -f /etc/tor/torrc && tail -f /dev/null
# EXPOSE 9050

RUN chmod +x ./bridge
ENTRYPOINT ["./bridge"]

# Updated
# RUN echo "deb http://deb.debian.org/debian stretch main" > /etc/apt/sources.list && \
#     apt-get update
# RUN apt-get update && apt-get install -y ca-certificates tor supervisor

# Runint
# RUN apt-get update && apt-get install runit


# COPY ./example/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# # Expose Tor SOCKS port
# # EXPOSE 9050

# # Start supervisord
# CMD ["/usr/bin/supervisord", "-n"]

# ENTRYPOINT ["./bridge"]


# Update 2
# Add runit repository and update apt cache
#RUN echo "deb http://smarden.org/runit/debian stretch main" > /etc/apt/sources.list.d/runit.list && \
#   apt-get update

# Install packages
# RUN apt-get update && apt-get install -y ca-certificates tor supervisor runit-systemd

# # Create the tor user and group
# RUN groupadd -g 2000 tor && useradd -u 2000 -g tor -s /bin/false tor

# # Copy the runit service configuration file
# # COPY ./example/tor.runit /etc/sv/tor/run
# COPY ./example/tor.runit /etc/sv/tor/run

# # Create the tor data directory and set permissions
# RUN mkdir -p /var/lib/tor && chown -R tor:tor /var/lib/tor

# # Expose Tor SOCKS port
# # EXPOSE 9050

# # Start runit
# CMD ["/usr/bin/runsvdir", "/etc/service"]

# ENTRYPOINT ["./bridge"]


