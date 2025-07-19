package config

import (
	"cryptobot-catch/internal/core/cheques"
	"cryptobot-catch/internal/utils"
	"encoding/json"
	"fmt"
)

type ExtractorsNotSelectedError struct{}

func (e ExtractorsNotSelectedError) Error() string {
	return "extractors not selected"
}

type Extractors []Extractor

func (es *Extractors) Build() []cheques.Extractor {
	cs := make([]cheques.Extractor, len(*es))
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
	var es_ []Extractor

	if err := json.Unmarshal(data, &es_); err != nil {
		return err
	}

	cs := make([]Extractor, len(es_))
	for i, e := range es_ {
		cs[i] = e
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

type Extractor string

func (n *Extractor) checkName() error {
	if _, ok := extractors[string(*n)]; ok {
		return nil
	}
	return ExtractorNotExistsError{string(*n)}
}

func (n *Extractor) UnmarshalJSON(data []byte) error {
	var n_ string
	if err := json.Unmarshal(data, &n_); err != nil {
		return err
	}
	*n = Extractor(n_)
	return n.checkName()
}

func (n *Extractor) Build() cheques.Extractor {
	return extractors[string(*n)]
}

type InvalidExtractorError struct {
	Name string
}

func (e InvalidExtractorError) Error() string {
	return fmt.Sprintf("invalid extractor name %s", e.Name)
}

var available = []cheques.Extractor{&cheques.InlineExtractor{}, &cheques.TextExtractor{}}

var extractors = func() map[string]cheques.Extractor {
	m := make(map[string]cheques.Extractor)
	for _, extractor := range available {
		m[extractor.Name()] = extractor
	}
	return m
}()
