# Contributing to go-common

Thank you for your interest in contributing to go-common! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful, constructive, and professional in all interactions.

## Getting Started

### Prerequisites

- Go 1.25 or higher
- Git
- Make
- (Optional) golangci-lint for linting

### Setting Up Development Environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/go-common.git
   cd go-common
   ```

3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/soderasen-au/go-common.git
   ```

4. Install dependencies:
   ```bash
   make deps
   ```

5. Verify everything works:
   ```bash
   make check
   ```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/my-awesome-feature
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Adding tests

### 2. Make Your Changes

#### Coding Standards

**Error Handling:**
- All functions should return `*util.Result` for consistent error handling
- Use `util.Error()` for wrapping errors
- Use `util.OK()` for success
- Chain errors with `.With()` method

Example:
```go
func YourFunction(ctx string) *util.Result {
    if res := InnerFunction(ctx); res.Code != 0 {
        return util.With(ctx).ResultError(res)
    }
    return util.OK(ctx)
}
```

**Testing:**
- Write tests for all new functionality
- Aim for >75% code coverage
- Use table-driven tests where appropriate
- Test both success and error paths

Example:
```go
func TestYourFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "result", false},
        {"invalid input", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := YourFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("YourFunction() error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("YourFunction() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Documentation:**
- Add godoc comments for all exported functions, types, and constants
- Include examples in documentation where helpful
- Update README.md if adding new packages or features

**Code Style:**
- Run `make fmt` before committing
- Follow standard Go conventions
- Keep functions small and focused
- Use descriptive variable names

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run tests for specific package
make test-util
make test-fx
make test-crypto

# Generate coverage report
make coverage

# Check code formatting
make fmt-check

# Run linters
make lint

# Run all quality checks
make check
```

### 4. Commit Your Changes

Write clear, descriptive commit messages:

```bash
git add .
git commit -m "feat: add new encryption algorithm support

- Add AES-256-GCM encryption
- Add corresponding tests
- Update documentation"
```

Commit message format:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Adding or updating tests
- `refactor:` - Code refactoring
- `chore:` - Maintenance tasks

### 5. Push and Create Pull Request

```bash
git push origin feature/my-awesome-feature
```

Then create a Pull Request on GitHub with:
- Clear title describing the change
- Description of what changed and why
- Reference any related issues
- Screenshots/examples if applicable

## Pull Request Checklist

Before submitting your PR, ensure:

- [ ] Code follows project conventions
- [ ] Tests are added/updated and passing (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] All checks pass (`make check`)
- [ ] Coverage is maintained/improved (`make coverage`)
- [ ] Documentation is updated
- [ ] Commit messages are clear and descriptive
- [ ] PR description explains the changes

## Testing Guidelines

### Unit Tests

- Test files should be named `*_test.go`
- Use table-driven tests for multiple cases
- Test error conditions and edge cases
- Mock external dependencies when needed

### Coverage

- Minimum coverage target: 75%
- Focus on testing critical paths
- Don't sacrifice quality for coverage numbers

### Running Tests

```bash
# All tests
make test

# With race detection
make test

# Verbose output
make test-v

# Coverage report
make coverage

# HTML coverage report
make coverage-html
```

## Package-Specific Guidelines

### util/

- This is the foundation package
- Changes here affect all other packages
- Extra care needed for backward compatibility

### crypto/

- Security-sensitive code
- Must include tests with known test vectors
- Document security assumptions

### fx/

- Expression evaluation engine
- Add tests for new operators
- Update operator registry properly

## Code Review Process

1. Automated checks run on PR creation
2. Maintainer review
3. Address feedback
4. Approval and merge

## Questions?

- Open an issue for discussion
- Check existing issues/PRs
- Review documentation

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
