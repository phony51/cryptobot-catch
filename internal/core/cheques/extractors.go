package cheques

import (
	"github.com/flier/gohs/hyperscan"
	"github.com/gotd/td/tg"
	"strings"
	"sync"
)

type ExtractFunc func(message *tg.Message) (string, bool)

type Extractor interface {
	Name() string
	Extract(message *tg.Message) (string, bool)
}

type InlineExtractor struct{}

func (ie *InlineExtractor) Name() string {
	return "inline"
}

const inlineChequeURLPrefix = "http://t.me/send?start=CQ"

func (ie *InlineExtractor) Extract(message *tg.Message) (string, bool) {
	if keyboard, ok := message.ReplyMarkup.(*tg.ReplyInlineMarkup); ok {
		if button, ok := keyboard.Rows[0].Buttons[0].(*tg.KeyboardButtonURL); ok {
			if strings.HasPrefix(button.URL, inlineChequeURLPrefix) {
				return button.URL[len(inlineChequeURLPrefix):], true
			}
		}
	}
	return "", false
}

type TextExtractor struct {
	db      hyperscan.BlockDatabase
	scratch *hyperscan.Scratch
	mu      sync.Mutex
}

const chequePrefix = "CQ"

func NewTextExtractor() (*TextExtractor, error) {
	pattern := hyperscan.NewPattern(chequePrefix+"[A-Za-z0-9]{10}", hyperscan.SomLeftMost)
	db, err := hyperscan.NewBlockDatabase(pattern)
	if err != nil {
		return nil, err
	}

	scratch, err := hyperscan.NewScratch(db)
	if err != nil {
		err = db.Close()
		return nil, err
	}

	return &TextExtractor{
		db:      db,
		scratch: scratch,
	}, nil
}

func (te *TextExtractor) Name() string {
	return "text"
}

func (te *TextExtractor) Close() error {
	te.mu.Lock()
	defer te.mu.Unlock()

	var err error
	if te.scratch != nil {
		err = te.scratch.Free()
		te.scratch = nil
	}
	if te.db != nil {
		err2 := te.db.Close()
		if err == nil {
			err = err2
		}
		te.db = nil
	}
	return err
}

func (te *TextExtractor) Extract(message *tg.Message) (string, bool) {
	if message.Message == "" {
		return "", false
	}

	te.mu.Lock()
	defer te.mu.Unlock()

	var matchFound string

	handler := func(id uint, from, to uint64, flags uint, context interface{}) error {
		text := context.(string)
		matchFound = text[from:to]
		return hyperscan.ErrScanTerminated
	}

	err := te.db.Scan([]byte(message.Message), te.scratch, handler, message.Message)
	if err != nil {
		return "", false
	}

	if matchFound != "" {
		return matchFound[len("CQ"):], true
	}

	return "", false
}
