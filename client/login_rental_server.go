package client

import (
	"context"
	"fmt"
	"time"

	"github.com/OmineDev/flowers-for-machines/core/bunker/auth"
)

// LoginRentalServer ..
func LoginRentalServer(cfg Config) (client *Client, err error) {
	authClient, err := auth.CreateClient(&auth.ClientOptions{
		AuthServer: cfg.AuthServerAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("LoginRentalServer: %v", err)
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()

	authenticator := auth.NewAccessWrapper(
		authClient,
		cfg.RentalServerCode,
		cfg.RentalServerPasscode,
		cfg.AuthServerToken,
		"", "",
	)
	conn, err := openConnection(ctx, authenticator)
	if err != nil {
		return nil, fmt.Errorf("LoginRentalServer: %v", err)
	}

	client = &Client{connection: conn, authClient: authClient}
	err = NewChallengeSolver(client).CopeChallenge()
	if err != nil {
		return nil, fmt.Errorf("LoginRentalServer: %v", err)
	}

	return client, nil
}
