#!/bin/bash
# Start the Next.js server
docker run -d --name listmonk_connector -e MYSQL_ROOT_PASSWORD=password -p 3307:3306 mariadb:latest