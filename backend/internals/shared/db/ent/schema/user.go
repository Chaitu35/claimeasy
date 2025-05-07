package schema

import (
	"entgo.io/ent"
	 "entgo.io/ent/schema/field"

	 "github.com/chaitu35/claimeasy/shared/types"

)


type User struct {ent.Schema}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").Unique(),
		field.String("username").Unique(),
		field.String("password").Sensitive(),
		field.String("email").Unique(),
		field.String("role").Default("user"),
		field.String("first_name").Optional(),
		field.String("last_name").Optional(),	
		field.String("phone").Optional(),		
		field.String("country").Optional(),
		field.String("national_id").Optional(),
	}
}