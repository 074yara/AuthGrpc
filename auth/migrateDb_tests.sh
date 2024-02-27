#!/bin/bash

go run ./cmd/migrator \
--storage-path=./storage/auth.db \
--migrations-path=./tests/migrations \
--migrations-table=migrations_test