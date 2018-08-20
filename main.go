package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/magento-mcom/fake-mom-api/api"
	"github.com/magento-mcom/fake-mom-api/api/handler"
	"github.com/magento-mcom/fake-mom-api/app"
	"github.com/magento-mcom/fake-mom-api/order"
	"github.com/satori/go.uuid"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"github.com/magento-mcom/fake-mom-api/consumer"
)

func main() {
	content, err := ioutil.ReadFile("config.yml")

	if err != nil {
		panic(fmt.Sprintf("Error reading config file: %v", err))
	}

	config := app.Config{}

	if err = yaml.Unmarshal(content, &config); err != nil {
		panic(fmt.Sprintf("Error composing configuration: %v", err))
	}

	router := gin.New()

	r := api.NewRegistry()
	p := api.NewPublisher(r)
	or := order.NewOrderRegistry()
	q := consumer.NewConsumerQueue()
	c := consumer.NewConsumer(q, p, config.DelayBetweenMessages)
	mh := map[string]api.Handler{
		"magento.service_bus.remote.register":              handler.NewRegisterHandler(r),
		"magento.sales.order_management.create":            handler.NewCreateOrderHandler(q, config.StatusToExport, or),
		"magento.inventory.source_stock_management.update": handler.NewSourceUpdateHandler(p, config.AggregatesToExport),
	}


	d := api.NewDispatcher(mh)

	apiGatewayHandler := func(ctx *gin.Context) {
		data := api.Request{}

		b, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			return
		}

		if err := json.Unmarshal(b, &data); err != nil {
			fmt.Printf("%v", err)
			return
		}

		res, err := d.Dispatch(data)

		respBody := api.Response{
			ID:      data.ID,
			JsonRpc: "2.0",
		}

		if err == nil {
			m, err := json.Marshal(res)
			if err == nil {
				raw := json.RawMessage(m)
				respBody.Result = &raw
			}

		}

		if err != nil {
			respBody.Error = err.Error()
		}

		ctx.JSON(http.StatusOK, respBody)
	}

	router.POST("/api/", apiGatewayHandler)
	router.POST("/api/:mode/:endpoint", apiGatewayHandler)

	router.GET("/order/:id", func(ctx *gin.Context) {
		orderId := ctx.Param("id")

		id, _ := uuid.NewV4()
		respBody := api.Response{
			ID:      id.String(),
			JsonRpc: "2.0",
		}

		if !or.Exists(orderId) {
			respBody.Error = fmt.Sprintf("Order with id %v not exists.", orderId)
		}

		ctx.JSON(http.StatusOK, respBody)
	})

	go c.Run()

	srv := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		Addr:              ":24213",
		Handler:           router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
