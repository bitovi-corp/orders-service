<!--
  Sync Impact Report - Constitution Update
  ========================================
  Version Change: 1.0.0 → 1.0.1
  
  Modified Principles: Principle V - Standard Error Handling
  - Updated: writeErrorResponse function signature documentation to include 'details' parameter
  - Clarified: Full function signature is writeErrorResponse(w, statusCode, code, message, details)
  - Added: Parameter descriptions for each argument
  
  Templates Requiring Updates:
  ✅ spec-template.md - No changes needed (principle clarification only)
  ✅ plan-template.md - No changes needed (principle clarification only)
  ✅ tasks-template.md - No changes needed (principle clarification only)
  
  Follow-up TODOs: None - documentation update only
-->

# Example Go Server Constitution

## Core Principles

### I. Contract-First Development (NON-NEGOTIABLE)
The OpenAPI specification (`api/openapi.yaml`) is the source of truth for all API contracts.

- All new endpoints MUST be defined in the OpenAPI specification before implementation
- API request/response schemas MUST match OpenAPI definitions exactly
- HTTP status codes MUST align with OpenAPI response specifications
- Breaking changes MUST be documented in the OpenAPI spec with version updates

**Rationale**: Contract-first development ensures API consistency, enables client code generation, provides living documentation, and prevents implementation drift from intended design.

### II. Standard Go Project Layout
Follow standard Go project structure with clear separation of concerns.

- `/cmd` contains application entry points
- `/internal` contains private application code not intended for external import
- `/api` contains API contracts and specifications
- `/tests/integration` contains end-to-end integration tests
- Handlers delegate to services; services contain business logic
- Models define data structures without behavior

**Rationale**: Standard layout improves code navigability, enforces encapsulation boundaries, and aligns with Go community best practices.

### III. Test Coverage & Isolation
Comprehensive testing with proper isolation ensures reliability.

- Unit tests MUST be located alongside source files (`*_test.go`)
- Integration tests MUST be located in `/tests/integration`
- Each test MUST reset mock data to ensure isolation
- Table-driven test patterns SHOULD be used where applicable
- All packages MUST be testable via `go test ./...`

**Rationale**: Test isolation prevents flaky tests and ensures deterministic behavior. Co-located unit tests reduce cognitive load and encourage test writing.

### IV. Middleware Composition
Protected routes MUST use standardized middleware composition.

- All routes (except `/health`) MUST use `AuthMiddleware` for authentication
- All routes MUST use `LoggingMiddleware` for request/response logging
- Middleware MUST be composable: `LoggingMiddleware(AuthMiddleware(handler))`
- New middleware MUST follow the same composition pattern

**Rationale**: Consistent middleware application ensures security, observability, and maintainability. Composition pattern keeps code DRY and testable.

### V. Standard Error Handling
All HTTP errors MUST use consistent response format.

- Use `writeErrorResponse(w, statusCode, code, message, details)` pattern
  - `statusCode`: HTTP status code (e.g., 400, 404, 500)
  - `code`: Descriptive error code in uppercase (e.g., "INVALID_USER_ID")
  - `message`: Human-readable error message
  - `details`: Optional additional error details (use empty string if not needed)
- Error codes MUST be descriptive and uppercase (e.g., "INVALID_USER_ID")
- HTTP status codes MUST follow REST conventions
- Errors MUST be returned up the call stack and handled in handlers

**Rationale**: Standardized error responses improve API usability, enable consistent client-side error handling, and simplify debugging.

## Technical Standards

### Language & Dependencies
- Go version: 1.25.5 or higher
- Minimize external dependencies; prefer standard library
- All dependencies MUST be tracked in `go.mod`
- Dependency updates MUST include rationale and testing

### Code Conventions
- Keep handlers thin - delegate to services
- Models define data structures, not behavior
- Use standard `log` package with prefixes for logging
- Follow Go formatting standards (`gofmt`, `goimports`)

### Performance & Scale
- Target: Handle typical e-commerce workloads (<100 concurrent users for demo)
- Mock data storage acceptable for demonstration purposes
- Optimize for code clarity over premature performance optimization

## Development Workflow

### Feature Implementation Process
1. Define or update OpenAPI specification first
2. Create or update models if needed
3. Implement handler with proper middleware composition
4. Add service layer business logic
5. Wire up route in `cmd/server/main.go`
6. Write and run tests (unit + integration as appropriate)
7. Verify against OpenAPI spec

### Testing Requirements
- All new code SHOULD include unit tests
- Integration tests MUST be added for:
  - New endpoint workflows
  - Multi-step user journeys
  - Cascade operations (e.g., user deletion canceling orders)
  - State transitions (e.g., order status changes)

### Documentation Requirements
- README MUST be updated for new endpoints or features
- OpenAPI spec MUST include descriptions and examples
- Code comments SHOULD explain "why" not "what"

## Governance

This constitution supersedes all other development practices and guidelines. All feature specifications, implementation plans, and code reviews MUST verify compliance with these principles.

### Amendment Process
- Amendments require documented rationale and team approval
- Version increments follow semantic versioning:
  - **MAJOR**: Backward incompatible principle removals or redefinitions
  - **MINOR**: New principles or materially expanded guidance
  - **PATCH**: Clarifications, wording refinements, non-semantic changes
- All dependent templates and documentation MUST be updated with constitution changes

### Compliance
- All pull requests MUST demonstrate adherence to core principles
- Deviations from principles MUST be explicitly justified and documented
- Use `.specify/templates` for feature planning and task management aligned with these principles

**Version**: 1.0.1 | **Ratified**: 2026-01-09 | **Last Amended**: 2026-01-09
