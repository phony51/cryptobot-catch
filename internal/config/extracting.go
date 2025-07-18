package config

import (
	"cryptobot-catch/internal/core/cheques/extracting"
	"cryptobot-catch/internal/utils"
	"encoding/json"
	"fmt"
)

type ExtractorsNotSelectedError struct{}

func (e ExtractorsNotSelectedError) Error() string {
	return "extractors not selected"
}

type Extractors []ExtractorWorkerPoolConfig

func (es *Extractors) Build() []extracting.ExtractorWorkerPool {
	cs := make([]extracting.ExtractorWorkerPool, len(*es))
	for i, e := range *es {
		cs[i] = e.Build()
	}
	return cs
}

func (es *Extractors) checkLength() error {
	if len(*es) != 0 {
		return nil
	}
	return ExtractorsNotSelectedError{}
}

func (es *Extractors) UnmarshalJSON(data []byte) error {
	var es_ []ExtractorWorkerPoolConfig

	if err := json.Unmarshal(data, &es_); err != nil {
		return err
	}

	cs := make([]ExtractorWorkerPoolConfig, len(es_))
	for i, ext := range es_ {
		cs[i] = ext
	}
	*es = cs

	if err := es.checkLength(); err != nil {
		return err
	}
	utils.RemoveDuplicate(es_)
	return nil
}

type ExtractorNotExistsError struct {
	Name string
}

func (e ExtractorNotExistsError) Error() string {
	return fmt.Sprintf("extractor not exists: %s", e.Name)
}

type ExtractorWorkerPoolConfig struct {
	Name         ExtractorName `json:"name"`
	WorkersCount int           `json:"workersCount"`
}

func (c *ExtractorWorkerPoolConfig) checkWorkersCount() error {
	if c.WorkersCount < 1 {
		return WorkersCountMustBePositiveError{}
	}
	return nil
}

func (c *ExtractorWorkerPoolConfig) UnmarshalJSON(data []byte) error {
	type wrapper ExtractorWorkerPoolConfig
	var c_ wrapper

	if err := json.Unmarshal(data, &c_); err != nil {
		return err
	}
	*c = ExtractorWorkerPoolConfig(c_)
	return c.checkWorkersCount()
}

func (c *ExtractorWorkerPoolConfig) Build() extracting.ExtractorWorkerPool {
	return extracting.NewExtractorWorkerPool(c.Name.Build(), c.WorkersCount)
}

type WorkersCountMustBePositiveError struct{}

func (e WorkersCountMustBePositiveError) Error() string {
	return "workers count must be positive"
}

type ExtractorName string

func (n *ExtractorName) checkName() error {
	if _, ok := extractors[string(*n)]; ok {
		return nil
	}
	return ExtractorNotExistsError{string(*n)}
}

func (n *ExtractorName) UnmarshalJSON(data []byte) error {
	var n_ string
	if err := json.Unmarshal(data, &n_); err != nil {
		return err
	}
	*n = ExtractorName(n_)
	return n.checkName()
}

func (n *ExtractorName) Build() extracting.Extractor {
	return extractors[string(*n)]
}

type InvalidExtractorError struct {
	Name string
}

func (e InvalidExtractorError) Error() string {
	return fmt.Sprintf("invalid extractor name %s", e.Name)
}

var available = []extracting.Extractor{extracting.InlineExtractor, extracting.TextExtractor}

var extractors = func() map[string]extracting.Extractor {
	m := make(map[string]extracting.Extractor)
	for _, extractor := range available {
		m[extractor.Name] = extractor
	}
	return m
}()
