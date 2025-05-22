package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)


type RefreshToken struct {
	ent.Schema
}

func (RefreshToken) Fields() []ent.Field {
	return []ent.Field{
	field.String("token").Unique(),
	field.Int("user_id"),
	field.Time("issued_at").Default(time.Now),
	field.Time("expires_at").Default(time.Now).Nillable(),
	field.Bool("is_revoked").Default(false),
	field.String("user_agent").Nillable().Optional(),
	field.String("ip_address").Nillable().Optional(),
	}
}
