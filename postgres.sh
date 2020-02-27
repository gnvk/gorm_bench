#!/usr/bin/env sh

docker run -d -p 35432:5432 --name pgworld ghusta/postgres-world-db:2.4
