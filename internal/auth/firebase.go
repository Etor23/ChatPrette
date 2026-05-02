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

// FirebaseProvider verifica y gestiona autenticación con Firebase
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

// CreateUserWithEmail crea un nuevo usuario en Firebase con email y password
// Retorna el UID del usuario creado
func (f *FirebaseProvider) CreateUserWithEmail(ctx context.Context, email, password string) (uid string, idToken string, err error) {
	// Crear usuario en Firebase
	params := (&firebaseAuth.UserToCreate{}).
		Email(email).
		Password(password).
		EmailVerified(false)

	user, err := f.client.CreateUser(ctx, params)
	if err != nil {
		return "", "", fmt.Errorf("error creating user: %w", err)
	}

	// Retornar el UID
	return user.UID, "", nil
}

// SignInWithEmailPassword valida credenciales y retorna el UID del usuario
// NOTA: Esta función requiere que Firebase REST API esté habilitada
// Por ahora usamos GetUserByEmail como alternativa para verificar existencia
func (f *FirebaseProvider) GetUserByEmail(ctx context.Context, email string) (*firebaseAuth.UserRecord, error) {
	user, err := f.client.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UserExists verifica si un usuario existe en Firebase por email
func (f *FirebaseProvider) UserExists(ctx context.Context, email string) (bool, error) {
	_, err := f.client.GetUserByEmail(ctx, email)
	if err != nil {
		if firebaseAuth.IsUserNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// DeleteUser elimina un usuario de Firebase (útil para cleanup)
func (f *FirebaseProvider) DeleteUser(ctx context.Context, uid string) error {
	return f.client.DeleteUser(ctx, uid)
}
