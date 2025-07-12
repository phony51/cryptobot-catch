package config

import (
	"cryptobot-catch/internal/core/cheques"
	"iter"
	"maps"
	"reflect"
	"slices"
)

type CatchConfig struct {
	Catcher           Credentials `json:"catcher"`
	Activator         Credentials `json:"activator"`
	PollingIntervalMs int         `json:"pollingIntervalMs"`
	DetectBy          DetectBy    `json:"detectBy"`
}

type Credentials struct {
	AppID    int    `json:"appID"`
	AppHash  string `json:"appHash"`
	Phone    string `json:"phone"`
	Password string `json:"password,omitempty"`
}
