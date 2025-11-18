package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRepo struct {
	ID              bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name            string          `bson:"name" json:"name"`
	Username        string          `bson:"username" json:"username"`
	Password        string          `bson:"password" json:"-"`
	Email           string          `bson:"email" json:"email"`
	Roles           []bson.ObjectID `bson:"roles,omitempty" json:"roles,omitempty"`
	EmailVerifiedAt *time.Time      `bson:"email_verified_at,omitempty" json:"email_verified_at,omitempty"`
	RememberToken   *string         `bson:"remember_token,omitempty" json:"remember_token,omitempty"`
	Active          bool            `bson:"active" json:"active"`
	Audit           `bson:",inline"`
}

func (UserRepo) CollName() string { return "users" }

type RoleRepo struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Slug        string        `bson:"slug" json:"slug"`
	Name        string        `bson:"name" json:"name"`
	Permissions []string      `bson:"permissions" json:"permissions"`
	Audit       `bson:",inline"`
}

func (RoleRepo) CollName() string { return "roles" }
