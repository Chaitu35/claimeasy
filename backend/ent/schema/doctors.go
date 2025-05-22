package schema
import("entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field")


	type Doctors struct {
		ent.Schema
	}

	func (Doctors) Fields() []ent.Field {
		return []ent.Field{
			field.String("name").NotEmpty(),
			field.String("specialization").NotEmpty(),
			field.String("license").Unique().NotEmpty(),
		}
	}

	func (Doctors) Edges() []ent.Edge {
		return []ent.Edge{
			edge.From("clinics", Clinics.Type).Ref("doctors").Unique().Required(),
		}
	}