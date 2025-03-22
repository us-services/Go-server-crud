# Go-server-crud
Go server for crud operations

## Run Go server
1. ensure that go 1.24 is installed
2. clone this repository
3. install dependencies using go mod tidy
4. run go run main.go

## Operations
GET Operation
```
curl -X GET http://localhost:8080/items | jq .
```

POST Operation
```
curl -X POST http://localhost:8080/items -H "Content-Type: application/json" -d '{"name": "Sample Item"}'
```