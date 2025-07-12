package authenticators

import (
	"context"
	"fmt"
	"github.com/gotd/td/tg"
)

type PromptCodeAuthenticator struct {
	Prompt string
}

func (a *PromptCodeAuthenticator) Code(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
	var code string
	fmt.Println(a.Prompt)
	_, err := fmt.Scanln(&code)
	if err != nil {
		return "", err
	}
	return code, nil
}
