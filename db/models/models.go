// models/user.go
package models

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            string    `bson:"_id,omitempty"`
	Username      string    `bson:"username" json:"username"`
	Email         string    `bson:"email" json:"email"`
	Avatar        string    `bson:"avatar" json:"avatar"`
	Password      string    `bson:"password" json:"password"`
	Qualification string    `bson:"qualification" json:"qualification"`
	MobileNumber  string    `bson:"mobileNumber" json:"mobileNumber"`
	Country       string    `bson:"country" json:"country"`
	Role          string    `bson:"role" json:"role"`
	Status        string    `bson:"status" json:"status"`
	RefreshToken  string    `bson:"refreshToken" json:"refreshToken"`
	Points        string    `bson:"points" json:"points"` // Assuming Points is a string for simplicity
	CreatedAt     time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time `bson:"updatedAt" json:"updatedAt"`
}

// User roles and statuses
const (
	RoleStudent = "student"
	RoleTeacher = "teacher"
	RoleAdmin   = "admin"
	RolePro     = "pro"

	StatusMember = "member"
	StatusPro    = "pro"
)

// HashPassword hashes the user's password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password is correct
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// GenerateTokens generates access and refresh tokens
func (u *User) GenerateTokens(secretAccess, secretRefresh string) (string, string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  u.ID,
		"exp": time.Now().Add(15 * 24 * time.Hour).Unix(), // 15 days
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(secretRefresh))
	if err != nil {
		return "", "", err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
		"status":   u.Status,
		"role":     u.Role,
		"avatar":   u.Avatar,
		"exp":      time.Now().Add(2 * 24 * time.Hour).Unix(), // 2 days
	})
	accessTokenString, err := accessToken.SignedString([]byte(secretAccess))
	if err != nil {
		return "", "", err
	}

	return refreshTokenString, accessTokenString, nil
}

// SaveUser  saves the user to the database
func SaveUser(collection *mongo.Collection, user *User) error {
	_, err := collection.InsertOne(context.TODO(), user)
	return err
}

// FindUser ByEmail finds a user by email
func FindUserByEmail(collection *mongo.Collection, email string) (*User, error) {
	var user User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
