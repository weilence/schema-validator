package reflect

import (
	"reflect"
	"sync"
	"unsafe"
)

// StructInfo caches struct field metadata
type StructInfo struct {
	Type   reflect.Type
	Fields map[string]*FieldInfo
	names  []string
}

// FieldInfo describes a struct field (including embedded)
type FieldInfo struct {
	Name      string
	Index     []int // Field index path for reflect.FieldByIndex
	Type      reflect.Type
	Embedded  bool
	Anonymous bool
}

// GetValue accesses field value, including private embedded fields
func (f *FieldInfo) GetValue(structVal reflect.Value) reflect.Value {
	// Use FieldByIndex which works with embedded fields
	fieldVal := structVal.FieldByIndex(f.Index)

	// For private fields, use unsafe to make accessible
	if !fieldVal.CanInterface() && fieldVal.CanAddr() {
		// This allows reading private fields
		fieldVal = reflect.NewAt(fieldVal.Type(), unsafe.Pointer(fieldVal.UnsafeAddr())).Elem()
	}

	return fieldVal
}

var (
	structInfoCache sync.Map
)

// GetStructInfo returns cached struct info or builds it
func GetStructInfo(typ reflect.Type) *StructInfo {
	if info, ok := structInfoCache.Load(typ); ok {
		return info.(*StructInfo)
	}

	info := buildStructInfo(typ)
	structInfoCache.Store(typ, info)
	return info
}

func buildStructInfo(typ reflect.Type) *StructInfo {
	info := &StructInfo{
		Type:   typ,
		Fields: make(map[string]*FieldInfo),
		names:  []string{},
	}

	// Walk through fields including embedded
	walkFields(typ, nil, info)

	return info
}

// GetField returns field info by name
func (s *StructInfo) GetField(name string) (*FieldInfo, bool) {
	field, ok := s.Fields[name]
	return field, ok
}

// FieldNames returns all field names
func (s *StructInfo) FieldNames() []string {
	return s.names
}
