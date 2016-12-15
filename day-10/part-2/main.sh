#!/usr/bin/env bash

go run $(dirname $0)/../part-1/main.go < $(dirname $0)/../part-1/input.txt | grep -E 'Output [012] received' | tee >( cat >&2) | awk '{ print $5 }' | python -c 'import sys; print reduce(lambda a, b: a*b, (int(n.rstrip()) for n in iter(sys.stdin)), 1)'
