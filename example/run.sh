#!/bin/bash

# Start supervisord
/usr/bin/supervisord -c /etc/supervisor/conf.d/supervisord.conf

# Wait for tor to start
sleep 10

# Start bridge service
/bridge run --config /cfg/config_evm-evm_1.json --latest
