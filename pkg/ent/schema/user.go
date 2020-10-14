package schema

import (
	"regexp"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
	"github.com/facebook/ent/schema/mixin"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

const (
	userRoleAdmin   string = "admin"
	userRoleCompany string = "company"
	userRoleStudent string = "student"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.
			UUID("id", uuid.UUID{}).
			Immutable().
			Default(uuid.New),
		field.
			String("email").
			Unique().
			Match(emailRegex),
		field.
			String("password").
			Sensitive(),
		field.
			String("name").
			Default("unknown"),
		field.
			Enum("role").
			Values(userRoleAdmin, userRoleCompany, userRoleStudent).
			Default(userRoleStudent),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

// Mixin of the User
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
