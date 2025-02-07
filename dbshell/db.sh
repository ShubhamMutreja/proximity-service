#!/bin/bash
set -e
export PGPASSWORD=postgres123;
psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "businessdata" <<-EOSQL
  CREATE DATABASE businessdata;
  GRANT ALL PRIVILEGES ON DATABASE businessdata TO "postgres";
EOSQL