package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OAuth struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AccessToken  string             `bson:"access_token,omitempty" json:"access_token,omitempty"`
	RefreshToken string             `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	ExpiresIn    int64              `bson:"expires_in,omitempty" json:"expires_in,omitempty"`
	TokenType    string             `bson:"token_type,omitempty" json:"token_type,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
