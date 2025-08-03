# GoOse Code Improvement Plan

## Executive Summary

This document outlines a comprehensive plan to improve the GoOse codebase, focusing on code organization, error handling, logging, documentation, and testing. The analysis reveals several critical areas that require attention to make the project more maintainable, robust, and production-ready.

## Current State Analysis

### Project Structure ✅ (Recently Improved)
- **Good**: Standard Go project layout with `pkg/`, `internal/`, and `cmd/` directories
- **Good**: Separation of concerns between public API and internal implementation
- **Good**: CLI interface with Cobra/Viper integration

### Critical Issues Identified

#### 1. Import Cycles and Broken Dependencies 🚨 **HIGH PRIORITY**
- Multiple compilation errors due to circular imports between packages
- Missing function implementations across packages
- Undefined types and methods causing build failures

#### 2. Inconsistent Error Handling 🔴 **HIGH PRIORITY**
- Mixed use of `log.Printf` vs proper error handling
- No structured error types
- Silent failures in several places
- Missing context in error messages

#### 3. Poor Logging Practices 🔴 **HIGH PRIORITY**
- Hardcoded `log.Printf` statements throughout the codebase
- No log levels (debug, info, warn, error)
- No structured logging
- Debug logs mixed with production code

#### 4. Minimal Documentation 🔴 **HIGH PRIORITY**
- Many functions lack proper GoDoc comments
- No architectural documentation
- Missing examples and usage guides
- No API documentation

#### 5. Limited Test Coverage 🔴 **HIGH PRIORITY**
- Only basic tests in `pkg/goose/` directory
- No tests for internal packages
- No integration tests
- No benchmark tests

## Improvement Plan

### Phase 1: Critical Fixes (Week 1-2)

#### 1.1 Fix Import Cycles and Build Issues
**Priority**: Critical
**Effort**: 3-5 days

**Tasks**:
- [ ] Resolve circular dependencies between `pkg/goose` and `internal/` packages
- [ ] Implement missing functions and methods:
  - `StopWords.stopWordsCount()`
  - `Parser.name()`, `Parser.setAttr()`, `Parser.removeNode()`
  - `NormaliseCharset()`, `UTF8encode()`
  - `NewExtractor()`, `NewCleaner()`, etc.
- [ ] Fix undefined types (Configuration, Article) in internal packages
- [ ] Ensure all packages compile successfully

**Files to Fix**:
```
internal/crawler/crawler.go          (13 errors)
internal/extractor/extractor.go      (15 errors)
internal/crawler/crawlershort.go     (similar issues)
internal/extractor/cleaner.go        (similar issues)
```

#### 1.2 Implement Proper Error Handling
**Priority**: High
**Effort**: 2-3 days

**Tasks**:
- [ ] Create custom error types for different error categories:
  ```go
  type ExtractionError struct {
      URL string
      Cause error
      Phase string
  }
  
  type NetworkError struct {
      URL string
      StatusCode int
      Cause error
  }
  
  type ParseError struct {
      Content string
      Position int
      Cause error
  }
  ```
- [ ] Replace `log.Printf` error statements with proper error returns
- [ ] Add error wrapping with context using `fmt.Errorf` or `pkg/errors`
- [ ] Implement error handling in CLI commands

### Phase 2: Logging and Observability (Week 3)

#### 2.1 Implement Structured Logging
**Priority**: High
**Effort**: 2-3 days

**Tasks**:
- [ ] Replace standard `log` package with structured logging (recommend `slog` or `logrus`)
- [ ] Implement configurable log levels
- [ ] Add request tracing for debugging
- [ ] Remove debug `log.Printf` statements from production code

**Example Implementation**:
```go
// internal/logging/logger.go
package logging

import (
    "log/slog"
    "os"
)

type Logger interface {
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
}

func NewLogger(level string) Logger {
    // Implementation with configurable levels
}
```

#### 2.2 Add Metrics and Monitoring
**Priority**: Medium
**Effort**: 1-2 days

**Tasks**:
- [ ] Add basic metrics (extraction time, success/failure rates)
- [ ] Implement health check endpoints
- [ ] Add performance monitoring for large documents

### Phase 3: Documentation (Week 4)

#### 3.1 Code Documentation
**Priority**: High
**Effort**: 3-4 days

**Tasks**:
- [ ] Add comprehensive GoDoc comments to all public functions
- [ ] Document all struct fields and their purposes
- [ ] Add examples in documentation
- [ ] Document configuration options

**Example**:
```go
// ExtractFromURL fetches content from the specified URL and extracts the main article content.
// It automatically handles character encoding detection and content cleaning.
//
// Parameters:
//   - url: The target URL to extract content from
//
// Returns:
//   - *Article: Extracted article with cleaned content and metadata
//   - error: Any error encountered during extraction
//
// Example:
//   g := goose.New()
//   article, err := g.ExtractFromURL("https://example.com/article")
//   if err != nil {
//       return err
//   }
//   fmt.Println(article.CleanedText)
func (g Goose) ExtractFromURL(url string) (*Article, error)
```

#### 3.2 Project Documentation
**Priority**: High
**Effort**: 2-3 days

**Tasks**:
- [ ] Create comprehensive README.md
- [ ] Add architecture documentation
- [ ] Create API documentation
- [ ] Add configuration guide
- [ ] Create troubleshooting guide

**Documentation Structure**:
```
docs/
├── README.md              # Project overview
├── ARCHITECTURE.md        # System design
├── API.md                 # API reference
├── CONFIGURATION.md       # Configuration guide
├── TROUBLESHOOTING.md     # Common issues
├── CONTRIBUTING.md        # Development guide
└── examples/              # Usage examples
    ├── basic-usage.md
    ├── advanced-config.md
    └── custom-extractors.md
```

### Phase 4: Testing (Week 5-6)

#### 4.1 Unit Tests
**Priority**: High
**Effort**: 4-5 days

**Tasks**:
- [ ] Write comprehensive unit tests for all packages
- [ ] Test error conditions and edge cases
- [ ] Mock external dependencies (HTTP requests)
- [ ] Achieve minimum 80% code coverage

**Test Structure**:
```
internal/
├── extractor/
│   ├── extractor_test.go
│   ├── cleaner_test.go
│   └── testdata/
├── crawler/
│   ├── crawler_test.go
│   └── testdata/
└── utils/
    ├── charset_test.go
    ├── stopwords_test.go
    └── testdata/
```

#### 4.2 Integration Tests
**Priority**: Medium
**Effort**: 2-3 days

**Tasks**:
- [ ] Test complete extraction workflows
- [ ] Test CLI commands end-to-end
- [ ] Test with various website types
- [ ] Performance tests with large documents

#### 4.3 Benchmark Tests
**Priority**: Medium
**Effort**: 1-2 days

**Tasks**:
- [ ] Benchmark extraction performance
- [ ] Memory usage analysis
- [ ] Regression testing

### Phase 5: Code Quality and Refactoring (Week 7-8)

#### 5.1 Code Quality Improvements
**Priority**: Medium
**Effort**: 3-4 days

**Tasks**:
- [ ] Run `go vet`, `golint`, `staticcheck` and fix issues
- [ ] Implement consistent naming conventions
- [ ] Remove dead code and unused variables
- [ ] Add input validation
- [ ] Improve configuration validation

#### 5.2 Performance Optimizations
**Priority**: Medium
**Effort**: 2-3 days

**Tasks**:
- [ ] Profile memory usage and optimize allocations
- [ ] Optimize HTML parsing performance
- [ ] Add caching for frequently accessed data
- [ ] Implement concurrent processing where appropriate

#### 5.3 Security Improvements
**Priority**: Medium
**Effort**: 1-2 days

**Tasks**:
- [ ] Input sanitization for URLs
- [ ] Validate file paths for output
- [ ] Add rate limiting for HTTP requests
- [ ] Secure configuration handling

### Phase 6: Advanced Features (Week 9-10)

#### 6.1 Enhanced Configuration
**Priority**: Low
**Effort**: 2-3 days

**Tasks**:
- [ ] Add configuration validation
- [ ] Support for configuration profiles
- [ ] Runtime configuration updates
- [ ] Environment variable support

#### 6.2 Plugin System
**Priority**: Low
**Effort**: 3-4 days

**Tasks**:
- [ ] Design plugin interface for custom extractors
- [ ] Implement plugin loading mechanism
- [ ] Create example plugins
- [ ] Document plugin development

## Implementation Guidelines

### Code Standards

#### Error Handling Pattern
```go
// Bad
if err != nil {
    log.Printf("Error: %v", err)
    return nil
}

// Good
if err != nil {
    return nil, fmt.Errorf("failed to extract content from %s: %w", url, err)
}
```

#### Logging Pattern
```go
// Bad
log.Printf("Starting extraction for %s", url)

// Good
logger.Info("starting content extraction", 
    "url", url, 
    "timeout", config.Timeout,
    "user_agent", config.UserAgent)
```

#### Function Documentation Pattern
```go
// FunctionName does X by doing Y.
// It returns Z when W happens.
//
// Parameters:
//   - param1: description
//   - param2: description
//
// Returns:
//   - type1: description
//   - error: description of possible errors
//
// Example:
//   result, err := FunctionName(param1, param2)
//   if err != nil {
//       // handle error
//   }
func FunctionName(param1 type1, param2 type2) (type1, error)
```

### Testing Standards

#### Test Structure
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        {
            name:     "invalid input",
            input:    invalidInput,
            expected: nil,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            
            if tt.wantErr && err == nil {
                t.Error("expected error but got none")
                return
            }
            
            if !tt.wantErr && err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }
            
            if !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Success Metrics

### Code Quality Metrics
- [ ] 100% successful builds with `go build ./...`
- [ ] Zero critical issues from static analysis tools
- [ ] Minimum 80% test coverage
- [ ] All public APIs documented with examples

### Performance Metrics
- [ ] < 2 seconds extraction time for typical articles
- [ ] < 100MB memory usage for large documents
- [ ] Support for 100+ concurrent extractions

### Maintainability Metrics
- [ ] Clear separation of concerns
- [ ] Consistent error handling patterns
- [ ] Comprehensive documentation
- [ ] Easy onboarding for new developers

## Timeline Summary

| Phase | Duration | Priority | Deliverables |
|-------|----------|----------|-------------|
| 1 | 1-2 weeks | Critical | Fixed build issues, basic error handling |
| 2 | 1 week | High | Structured logging, monitoring |
| 3 | 1 week | High | Complete documentation |
| 4 | 2 weeks | High | Comprehensive test suite |
| 5 | 2 weeks | Medium | Code quality improvements |
| 6 | 2 weeks | Low | Advanced features |

**Total Estimated Time**: 8-10 weeks

## Immediate Next Steps

1. **Fix build issues** (Days 1-3)
   - Resolve import cycles
   - Implement missing functions
   - Ensure `go build ./...` succeeds

2. **Implement basic error handling** (Days 4-5)
   - Create error types
   - Replace log statements with error returns
   - Add error context

3. **Set up testing framework** (Week 2)
   - Create test structure
   - Write first batch of unit tests
   - Set up CI/CD pipeline

This plan provides a structured approach to transforming GoOse from its current state into a production-ready, maintainable, and well-documented Go project.
