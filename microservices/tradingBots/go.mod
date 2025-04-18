module cryptobotmanager.com/cbm-backend/microservices/tradingBots

go 1.24.2

replace cryptobotmanager.com/cbm-backend/shared => ../../shared

replace cryptobotmanager.com/cbm-backend/shared/graph => ../graph

replace cryptobotmanager.com/cbm-backend/resolvers => ../../resolvers

require (
	cryptobotmanager.com/cbm-backend/resolvers v0.0.0-00010101000000-000000000000
	cryptobotmanager.com/cbm-backend/shared v0.0.0-00010101000000-000000000000
	github.com/Khan/genqlient v0.8.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/vektah/gqlparser/v2 v2.5.19 // indirect
)
