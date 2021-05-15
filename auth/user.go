package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/soumitra003/go-webrtc/internal/oauth2/google"
)

// User struct
type User struct {
	ID           string    `json:"id" bson:"_id"`
	FirstName    string    `json:"firstName" bson:"firstName"`
	LastName     string    `json:"lastName" bson:"lastName"`
	Email        string    `json:"email" bson:"email"`
	Mobile       string    `json:"mobile" bson:"mobile"`
	FirstLoginAt time.Time `json:"firstLoginAt" bson:"firstLoginAt"`
	LastLoginAt  time.Time `json:"lastLoginAt" bson:"lastLoginAt"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt"`
}

// UserClaims JWT claims for User struct
type UserClaims struct {
	ID            string `json:"id"`
	jwt.MapClaims `json:"-"`
}

// NewUser creates a new user with assigned uuid
func NewUser() User {
	return User{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func generateTokenWithClaims(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{ID: user.ID})
	tokenString, err := token.SignedString([]byte("secretKey"))
	return tokenString, err
}

func newUserFromGoogleIDToken(idToken string) (User, error) {
	var user User
	idTokenPayload, err := google.DecodeIDToken(idToken)
	if err != nil {
		return user, err
	}

	user = NewUser()
	user.Email = idTokenPayload.Email
	return user, nil
}
