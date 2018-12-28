#!/usr/bin/env bash

go vet ./services/chat/ && \
  go vet ./services/cli-tool/*/ && \
  go vet ./services/mq/*/ && \
  go vet ./services/trace/

