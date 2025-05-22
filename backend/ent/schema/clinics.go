package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	regulatory "github.com/Chaitu35/claimeasy/backend/pkg/core"
)



type Clinics struct {
	ent.Schema
}
func (Clinics) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("city").NotEmpty(),
		field.String("regulatory").
		GoType(regulatory.Regulatory("")).
		Validate( func(s string) error {
			if !regulatory.IsValidRegulatory(s) {
				return fmt.Errorf("invalid regulatory: %s", s)
			}
			return nil
		}),
		field.String("email").Nillable().Optional(),
		field.String("phone").Nillable().Optional(),
		field.String("facilityID").Unique(),		
	}
	}

	func (Clinics) Edges() []ent.Edge {
		return []ent.Edge{
			edge.To("clinic_creds", ClinicCreds.Type),
			edge.To("doctors", Doctors.Type),	
		}
	}