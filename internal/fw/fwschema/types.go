package fwschema

// This is copied from terraform-plugin-framework/internal/fwschema/schema.go@6246b0c15d868c2de240bd0ce7a49c861056fe3d

type BlockNestingMode uint8

const (
	// BlockNestingModeUnknown is an invalid nesting mode, used to catch when a
	// nesting mode is expected and not set.
	BlockNestingModeUnknown BlockNestingMode = 0

	// BlockNestingModeList is for attributes that represent a list of objects,
	// with multiple instances of those attributes nested inside a list
	// under another attribute.
	BlockNestingModeList BlockNestingMode = 1

	// BlockNestingModeSet is for attributes that represent a set of objects,
	// with multiple, unique instances of those attributes nested inside a
	// set under another attribute.
	BlockNestingModeSet BlockNestingMode = 2

	// BlockNestingModeSingle is for attributes that represent a single object.
	// The object cannot be repeated in the practitioner configuration.
	//
	// While the framework implements support for this block nesting mode, it
	// is not thoroughly tested in production Terraform environments beyond the
	// resource timeouts block from the older Terraform Plugin SDK. Use single
	// nested attributes for new implementations instead.
	BlockNestingModeSingle BlockNestingMode = 3
)

type NestingMode uint8

const (
	// NestingModeUnknown is an invalid nesting mode, used to catch when a
	// nesting mode is expected and not set.
	NestingModeUnknown NestingMode = 0

	// NestingModeSingle is for attributes that represent a struct or
	// object, a single instance of those attributes directly nested under
	// another attribute.
	NestingModeSingle NestingMode = 1

	// NestingModeList is for attributes that represent a list of objects,
	// with multiple instances of those attributes nested inside a list
	// under another attribute.
	NestingModeList NestingMode = 2

	// NestingModeSet is for attributes that represent a set of objects,
	// with multiple, unique instances of those attributes nested inside a
	// set under another attribute.
	NestingModeSet NestingMode = 3

	// NestingModeMap is for attributes that represent a map of objects,
	// with multiple instances of those attributes, each associated with a
	// unique string key, nested inside a map under another attribute.
	NestingModeMap NestingMode = 4
)
