FAKE API
========

How to set run:

```go run main.go --config config.yml```

Register:

```
curl -X POST \
     http://localhost:24213 \
     -d '{
     "id": "1",
     "jsonrpc": "2.0",
     "method": "magento.service_bus.remote.register",
     "params": {
         "id": "test",
         "url": "http://localhost:1337/"
       }
     }'
```
   
Once you have your endpoint registered every time you send order created message you will receive order created message and then order updated message with the status configured in the config.yml file in the field StatusToExport.

Also if you send magento.inventory.source_stock_management.update you will receive as well a magento.inventory.aggregate_stock_management.updated messages to the aggregates configured in config.yml in the field AggregatesToExport.

  