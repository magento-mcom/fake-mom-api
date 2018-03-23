package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/magento-mcom/fake-api/api"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	router := gin.New()

	router.POST("/", func(ctx *gin.Context) {
		data := api.Request{}

		b, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			return
		}

		if err := json.Unmarshal(b, &data); err != nil {
			fmt.Printf("%v", err)
			return
		}

		fmt.Printf("%v\n", data)
		fmt.Printf("%v\n", string(*data.Params))
	})

	srv := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		Addr:              ":1337",
		Handler:           router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
