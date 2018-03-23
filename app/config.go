package app

import "github.com/magento-mcom/fake-api/api/handler"

type Config struct {
	StatusToExport     []handler.OrderStatus `yaml:"status_to_export"`
	AggregatesToExport []string              `yaml:"aggregates_to_export"`
}
