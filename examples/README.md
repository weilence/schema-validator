# Examples

This directory contains examples demonstrating various features of the schema-validator library.

## Running Examples

Each example is in its own subdirectory. To run an example:

```bash
# Run basic validation examples
go run examples/basic/main.go

# Run schema ToString examples
go run examples/tostring/main.go
```

## Available Examples

### 1. Basic Validation (`basic/`)

Comprehensive examples covering the core validation features:

- **Example 1**: Tag-based validation with structs
- **Example 2**: Code-based validation with maps
- **Example 3**: Embedded struct with private fields
- **Example 4**: Array validation
- **Example 5**: Cross-field validation
- **Example 6**: Error handling patterns

**Run it:**
```bash
go run examples/basic/main.go
```

### 2. Schema ToString (`tostring/`)

Examples demonstrating the schema JSON serialization feature:

- Simple field schemas
- Array schemas
- Object schemas
- Nested object schemas
- Complex schemas with arrays and objects

**Run it:**
```bash
go run examples/tostring/main.go
```

## Example Output

### Basic Validation Example

```
=== Example 1: Tag-based Validation ===
Valid user result: true
Invalid user result: false
  - email: email
  - password: min_length
  - confirmPassword: eqfield
  - age: min
```

### Schema ToString Example

```
=== Schema ToString Examples ===

1. Simple Field Schema:
{
  "type": "field",
  "optional": false,
  "validators": [
    {
      "name": "required"
    },
    {
      "name": "min_length",
      "value": 5
    },
    {
      "name": "email"
    }
  ]
}
```

## Learn More

For detailed documentation, see the main [README](../README.md) or [CLAUDE.md](../CLAUDE.md).
