package firebase

// INICIALIZACION DEL SDK
import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type FirebaseClient struct {
	Auth *auth.Client
}

func NewFirebaseClient(credentialsPath string) (*FirebaseClient, error) {
	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return &FirebaseClient{Auth: authClient}, nil
}
