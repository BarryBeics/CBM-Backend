module cryptobotmanager.com/cbm-backend/microservices/priceData

go 1.23

replace cryptobotmanager.com/cbm-backend/shared => ../../shared

replace cryptobotmanager.com/cbm-backend/shared/graph => ../graph

require (
	cryptobotmanager.com/cbm-backend/shared v0.0.0-00010101000000-000000000000
	github.com/Khan/genqlient v0.7.0
	github.com/adshao/go-binance/v2 v2.7.0
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/vektah/gqlparser/v2 v2.5.11 // indirect
	golang.org/x/sys v0.28.0 // indirect
)
