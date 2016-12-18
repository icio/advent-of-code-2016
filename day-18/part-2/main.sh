#!/usr/bin/env bash

go run $(dirname $0)/../part-1/main.go 400000 < $(dirname $0)/../part-1/input.txt 2>/dev/null | tee README
