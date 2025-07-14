package main

import (
	"context"
	"cryptobot-catch/internal/config"
	"cryptobot-catch/internal/core"
	"cryptobot-catch/internal/utils"
	"cryptobot-catch/pkg/cryptobot"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var (
	encCfg = zapcore.EncoderConfig{
		MessageKey: "msg",
		LevelKey:   "level",
		TimeKey:    "timestamp",
		EncodeTime: zapcore.ISO8601TimeEncoder,
	}
	devCfg = zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    encCfg,
	}
	prodCfg = zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout", "logs.json"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    encCfg,
	}
)

func main() {
	logger, err := prodCfg.Build()
	utils.Must(err)
	zap.ReplaceGlobals(logger)

	var catchConfig config.CatchConfig
	raw, err := os.ReadFile("configuration.json")
	utils.Must(err)
	utils.Must(json.Unmarshal(raw, &catchConfig))
	logger.Info("configuration loaded", zap.String("configuration", fmt.Sprint(catchConfig)))

	ctx, shutdown := context.WithCancel(context.Background())
	defer shutdown()

	activatorClient := telegram.NewClient(catchConfig.Activator.AppID, catchConfig.Activator.AppHash,
		telegram.Options{
			SessionStorage: &session.FileStorage{Path: "sessions/activator.json"},
			Logger:         logger,
		},
	)

	catcherClient := telegram.NewClient(catchConfig.Catcher.AppID, catchConfig.Catcher.AppHash,
		telegram.Options{
			SessionStorage: &session.FileStorage{Path: "sessions/catcher.json"},
			Logger:         logger,
		},
	)

	err = activatorClient.Run(ctx, func(ctx context.Context) error {
		if _, err := activatorClient.Auth().Status(ctx); err != nil {
			return errors.Join(fmt.Errorf("failed to authorize activator"), err)
		}

		resolvedCryptoBot, err := activatorClient.API().ContactsResolveUsername(ctx, &tg.ContactsResolveUsernameRequest{
			Username: "send",
		})
		if err != nil {
			return err
		}
		cryptoBot := cryptobot.NewCryptoBot(message.NewSender(activatorClient.API()),
			&tg.InputPeerUser{
				UserID:     resolvedCryptoBot.Users[0].GetID(),
				AccessHash: resolvedCryptoBot.Users[0].(*tg.User).AccessHash,
			},
		)
		_ = catcherClient.Run(ctx, func(ctx context.Context) error {
			if _, err := catcherClient.Auth().Status(ctx); err != nil {
				return errors.Join(fmt.Errorf("failed to authorize catcher"), err)
			}
			chequeCatcher := core.NewCatcher(cryptoBot, &core.CatchOptions{
				PollingInterval:  time.Duration(catchConfig.PollingIntervalMs) * time.Millisecond,
				DetectStrategies: catchConfig.DetectBy.Strategies,
			})
			return chequeCatcher.Run(ctx, catcherClient.API())
		})
		return err
	})
}
