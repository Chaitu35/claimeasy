package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)


type ClinicCreds struct {
	ent.Schema
}

func (ClinicCreds) Fields() []ent.Field {
	return []ent.Field{
		field.String("system").NotEmpty(),
		field.String("username").NotEmpty().Sensitive(),
		field.String("password").NotEmpty().Sensitive(),
	}
}


func (ClinicCreds) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("clinics", Clinics.Type).Ref("clinic_creds").Unique().Required(),
	
}
}