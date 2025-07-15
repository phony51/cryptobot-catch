package main

import (
	"context"
	"cryptobot-catch/internal/config"
	"cryptobot-catch/internal/utils"
	"cryptobot-catch/pkg/authenticators"
	"encoding/json"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"os"
)

func main() {

	var catchConfig config.CatchConfig
	raw, err := os.ReadFile("configuration.json")
	utils.Must(err)
	utils.Must(json.Unmarshal(raw, &catchConfig))

	catcherClient := telegram.NewClient(catchConfig.Catcher.AppID, catchConfig.Catcher.AppHash,
		telegram.Options{
			SessionStorage: &session.FileStorage{Path: "sessions/catcher.json"},
		},
	)

	activatorClient := telegram.NewClient(catchConfig.Activator.AppID, catchConfig.Activator.AppHash,
		telegram.Options{
			SessionStorage: &session.FileStorage{Path: "sessions/activator.json"},
		},
	)
	ctx := context.Background()

	utils.Must(catcherClient.Run(ctx, func(ctx context.Context) error {
		return catcherClient.Auth().IfNecessary(ctx, auth.NewFlow(
			auth.Constant(
				catchConfig.Catcher.Phone,
				catchConfig.Catcher.Password,
				&authenticators.PromptCodeAuthenticator{
					Prompt: "enter the confirmation code for catcher: ",
				},
			),
			auth.SendCodeOptions{},
		))
	}))

	utils.Must(activatorClient.Run(ctx, func(ctx context.Context) error {
		return activatorClient.Auth().IfNecessary(ctx, auth.NewFlow(
			auth.Constant(
				catchConfig.Activator.Phone,
				catchConfig.Activator.Password,
				&authenticators.PromptCodeAuthenticator{
					Prompt: "enter the confirmation code for activator: ",
				},
			),
			auth.SendCodeOptions{},
		))
	}))

}
