#!/usr/bin/env bash

go run $(dirname $0)/main.go < input.txt | grep "value 61" | grep "value 17"
