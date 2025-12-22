package data

// DataKind represents the kind of data
type DataKind int

const (
	KindPrimitive DataKind = iota
	KindArray
	KindObject
)

// Accessor provides unified interface for accessing different data types
type Accessor interface {
	// Kind returns the kind of data (struct, map, slice, primitive)
	Kind() DataKind

	// IsNil checks if the underlying data is nil
	IsNil() bool

	// AsField returns field accessor (for primitives)
	AsField() (FieldAccessor, error)

	// AsObject returns object accessor (for structs/maps)
	AsObject() (ObjectAccessor, error)

	// AsArray returns array accessor (for slices/arrays)
	AsArray() (ArrayAccessor, error)
}

// FieldAccessor accesses primitive values
type FieldAccessor interface {
	Accessor

	// Value returns the underlying value
	Value() interface{}

	// String returns string representation
	String() string

	// Int returns int64 value
	Int() (int64, error)

	// Float returns float64 value
	Float() (float64, error)

	// Bool returns bool value
	Bool() (bool, error)
}

// ObjectAccessor accesses struct/map fields
type ObjectAccessor interface {
	Accessor

	// GetField returns field by name (supports dot notation for embedded)
	GetField(name string) (Accessor, bool)

	// Fields returns all field names
	Fields() []string

	// Len returns number of fields
	Len() int
}

// ArrayAccessor accesses slice/array elements
type ArrayAccessor interface {
	Accessor

	// GetIndex returns element at index
	GetIndex(idx int) (Accessor, bool)

	// Len returns array length
	Len() int

	// Iterate calls fn for each element
	Iterate(fn func(idx int, elem Accessor) error) error
}
