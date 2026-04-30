package auth

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// TokenInfo contiene los datos que extraemos del token de Firebase
type TokenInfo struct {
	UID   string
	Email string
}

// FirebaseProvider verifica tokens de Firebase
type FirebaseProvider struct {
	client *firebaseAuth.Client
}

// NewFirebaseProvider inicializa el Admin SDK de Firebase
func NewFirebaseProvider(credentialsPath string) (*FirebaseProvider, error) {
	ctx := context.Background()

	opt := option.WithCredentialsFile(credentialsPath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting firebase auth client: %w", err)
	}

	fmt.Println("Firebase Auth initialized")

	return &FirebaseProvider{client: client}, nil
}

// VerifyToken verifica un ID token y retorna UID + email
func (f *FirebaseProvider) VerifyToken(ctx context.Context, idToken string) (*TokenInfo, error) {
	decoded, err := f.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	email, _ := decoded.Claims["email"].(string)

	return &TokenInfo{
		UID:   decoded.UID,
		Email: email,
	}, nil
}
