module cryptobotmanager.com/cbm-backend/shared

go 1.24.2

require (
	cryptobotmanager.com/cbm-backend/cbm-api v0.0.0-00010101000000-000000000000
	cryptobotmanager.com/cbm-backend/microservices/backTesting v0.0.0-00010101000000-000000000000
	github.com/Khan/genqlient v0.8.0
	github.com/nats-io/nats.go v1.39.1
	github.com/rs/zerolog v1.34.0
)

replace cryptobotmanager.com/cbm-backend/microservices/backTesting => ../microservices/backTesting

replace cryptobotmanager.com/cbm-backend/microservices/filters => ../microservices/filters

replace cryptobotmanager.com/cbm-backend/cbm-api => ../cbm-api

replace cryptobotmanager.com/cbm-backend/microservices/binance => ../microservices/binance

replace cryptobotmanager.com/cbm-backend/microservices/reports => ../microservices/reports

replace cryptobotmanager.com/cbm-backend/microservices/tradingBots => ../microservices/tradingBots

require (
	cryptobotmanager.com/cbm-backend/microservices/binance v0.0.0-00010101000000-000000000000 // indirect
	cryptobotmanager.com/cbm-backend/microservices/filters v0.0.0-00010101000000-000000000000 // indirect
	cryptobotmanager.com/cbm-backend/microservices/reports v0.0.0-00010101000000-000000000000 // indirect
	cryptobotmanager.com/cbm-backend/microservices/tradingBots v0.0.0-00010101000000-000000000000 // indirect
	github.com/go-gota/gota v0.12.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/vektah/gqlparser/v2 v2.5.25 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	gonum.org/v1/gonum v0.9.1 // indirect
)
