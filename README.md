docker build -f cbm-api/Dockerfile -t cbm-api .

docker build --no-cache -f cbm-api/Dockerfile -t cbm-api .

 How to regenerate the GraphQL resolvers
Since your gqlgen.yml is already correctly set up inside cbm-api/,
you just need to run gqlgen.

From inside CBM-Backend/cbm-api/:

cd CBM-Backend/cbm-api
then simply run:

go run github.com/99designs/gqlgen generate

# find broken GraphQL queries quickly

grep -rnw . -e 'Pair' --include \*.graphql --include \*.go

docker compose down --volumes --remove-orphans
docker compose build --no-cache
docker compose up

|--  Purpose  -|- Directory -|- Tool/Command --|
|--  Generate client code (queries/mutations) for consuming GraphQL APIs   -|- /workspaces/CBM-Backend/shared/graph/   -|-  go run github.com/Khan/genqlient generate --|
|--  Generate server code (resolvers, models, schema stitching) for exposing GraphQL APIs   -|-  /workspaces/CBM-Backend/cbm-api/   -|-  go run github.com/99designs/gqlgen generate\ --|
