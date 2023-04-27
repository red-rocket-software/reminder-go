package firestore

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context, options option.ClientOption) (client *auth.Client, err error) {
	fireApp, err := firebase.NewApp(ctx, nil, options)
	if err != nil {
		return nil, err
	}

	fireClient, err := fireApp.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return fireClient, nil
}
