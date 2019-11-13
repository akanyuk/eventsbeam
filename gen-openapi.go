package main

// go:generate go get -u github.com/go-swagger/go-swagger/cmd/swagger

//go:generate swagger generate spec -i web/generate/swagger-api.yaml -o web/swagger.json --scan-models
//go:generate go run web/generate/main.go -i web/swagger.json -o web/swagger.json.go
