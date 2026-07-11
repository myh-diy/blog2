#!/bin/sh
set -eu

./mini-node-exporter &
exec ./blog-server
