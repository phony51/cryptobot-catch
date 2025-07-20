package config

import "time"

type Connection struct {
	DC             int           `json:"dc"`
	MaxConnections int64         `json:"maxConnections"`
	AckIntervalMs  time.Duration `json:"ackIntervalMs"`
	AckBatchSize   int           `json:"ackBatchSize"`
}
