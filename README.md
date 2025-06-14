docker build -f cbm-api/Dockerfile -t cbm-api .

docker build --no-cache -f cbm-api/Dockerfile -t cbm-api .

Mongo shell login
docker exec -it 06f849d40a8c mongosh --username fudgebot --password cookiebot --authenticationDatabase admin

# How to regenerate the GraphQL resolvers

graph/schema - define objects in YOURFILENAME.graphqls

RUN 'go run github.com/99designs/gqlgen generate'

Thie creates templates in resolver/YOURFILENAME.resolver.go
inside these functions you call

database/YOURNAME.go (NOT GENERATED)

Fetch new SDL.gql file from Altair

Now in shared/graph add the sdl.gql file

create YOURNAME.graphql files

in here define the queries and mutaions you have tested in Altair

now you can use these funcitons in your code by inporting graph then
calling graph.QUERYNAME or MUTATIONNAME
Since your gqlgen.yml is already correctly set up inside cbm-api/,
you just need to run gqlgen.

From inside CBM-Backend/cbm-api/:

cd CBM-Backend/cbm-api
then simply run:

# find broken GraphQL queries quickly

grep -rnw . -e 'Pair' --include \*.graphql --include \*.go

docker compose down --volumes --remove-orphans
docker compose build --no-cache
docker compose up

|--  Purpose  -|- Directory -|- Tool/Command --|
|--  Generate client code (queries/mutations) for consuming GraphQL APIs   -|- /workspaces/CBM-Backend/shared/graph/   -|-  go run github.com/Khan/genqlient generate --|
|--  Generate server code (resolvers, models, schema stitching) for exposing GraphQL APIs   -|-  /workspaces/CBM-Backend/cbm-api/   -|-  go run github.com/99designs/gqlgen generate\ --|
