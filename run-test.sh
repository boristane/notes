#!/bin/bash
set -e
pipeline=${PIPELINE_LABEL:-local}
# docker-compose down

docker-compose up -d

# Check that mysql is ready before creating the database
# while ! docker exec mysql_users mysqladmin --user=root --password=password --host "127.0.0.1" ping --silent &> /dev/null ; do
#   sleep 2
# done

# SNS Topic for emails

docker exec -t mysql mysql -h 127.0.0.1 -u root -ppassword -e "create database notes" || true

# docker-compose down