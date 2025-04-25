# Changelog

## [v1.0.0] - 2025-04-24
### Initial Release
- `.env` file parser with typed value support in Go
- Supported types include:
  - `int`
  - `float64`
  - `bool`
  - `string`
  - `[]string` (comma-separated list)
  - `[]int` (comma-separated list)
  - `map[string]interface{}` (from JSON string)
- Default fallback support for all types **(except placeholders — coming soon)**
- Placeholder replacement using `${VAR_NAME}` syntax
- Recursive placeholder resolution with circular reference detection
- Function to print all env variables with masking for sensitive values like `PASSWORD`, `SECRET`, etc.
- Configurable key masking for logs and debug output