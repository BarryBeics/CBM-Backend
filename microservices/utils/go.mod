module cryptobotmanager.com/cbm-backend/microservices/utils

go 1.23

replace cryptobotmanager.com/cbm-backend/shared => ../../shared

require (
	cryptobotmanager.com/cbm-backend/shared v0.0.0-00010101000000-000000000000
	github.com/nats-io/nats.go v1.39.1
)

require (
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)
