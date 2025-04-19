module cryptobotmanager.com/cbm-backend/microservices/backTesting

go 1.24.2

replace cryptobotmanager.com/cbm-backend/shared => ../../shared

replace cryptobotmanager.com/cbm-backend/microservices/filters => ../filters

replace cryptobotmanager.com/cbm-backend/microservices/reports => ../reports

replace cryptobotmanager.com/cbm-backend/microservices/tradingBots => ../tradingBots

replace cryptobotmanager.com/cbm-backend/shared/graph => ../graph

replace cryptobotmanager.com/cbm-backend/resolvers => ../../resolvers

require (
	cryptobotmanager.com/cbm-backend/microservices/filters v0.0.0-00010101000000-000000000000
	cryptobotmanager.com/cbm-backend/microservices/reports v0.0.0-00010101000000-000000000000
	cryptobotmanager.com/cbm-backend/microservices/tradingBots v0.0.0-00010101000000-000000000000
	cryptobotmanager.com/cbm-backend/resolvers v0.0.0-00010101000000-000000000000
	cryptobotmanager.com/cbm-backend/shared v0.0.0-00010101000000-000000000000
	github.com/Khan/genqlient v0.8.0
	github.com/adshao/go-binance/v2 v2.8.2
	github.com/joho/godotenv v1.5.1
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/go-gota/gota v0.12.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/vektah/gqlparser/v2 v2.5.19 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	gonum.org/v1/gonum v0.9.1 // indirect
)
