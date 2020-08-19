package main

// go:generate go get -u github.com/go-swagger/go-swagger/cmd/swagger
//go:generate swagger generate spec -i web/generate/swagger-api.yaml -o swagger.json --scan-models
//go:generate swagger2openapi swagger.json -o openapi.yaml
//go:generate rm -f swagger.json
//go:generate go run web/generate/main.go -i openapi.yaml -o web/openapi.yaml.go
