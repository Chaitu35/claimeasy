package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	//"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
)

type PasswordResetToken struct{
	ent.Schema
}

func (PasswordResetToken) Fields() []ent.Field {
return []ent.Field{
	field.String("token").Unique().NotEmpty().Sensitive().MaxLen(100),
	field.Time("expires_at"),
	field.Time("created_at").Default(time.Now),
}
}

func (PasswordResetToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("password_reset_tokens"),
	}
}