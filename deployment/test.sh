#!/usr/bin/env bash

go test -v ./services/trace/ && \
  go test -v ./services/meander/ && \
  go test -v ./services/backup/ && \
  go test -v ./services/backup/*/
