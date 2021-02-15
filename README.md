# How it works
Postgres is used for storing IXP monitoring data into the `IXP_SERVER_DATA` table. Database instance
is populated with DDL for creating the table at start-up.

Hasura is used for exposing the data from Postgres using GraphQL. The Hasura instance is configured
using `metadata` in order to monitor the `IXP_SERVER_DATA` table.

Parser client is used for updating database with IXP monitoring data from PCH. A bash script ran as
docker container entry point initiates the parser client every 10 seconds. The client fetches and 
updates records in the `IXP_SERVER_DATA` table. The parser is configured by default a YAML file
`./parser/config/parser.yml`

UI service contains a React app which is served using nginx. The JS application is using Apollo 
GraphQL library to connect to Hasura and update UI with realtime data.


# Installation and Running

1. Start all services using docker-compose:

   ```
   docker-compose up -d
   ```
   This step is expected to take long due to downloading the docker images of postgres, hasura 
   and base images for golang and nginx on which the parser and ui service depend. This step will
   build the `parser` docker image by building the golang app parser that was implemented using 
   golang and the `ui` docker image by attaching the React app artifact into nginx image.
   

2. Configure parser (./parser/config/parser.yml):

   ```
   parser:
      server_limit: 50
      protocol: IPv6                      # IPv4 or IPv6
      ixp:
      - DE-CIX Marseille
       city:
      - Marseille
      country:
      - France
   ```
   
3. Open UI URL:
   ```
   http://localhost:8000/   
   ```
   
   HTML & JS app presenting table with IXP server data receiving live updates. 
   The table presents all IXP servers data stored in the database and receives updates 
   from the database using WebSocket and Hasura GraphQL subscription.


4. Stop all services:

   ```
   docker-compose stop/down
   ```
   
5. Reset:

   Delete database files.
   ```
   rm -r ./db/data/*
   ```
