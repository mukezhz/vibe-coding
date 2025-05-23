# Go Coding Guidelines for GitHub Copilot

## General Go Principles

### Code Organization
- Follow the wesionaryTEAM/go_clean_architecture Go project layout
- Use packages to organize code by functionality, not by type
- Keep package names simple, short, and meaningful
- One package per directory

### Naming Conventions
- Use `camelCase` for private variables, functions, and methods
- Use `PascalCase` for exported (public) variables, functions, methods, and types
- Use short but descriptive names for variables
- Prefer clarity over brevity for function and method names
- Use acronyms consistently (e.g., `HTTP`, `URL`, `ID`) - keep them all uppercase or all lowercase based on where they appear in the name

### Error Handling
- Check errors immediately after function calls
- Don't use `panic` or `recover` in production code; use proper error handling
- Return errors rather than using panic
- Use custom error types for specific error conditions
- Use the `errors.Is()` and `errors.As()` functions for error checking
- Wrap errors with context using `fmt.Errorf("doing X: %w", err)`

```go
if err != nil {
    return fmt.Errorf("failed to fetch user data: %w", err)
}