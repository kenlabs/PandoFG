#!/usr/bin/env bash

if [ -f "$HOME/.pando/config.yaml" ]; then
  /opt/pando-server daemon
else
  /opt/pando-server init && /opt/pando-server daemon
fi