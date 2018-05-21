#!/bin/sh

export MESSENGER_ADDRESS=$HOST
export MESSENGER_PORT=$PORT1

echo "ok"
env

ls -al /opt/eremetic/eremetic
exec /opt/eremetic/eremetic
