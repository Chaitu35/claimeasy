package schema

import (
	"entgo.io/ent"
	//"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)
 type Permissions struct {
	ent.Schema
 }
func (Permissions) Fields() []ent.Field {
    return []ent.Field{
        field.String("module"),
        field.Bool("can_view").Default(false),
        field.Bool("can_create").Default(false),
        field.Bool("can_edit").Default(false),  // Changed from "update"
        field.Bool("can_delete").Default(false),
        field.Bool("can_export").Default(false),
        field.Bool("can_print").Default(false),
    }
}

func (Permissions) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user",User.Type).Ref("permissions").Unique(),
	}

}