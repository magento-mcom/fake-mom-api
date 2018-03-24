package app

import "github.com/magento-mcom/fake-mom-api/api/handler"

type Config struct {
	StatusToExport       []handler.OrderStatus `yaml:"status_to_export"`
	AggregatesToExport   []string              `yaml:"aggregates_to_export"`
	DelayBetweenMessages int                   `yaml:"delay_between_messages"`
}
