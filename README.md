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

POST Operation ( Add an item )
```
curl -X POST http://localhost:8080/items -H "Content-Type: application/json" -d '{"name": "Sample Item"}'
```

PUT Operation ( Update an item )
```
curl -X PUT http://localhost:8080/items/update -H "Content-Type: application/json" -d '{"id":1,"name": "Another new Sample Item"}'
```

DELETE Operation ( Delete an item )
```
curl -X DELETE http://localhost:8080/items/delete -H "Content-Type: application/json" -d '{"id":1,name": "Another Sample Item"}'
```