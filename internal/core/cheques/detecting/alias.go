package detecting

type MappedDetectStrategy interface {
	DetectStrategy
	Mapper
}

type Mapper interface {
	Alias() string
}

func (s *InlineDetectStrategy) Alias() string {
	return "inline-detect"
}

func (s *RegexChequeIDDetectStrategy) Alias() string {
	return "regex-cheque-id"
}
