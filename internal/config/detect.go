package config

import (
	"cryptobot-catch/internal/core/cheques/detecting"
	"cryptobot-catch/internal/utils"
	"encoding/json"
	"fmt"
)

type DetectStrategyNotSelectedError struct{}

func (e DetectStrategyNotSelectedError) Error() string {
	return "detect strategy not selected"
}

type InvalidDetectStrategyError struct {
	Name string
}

func (e InvalidDetectStrategyError) Error() string {
	return fmt.Sprintf("invalid detect strategy name %s", e.Name)
}

var detectStrategyAliasMap = map[string]detecting.DetectStrategy{
	(&detecting.RegexChequeIDDetectStrategy{}).Alias(): &detecting.RegexChequeIDDetectStrategy{},
	(&detecting.InlineDetectStrategy{}).Alias():        &detecting.InlineDetectStrategy{},
}

type DetectStrategies struct {
	Strategies []detecting.DetectStrategy
}

func (ss *DetectStrategies) checkLength() error {
	if len(ss.Strategies) > 0 {
		return nil
	}
	return DetectStrategyNotSelectedError{}
}

func (ss *DetectStrategies) UnmarshalJSON(data []byte) error {
	var ss_ []DetectStrategy
	if err := json.Unmarshal(data, &ss_); err != nil {
		return err
	}

	utils.RemoveDuplicate(ss_)
	if len(ss_) == 0 {
		return DetectStrategyNotSelectedError{}
	}

	ss.Strategies = make([]detecting.DetectStrategy, len(ss_))
	for i, s := range ss_ {
		ss.Strategies[i] = s.Strategy()
	}
	return nil
}

type DetectStrategyNotExistsError struct {
	Name string
}

func (e DetectStrategyNotExistsError) Error() string {
	return fmt.Sprintf("detect strategy not exists: %s", e.Name)
}

type DetectStrategy string

func (s *DetectStrategy) checkName() error {
	if _, ok := detectStrategyAliasMap[string(*s)]; ok {
		return nil
	}
	return DetectStrategyNotExistsError{string(*s)}
}

func (s *DetectStrategy) UnmarshalJSON(data []byte) error {
	var s_ string
	if err := json.Unmarshal(data, &s_); err != nil {
		return err
	}
	*s = DetectStrategy(s_)
	return s.checkName()
}

func (s *DetectStrategy) Strategy() detecting.DetectStrategy {
	return detectStrategyAliasMap[string(*s)]
}
