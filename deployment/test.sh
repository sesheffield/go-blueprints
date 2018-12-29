#!/usr/bin/env bash

go test -v ./services/trace/ && \
  go test -v ./services/meander/
