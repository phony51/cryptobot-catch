package config

import (
	"cryptobot-catch/internal/core/cheques"
	"fmt"
	"maps"
	"slices"
)

type DetectStrategyNotSelectedError struct{}

func (e DetectStrategyNotSelectedError) Error() string {
	return "detect strategy not selected"
}

type InvalidDetectStrategyError struct {
	StrategyName string
}

func (e InvalidDetectStrategyError) Error() string {
	return fmt.Sprintf("%s: %s", "invalid detect strategy name", e.StrategyName)
}

var detectStrategyAliasMap = map[string]cheques.DetectStrategy{
	"regexFullChequeID": &cheques.RegexFullChequeIDDetectStrategy{},
	"inline":            &cheques.InlineDetectStrategy{},
}

var detectStrategyAliasKeys = slices.Collect(maps.Keys(detectStrategyAliasMap))

type DetectStrategyAliases []DetectStrategyAlias

func (a *DetectStrategyAliases) Strategies() ([]cheques.DetectStrategy, error) {
	if len(*a) == 0 {
		return nil, DetectStrategyNotSelectedError{}
	}
	strategies := make([]cheques.DetectStrategy, len(*a))
	for i, alias := range *a {
		s, err := alias.Strategy()
		if err != nil {
			return nil, err
		}
		strategies[i] = s
	}
	return strategies, nil
}

type DetectStrategyAlias string

func (a *DetectStrategyAlias) Strategy() (cheques.DetectStrategy, error) {
	k := string(*a)
	if slices.Contains(detectStrategyAliasKeys, k) {
		return detectStrategyAliasMap[k], true
	}
	return nil, false
}
