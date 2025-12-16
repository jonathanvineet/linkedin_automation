# Contributing Guide

Thank you for your interest in contributing to this educational project! While this is primarily a proof-of-concept for learning, improvements and suggestions are always welcome.

## Code of Conduct

### Ethical Guidelines

1. **Educational Purpose Only**: This project is for learning about automation and anti-detection techniques
2. **No Production Use**: Never encourage or assist with violating LinkedIn's ToS
3. **Responsible Disclosure**: If you find new detection methods, share responsibly
4. **Test Accounts Only**: All testing must be on dedicated test accounts

## How to Contribute

### Reporting Bugs

**Before submitting**:
- Check existing [issues](https://github.com/jonathanvineet/linkedin-automation/issues)
- Verify the bug with latest code
- Test with a clean environment

**When reporting**:
```markdown
### Bug Description
Clear description of the issue

### Steps to Reproduce
1. Start application
2. Navigate to...
3. Click on...
4. See error

### Expected Behavior
What should happen

### Actual Behavior
What actually happens

### Environment
- OS: Ubuntu 22.04
- Go version: 1.21.5
- Node version: 18.17.0
- Browser: Chrome 120

### Logs
```
Relevant log excerpts
```
```

### Suggesting Enhancements

**Enhancement Template**:
```markdown
### Feature Request

**Problem**: Describe the problem this solves

**Proposed Solution**: How should it work?

**Alternatives**: Other approaches considered

**Educational Value**: How does this help learning?
```

### Pull Requests

**Process**:

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-improvement`
3. **Make** your changes
4. **Test** thoroughly
5. **Commit** with clear messages
6. **Push** to your fork
7. **Submit** a pull request

**PR Guidelines**:
- One feature/fix per PR
- Include tests for new functionality
- Update documentation
- Follow existing code style
- Add comments for complex logic

### Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/linkedin-automation.git
cd linkedin-automation

# Add upstream remote
git remote add upstream https://github.com/jonathanvineet/linkedin-automation.git

# Install dependencies
make install

# Create feature branch
git checkout -b feature/my-feature

# Make changes and test
make test

# Commit
git commit -m "Add amazing feature"

# Push
git push origin feature/my-feature
```

## Coding Standards

### Go Code Style

**Follow Go idioms**:
```go
// Good
func (s *Session) Navigate(url string) error {
    if url == "" {
        return errors.New("URL cannot be empty")
    }
    return s.page.Navigate(url)
}

// Bad
func (s *Session) navigate(Url string) (err error) {
    if Url != "" {
        err = s.page.Navigate(Url)
    } else {
        err = errors.New("URL cannot be empty")
    }
    return
}
```

**Naming Conventions**:
- Exported: `PascalCase`
- Unexported: `camelCase`
- Constants: `UPPER_SNAKE_CASE`
- Interfaces: `er` suffix (e.g., `Reader`, `Writer`)

**Error Handling**:
```go
// Wrap errors with context
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

**Comments**:
```go
// Function comments start with function name
// Example: SendRequest sends a connection request to a LinkedIn profile.
func SendRequest(profileURL string, note string) error {
    // Implementation details
}
```

### TypeScript/React Style

**Component Structure**:
```typescript
interface ComponentProps {
    prop1: string;
    prop2?: number; // Optional props
}

export function Component({ prop1, prop2 = 0 }: ComponentProps) {
    const [state, setState] = useState<Type>(initialValue);
    
    useEffect(() => {
        // Side effects
    }, [dependencies]);
    
    return (
        <div className="...">
            {/* JSX */}
        </div>
    );
}
```

**Naming**:
- Components: `PascalCase`
- Functions/variables: `camelCase`
- Constants: `UPPER_SNAKE_CASE`
- Types/Interfaces: `PascalCase`

### File Organization

**Go**:
```
internal/
├── module/
│   ├── module.go        # Main logic
│   ├── module_test.go   # Tests
│   └── types.go         # Type definitions
```

**React**:
```
src/
├── components/
│   ├── Component.tsx
│   └── ui/              # Reusable UI components
├── lib/
│   └── utility.ts       # Utility functions
└── pages/
    └── Page.tsx         # Page components
```

## Testing

### Go Tests

**Unit Test Example**:
```go
func TestMouseBezierCurve(t *testing.T) {
    mouse := NewMouseMovement()
    start := Point{X: 0, Y: 0}
    end := Point{X: 100, Y: 100}
    
    points := mouse.GenerateBezierCurve(start, end, 10)
    
    if len(points) != 10 {
        t.Errorf("Expected 10 points, got %d", len(points))
    }
    
    if points[0] != start {
        t.Error("First point should equal start")
    }
    
    if points[len(points)-1] != end {
        t.Error("Last point should equal end")
    }
}
```

**Run Tests**:
```bash
go test ./...              # All tests
go test -v ./internal/stealth  # Specific package
go test -cover ./...       # With coverage
```

### Integration Tests

**API Test Example**:
```bash
# Start server
./bin/automation &

# Test endpoints
curl -X GET http://localhost:8090/api/status
curl -X POST http://localhost:8090/api/start

# Cleanup
pkill automation
```

## Documentation

### Code Comments

**What to comment**:
- ✅ Complex algorithms
- ✅ Non-obvious decisions
- ✅ Public APIs
- ✅ Stealth techniques rationale

**What not to comment**:
- ❌ Obvious code
- ❌ Every line
- ❌ Bad code (refactor instead)

**Example**:
```go
// GenerateBezierCurve creates a natural mouse movement path using
// quadratic Bézier curve interpolation. This prevents detection of
// straight-line mouse movements which are a telltale sign of automation.
//
// The control point is randomized to create varied curves, and occasional
// overshoot is added to mimic human imprecision.
func (m *MouseMovement) GenerateBezierCurve(start, end Point, steps int) []Point {
    // ... implementation
}
```

### README Updates

**When to update**:
- New features added
- Configuration changes
- New dependencies
- Setup process changes

## Areas for Contribution

### High Priority

- ✅ **Improved Selectors**: LinkedIn UI changes frequently
- ✅ **Enhanced Stealth**: New detection evasion techniques
- ✅ **Better Error Handling**: Graceful degradation
- ✅ **Test Coverage**: More unit and integration tests

### Medium Priority

- 📊 **Analytics Dashboard**: Better metrics visualization
- 🔧 **Configuration UI**: Web-based config editor
- 📝 **Templates**: More message templates
- 🌐 **Internationalization**: Multi-language support

### Future Ideas

- 🐳 **Docker Support**: Containerization
- ☸️ **Kubernetes**: Deployment configs
- 📱 **Mobile UI**: Responsive dashboard
- 🤖 **AI Integration**: GPT-powered message generation

## Review Process

### Checklist

Before submitting PR, verify:

- [ ] Code follows style guide
- [ ] Tests pass (`make test`)
- [ ] Documentation updated
- [ ] Commits are clean and descriptive
- [ ] No sensitive data in commits
- [ ] Changes are educational in nature
- [ ] License headers present (if new files)

### Review Criteria

Reviewers will check:

1. **Code Quality**: Clean, readable, idiomatic
2. **Functionality**: Works as intended
3. **Testing**: Adequate test coverage
4. **Documentation**: Clear comments and README
5. **Ethics**: Aligns with educational purpose
6. **Security**: No credentials or secrets
7. **Performance**: No obvious inefficiencies

## Questions?

- 📧 **Email**: Open an issue for questions
- 💬 **Discussions**: Use GitHub Discussions
- 🐛 **Bugs**: Create an issue with details
- 💡 **Ideas**: Share in Discussions

---

## License

By contributing, you agree that your contributions will be licensed under the same MIT License that covers this project.

---

**Thank you for helping make this educational project better!** 🎓

*Remember: This is for learning only. Always respect platform terms of service.*
