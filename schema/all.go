package schema

// All returns all schema instances.
func All() []interface{} {
	return []interface{}{
		new(TagDefinition),
	}
}
