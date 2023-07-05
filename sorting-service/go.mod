module github.com/sandio/sort/sorting-service

go 1.16

replace github.com/sandio/sort/gen => ../gen

require (
	github.com/sandio/sort/gen v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.53.0
)
