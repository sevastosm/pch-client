# pch-parser
A PCH parser to fetch store and present BGP summary update

Config file (./parser/config/parser.yml) is mount as a volume,
so any config updates are loaded in the next run of 
the parser without stopping the docker container.

# postgres
A Postgres instance to store the BGP summaries updated by the parser

# hasura
GraphQL engine to access postgres and serve data realtime
Web URL: http://localhost:8001

# ui
A react application using websockets to fetch realtime BGP summaries and update an HTML table
Web URL: http://localhost:8000
