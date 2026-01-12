module github.com/Bitovi/example-go-server

go 1.25.5

require (
	github.com/bitovi-corp/auth-middleware-go v0.0.0
	github.com/google/uuid v1.6.0
)

replace github.com/bitovi-corp/auth-middleware-go => ../auth-middleware-go
