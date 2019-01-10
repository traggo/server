package graphql

import (
	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlschema"
)

// NewDirective creates a new directive.
func NewDirective() gqlschema.DirectiveRoot {
	return gqlschema.DirectiveRoot{
		HasRole: auth.HasRole(),
	}
}
