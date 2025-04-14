module cryptobotmanager.com/cbm-backend/shared

go 1.23

require (
	cryptobotmanager.com/cbm-backend/microservices/backTesting v0.0.0-00010101000000-000000000000
	github.com/Khan/genqlient v0.8.0
	github.com/nats-io/nats.go v1.39.1
	github.com/rs/zerolog v1.34.0
)

replace cryptobotmanager.com/cbm-backend/microservices/backTesting => ../microservices/backTesting

replace cryptobotmanager.com/cbm-backend/microservices/filters => ../microservices/filters

replace cryptobotmanager.com/cbm-backend/resolvers => ../resolvers

require (
	cryptobotmanager.com/cbm-backend/microservices/filters v0.0.0-00010101000000-000000000000 // indirect
	cryptobotmanager.com/cbm-backend/resolvers v0.0.0-00010101000000-000000000000 // indirect
	github.com/adshao/go-binance/v2 v2.8.2 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/vektah/gqlparser/v2 v2.5.19 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)
