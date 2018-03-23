FAKE MAGENTO ORDER MANAGEMENT API
=================================

How to run in local:

```go run main.go --config config.yml```

How to run in docker:

```
make build
make run
```

You can then register your integration with :

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
   
Once you have your endpoint registered every time you send magento.sales.order_management.create message you will receive magento.sales.order_management.created message and then magento.sales.order_management.updated messages with the status configured in the config.yml file in the field StatusToExport.

Also if you send magento.inventory.source_stock_management.update you will receive as well a magento.inventory.aggregate_stock_management.updated messages to the aggregates configured in config.yml in the field AggregatesToExport.

  