#!/bin/sh

command=$1
shift

case "$command" in
  "echo")
    echo "$1"
    ;;
esac
