---
name: "project-standard-checker"
description: "Validates project structure, file naming, documentation location, and code standards. Invoke when creating new files, organizing project structure, or checking compliance with project norms."
---

# Project Standard Checker

This skill helps maintain project organization and coding standards for the FG-ABYSS project by validating file structure, naming conventions, documentation placement, and development practices.

## When to Invoke

Invoke this skill when:
1. **Creating new files or directories** - Validate naming and location
2. **Organizing project structure** - Ensure proper module placement
3. **Reviewing code compliance** - Check adherence to standards
4. **Setting up documentation** - Verify correct docs/ directory usage
5. **Preparing commits** - Validate commit message format
6. **Onboarding new team members** - Guide on project standards

## Core Responsibilities

### 1. Document Management Validation

**Rule**: All documentation MUST be stored in `docs/` directory (except README.md, LICENSE).

```bash
# ✅ Correct
/docs/architecture/system-design.md
/docs/testing/unit-test-report.md
/docs/development/getting-started.md

# ❌ Incorrect (in project root)
/PROJECT_REPORT.md
/TEST_PLAN.md
```

**Validation Checklist**:
- [ ] Document is in `docs/<category>/` directory
- [ ] Category is appropriate (architecture/, development/, testing/, api/, deployment/, reports/)
- [ ] File name follows naming conventions
- [ ] Document structure includes: title, overview, content, summary, references

### 2. File Naming Validation

**General Rules**:
- Use lowercase letters
- Use hyphens `-` as separators
- Semantic, descriptive names
- No camelCase, underscores, or spaces

**By File Type**:

```bash
# Go Files
✅ project_service.go
✅ webshell_handler_test.go
❌ ProjectService.go
❌ project-service.go
❌ project_service.go

# Vue Components
✅ ProjectsContent.vue
✅ CreateProjectModal.vue
❌ projects-content.vue
❌ projectsContent.vue

# Test Files
✅ project_service_test.go
✅ formatTime.test.ts
❌ test_project.go
❌ ProjectTest.go

# Documentation
✅ unit-test-report.md
✅ bugfix-webshell-restore.md
❌ UnitTestReport.md
❌ unit_test_report.md
```

### 3. Git Commit Message Validation

**Format**: `<type>(<scope>): <subject>`

**Type Values**:
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `style` - Code formatting (no logic change)
- `refactor` - Code refactoring
- `test` - Test additions/modifications
- `chore` - Build/config changes
- `perf` - Performance improvements
- `ci` - CI/CD changes

**Scope Values**:
- Backend: `service`, `repository`, `handler`, `entity`
- Frontend: `component`, `utils`, `view`, `api`
- General: `deps`, `config`, `readme`

**Examples**:

```bash
# ✅ Correct
feat(webshell): add batch delete support
fix(repository): correct GORM model specification
test(service): add unit tests for CRUD operations
docs(readme): update installation guide

# ❌ Incorrect
fixed bug
update code
Fix: Fixed webshell restore (past tense, capitalized)
```

**Subject Rules**:
- Use imperative mood (create, not created/creates)
- No capital first letter
- No period at the end
- Keep under 50 characters

### 4. Project Structure Validation

**Standard Structure**:

```
FG-ABYSS/
├── docs/                    # 📚 ALL documentation (except README)
│   ├── architecture/       # Architecture docs
│   ├── development/        # Development guides
│   ├── api/               # API documentation
│   ├── testing/           # Test reports
│   ├── deployment/        # Deployment guides
│   └── reports/           # Temporary reports
├── internal/               # Private business logic
│   ├── app/               # Application layer
│   │   ├── handlers/      # Request handlers
│   │   └── services/      # Business logic
│   ├── domain/            # Domain layer
│   │   ├── entity/        # Entities
│   │   └── repository/    # Repository interfaces
│   └── infrastructure/    # Infrastructure
│       └── repositories/  # Repository implementations
├── frontend/              # Frontend code
│   ├── src/
│   │   ├── components/   # Vue components
│   │   ├── views/        # Page views
│   │   └── utils/        # Utility functions
│   └── ...
├── tests/                 # Test files
│   ├── fixtures/         # Test data
│   └── integration/      # Integration tests
└── ...
```

**Module Responsibilities**:

- `handlers`: Request processing, parameter validation, response formatting
- `services`: Business logic, transactions, business rules
- `entity`: Business entities, validation logic
- `repository` (interface): Data access contracts
- `repositories` (implementation): Database operations, ORM mapping

### 5. Code Development Standards

**Go Standards**:
```go
// ✅ Proper error handling
func (s *Service) Create(name string) error {
    if err := validate(name); err != nil {
        return fmt.Errorf("validate: %w", err)
    }
    // ...
}

// ✅ Exported functions have comments
// CreateProject creates a new project
func (s *Service) CreateProject(name string) (*entity.Project, error) {
    // ...
}
```

**Vue/TypeScript Standards**:
```typescript
// ✅ Proper type definitions
interface Props {
  projectId?: string
  showDeleted?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  projectId: '',
  showDeleted: false
})

// ✅ Proper event naming
emit('project-created', data)
emit('update:loading', false)
```

### 6. Testing Standards

**Test File Location**:
```bash
# Go: Same directory as tested file
internal/app/services/
├── project_service.go
└── project_service_test.go

# Frontend: Same directory or tests/
frontend/src/utils/
├── formatTime.ts
└── formatTime.test.ts
```

**Test Naming**:
```go
// ✅ Clear naming: Test<Function>_<Scenario>
func TestProjectService_Create_Success(t *testing.T)
func TestProjectService_Create_NameExists(t *testing.T)
func TestWebShell_Validate_URLEmpty(t *testing.T)
```

**Coverage Requirements**:
- Service layer: ≥ 40%
- Entity validation: ≥ 80%
- Handler layer: ≥ 30%
- Repository layer: ≥ 50%

## Usage Examples

### Example 1: Creating New Documentation

**User**: "I need to create a bug fix report for the webshell restore issue"

**Response**:
```
I'll help you create a properly structured bug fix report.

✅ **Location**: `docs/testing/bugfixes/webshell-restore-fix.md`
✅ **Naming**: Uses lowercase with hyphens
✅ **Category**: testing/bugfixes/ for bug reports

The document should include:
1. Problem description
2. Root cause analysis
3. Fix details
4. Testing results
5. Prevention measures

Would you like me to create a template?
```

### Example 2: Validating File Creation

**User**: "I'm creating a new service for user authentication"

**Response**:
```
For a new user authentication service, follow these standards:

✅ **File Location**: `internal/app/services/user_auth_service.go`
✅ **Test File**: `internal/app/services/user_auth_service_test.go`
✅ **Naming**: lowercase with underscores for Go files

Structure:
```go
package services

// UserService handles user authentication
type UserService struct {
    userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository) *UserService {
    return &UserService{userRepo: repo}
}
```

Don't forget to:
- [ ] Add unit tests (target: ≥40% coverage)
- [ ] Add documentation in `docs/api/user-auth.md`
- [ ] Update commit message: `feat(service): add user authentication service`
```

### Example 3: Validating Commit Message

**User**: "git commit -m 'fixed the webshell restore bug'"

**Response**:
```
❌ **Commit message needs correction**

Issues:
- Uses past tense "fixed" instead of imperative "fix"
- Missing scope
- Too vague

✅ **Correct format**:
```bash
git commit -m "fix(webshell): restore deleted webshell functionality"
```

Or with more detail:
```bash
git commit -m "fix(repository): correct GORM model in Recover method

- Add .Model(&entity.WebShell{}) to Recover()
- Add .Unscoped() to FindByID() for soft delete support

Fixes #78"
```

Key rules:
- Use imperative mood (fix, not fixed)
- Add scope in parentheses (webshell, repository, etc.)
- No capital first letter
- No period at end
```

### Example 4: Project Reorganization

**User**: "I want to move all my test reports to the proper location"

**Response**:
```
I'll help you reorganize test reports according to project standards.

**Current** → **Target Location**:

/TEST_REPORT.md → /docs/testing/unit-test-report.md
/BUGFIX_REPORT.md → /docs/testing/bugfixes/webshell-restore-fix.md
/architecture/ANALYSIS.md → /docs/architecture/unit-test-analysis.md

**Commands to execute**:
```bash
mkdir -p docs/testing/bugfixes
mkdir -p docs/architecture

mv TEST_REPORT.md docs/testing/unit-test-report.md
mv BUGFIX_REPORT.md docs/testing/bugfixes/webshell-restore-fix.md
mv docs/architecture/ANALYSIS_REPORT.md docs/architecture/unit-test-analysis.md
```

**Naming improvements**:
- Use lowercase: `ANALYSIS_REPORT.md` → `unit-test-analysis.md`
- Use hyphens: `unit_test_report.md` → `unit-test-report.md`
```

## Quick Reference

### Validation Checklist

Before committing code, verify:

- [ ] **File Naming**: Lowercase, hyphens, semantic
- [ ] **Document Location**: In `docs/<category>/` directory
- [ ] **Test Coverage**: Meets minimum requirements
- [ ] **Commit Message**: `<type>(<scope>): <subject>` format
- [ ] **Code Structure**: Follows layer responsibilities
- [ ] **No Debug Code**: Removed console.log, fmt.Println debugging
- [ ] **No Sensitive Data**: No passwords, keys, tokens

### Common Violations & Fixes

```bash
# ❌ Document in root
/PROJECT_PLAN.md
# ✅ Fix: Move to docs
/docs/development/project-plan.md

# ❌ Wrong naming
/internal/app/ProjectService.go
# ✅ Fix: Rename
/internal/app/services/project_service.go

# ❌ Bad commit
git commit -m "fixed stuff"
# ✅ Fix: Use conventional format
git commit -m "fix(service): resolve null pointer in create method"

# ❌ Test in wrong location
/tests/project_test.go
# ✅ Fix: Move next to source
/internal/app/services/project_service_test.go
```

### Command Shortcuts

```bash
# Create properly named documentation
touch docs/<category>/<descriptive-name>.md

# Run all tests
task test              # Go
cd frontend && npm run test  # Frontend

# Check test coverage
task test:coverage
go test ./... -cover

# Validate project structure
# (Manual check against section 4)
```

## Related Skills

- **code-reviewer**: Reviews code quality and best practices
- **test-generator**: Generates unit test templates
- **git-commit-helper**: Assists with commit message formatting

## References

- [FG-ABYSS Project Specification](file://d:\Go\FG-ABYSS\docs\development\project-optimization-specification.md)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Vue.js Style Guide](https://vuejs.org/style-guide/)
