Fake Magento Order Management API
=================================

# Goal

The goal of this project is to enable Magento Commerce developers to be able to work with an isolated test environement. There is no logic and all orders will be shipped
upon arrival in the fake API.

# Build

How to build and run in local:

```go run main.go --config config.yml```

How to run in docker:

```
make build
make run
```

# Connecting your Magento Commerce

First register your Magento Commerce:

```
$ bin/magento oms:service-bus:register
Registering MDC to Service Bus...
Done.
```


# Behaviour of fake API

Once you have your endpoint registered every time you send `magento.sales.order_management.create` message you will receive:
* `magento.sales.order_management.created`
* `magento.sales.order_management.updated` messages with the status configured in the `config.yml` file in the field StatusToExport
* `magento.logistics.fulfillment_management.customer_shipment_done`

You can emulate the call from Magento Commerce to the Fake Magento Order Management API as such:

```
$ curl -X POST http://localhost:24213/api/ -d @test-data/order-test.json \
  --header "Content-Type: application/json"

{"jsonrpc":"2.0","result":null,"error":"","id":null}%   
``` 

Also if you send `magento.inventory.source_stock_management.update` you will receive as well a `magento.inventory.aggregate_stock_management.updated` messages to the aggregates configured in config.yml in the field AggregatesToExport.
 