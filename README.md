# Installation

1. Start all services using docker-compose:

   ```
   docker-compose up -d
   ```
   
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
   
   View table with IXP server data receiving live updates.


4. Stop/Remove all services:

   ```
   docker-compose stop/down
   ```



# Microservices

## parser
A PCH parser to fetch store and present BGP summary updates.

Configuration file is mounted as a volume. Updates are live loaded as the configuration
is loaded every 10 seconds that parser binary is executed.

## postgres
A Postgres instance to store the BGP summaries updated by the parser

## hasura
GraphQL engine to access postgres and serve data realtime

## ui
A React application using websockets to fetch realtime BGP summaries and update an HTML table