package main

import (
	"context"
	"cryptobot-catch/cmd/utils"
	"cryptobot-catch/internal/config"
	"cryptobot-catch/internal/core"
	"cryptobot-catch/internal/core/cheques"
	"cryptobot-catch/pkg/cryptobot"
	"encoding/json"
	"fmt"
	"github.com/gotd/contrib/bg"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
	"os"
)

func main() {
	ctx := context.Background()
	var catchConfig config.CatchConfig
	raw, err := os.ReadFile("configuration.json")

	zapLogger, _ := zap.NewProduction()
	zap.ReplaceGlobals(zapLogger)
	defer utils.Must(zapLogger.Sync())

	utils.Must(err)
	utils.Must(json.Unmarshal(raw, &catchConfig))

	catcherClient := telegram.NewClient(catchConfig.Catcher.AppID, catchConfig.Catcher.AppHash,
		telegram.Options{
			SessionStorage: &session.FileStorage{Path: "sessions/catcher.json"},
			Logger:         zapLogger,
		},
	)

	activatorClient := telegram.NewClient(catchConfig.Activator.AppID, catchConfig.Activator.AppHash,
		telegram.Options{
			SessionStorage: &session.FileStorage{Path: "sessions/activator.json"},
			Logger:         zapLogger,
		},
	)
	activatorStop, err := bg.Connect(activatorClient)
	utils.Must(err)
	if status, err := activatorClient.Auth().Status(ctx); (err == nil && !status.Authorized) || err != nil {
		utils.Must(fmt.Errorf("failed to authorize activator"))
	}
	catcherStop, err := bg.Connect(catcherClient)
	utils.Must(err)
	if status, err := catcherClient.Auth().Status(ctx); (err == nil && !status.Authorized) || err != nil {
		utils.Must(fmt.Errorf("failed to authorize catcher"))
	}
	defer func() {
		utils.Must(activatorStop())
		utils.Must(catcherStop())
	}()

	resolvedCryptoBot, err := activatorClient.API().ContactsResolveUsername(ctx, &tg.ContactsResolveUsernameRequest{
		Username: "send",
	})
	utils.Must(err)
	cryptoBot := cryptobot.NewCryptoBot(message.NewSender(activatorClient.API()),
		&tg.InputPeerUser{
			UserID:     resolvedCryptoBot.Users[0].GetID(),
			AccessHash: resolvedCryptoBot.Users[0].(*tg.User).AccessHash,
		},
	)
	chequeCatcher := core.NewCatcher(cryptoBot, &cheques.InlineDetectStrategy{}, &cheques.RegexFullChequeIDDetectStrategy{})
	utils.Must(chequeCatcher.Run(ctx, catcherClient.API()))

}
