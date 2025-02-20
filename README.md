# receipt-challenge



### Instructions 

Run `go run main.go`


#### Test receipts/process:
```
curl --request POST \
  --url http://localhost:8080/receipts/process \
  --header 'Content-Type: application/json' \
  --data '{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}'
```
This will return the id

Example:
```
{"id":"73eb91f9-518e-4425-aa9e-46c71e63d095"}
```


#### Test /receipts/{id}/points:

Use the newly generated `id` to receive the generated points

```
curl --request GET \
  --url http://localhost:8080/receipts/73eb91f9-518e-4425-aa9e-46c71e63d095/points
  ```

  