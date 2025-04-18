module cryptobotmanager.com/cbm-backend/microservices/reports

go 1.24.2

replace cryptobotmanager.com/cbm-backend/shared => ../../shared

replace cryptobotmanager.com/cbm-backend/shared/graph => ../graph

require (
	cryptobotmanager.com/cbm-backend/shared v0.0.0-00010101000000-000000000000
	github.com/Khan/genqlient v0.8.0
	github.com/go-gota/gota v0.12.0
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/vektah/gqlparser/v2 v2.5.19 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	gonum.org/v1/gonum v0.9.1 // indirect
)
