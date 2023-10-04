#!/bin/bash

# Read server port from command line arguments or use default value
SERVER_PORT=${2:-5001}

if [ "$1" == "server" ]; then
    exec /app/bin/server -port="$SERVER_PORT"
elif [ "$1" == "client" ]; then
    # Note: You can also modify the next line to accept client-specific arguments
    exec /app/bin/client -sPorts="${SERVER_PORT}"
else
    echo "Please specify 'client' or 'server' as the first argument."
    exit 1
fi
