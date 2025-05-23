package schema

import (
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"github.com/Chaitu35/claimeasy/backend/pkg/core/role"
)

type User struct{
	ent.Schema
}

func (User) Fields() []ent.Field{

	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("username").Unique().NotEmpty(),
		field.String("password_hash").NotEmpty().Sensitive(),
		field.String("email").Unique().NotEmpty(),
		field.String("mobile").Unique().NotEmpty(),
	field.String("role").GoType(role.Role("")).Default(string(role.RoleCoder)).
			Validate(func(s string) error {
				// Check if the role is valid
				if !role.IsValidRole(role.Role(s)) {
					return fmt.Errorf("invalid role: %s", s)
				}
				return nil
			}),	
		field.String("status").Default("active"),
		field.String("national_id").Optional().Nillable(),
		field.String("country").Optional().Nillable(),
		field.Bool("is_active").Default(true),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),

	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
	edge.To("permissions", Permissions.Type),
	edge.To("password_reset_tokens", PasswordResetToken.Type),
}
}