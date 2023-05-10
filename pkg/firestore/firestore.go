package firestore

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

//go:generate mockgen -source=firestore.go -destination=mocks/firestore.go

type Client interface {
	VerifyIDToken(idToken string) (*auth.Token, error)
	GetUser(userID string) (*auth.UserRecord, error)
}

type FirebaseClient struct {
	client  *auth.Client
	context context.Context
}

func NewClient(ctx context.Context, options option.ClientOption) (client Client, err error) {
	fireApp, err := firebase.NewApp(ctx, nil, options)
	if err != nil {
		return nil, err
	}

	fireClient, err := fireApp.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &FirebaseClient{context: ctx, client: fireClient}, nil
}

func (f *FirebaseClient) VerifyIDToken(idToken string) (*auth.Token, error) {
	token, err := f.client.VerifyIDToken(f.context, idToken)
	return token, err
}

func (f *FirebaseClient) GetUser(userID string) (*auth.UserRecord, error) {
	return f.client.GetUser(f.context, userID)
}
