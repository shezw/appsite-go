# Development Workflow

This document guides the autonomous development agent through the project lifecycle.

## 1. Continuous Development Loop
The agent must repeatedly perform the following cycle until all modules in `development_order.md` are marked as `Completed`.

1.  **Read Status**: Check `development_order.md` for the first module marked `Pending` or `In Progress`.
2.  **Execute Phase**: Perform the **Development Unit Cycle** for that module.
3.  **Update Status**: Mark the module as `Completed` in `development_order.md` only when the **Definition of Done** is met.
4.  **Repeat**: Move to the next module.

## 2. Definition of Done (DoD)
A module is considered `Completed` only when:
*   [x] **Code Implemented**: All features described are implemented in `.go` files.
*   [x] **Tests Passed**: `go test ./...` passes for the module.
*   [x] **Coverage Met**: Test coverage is > 70%.
*   [x] **Documentation**: Code is commented according to `author.md` standards.

## 3. Development Unit Cycle
For each module, strictly follow this sequence:

1.  **Plan**: Identify files to create/edit.
2.  **Code (Commit)**:
    *   Write implementation code.
    *   Add file header from `author.md`.
    *   *Commit*: `feat(module): implement [feature name]`
3.  **Test Write (Commit)**:
    *   **Location**: All tests must be placed in the `tests/` directory, mirroring the source structure (e.g., `src/pkg/utils/file.go` -> `tests/pkg/utils/file/file_test.go`).
    *   **Type**: Use black-box testing (test public interfaces).
    *   Create `_test.go` files.
    *   *Commit*: `test(module): add unit tests for [feature]`
4.  **Run Test**: Execute `go test -coverpkg=./... -v ./tests/...`.
5.  **Refine (Loop)**:
    *   If failed -> Fix code -> *Commit*: `fix(module): resolve test failure [reason]`
    *   If coverage low -> Add cases -> *Commit*: `test(module): improve coverage to [X]%`
6.  **Finalize**:
    *   Verify DoD.
    *   Update `development_order.md` status to `Completed`.
    *   *Commit*: `chore(status): mark [module ID] as completed`

## 4. Commit Message Standards
*   **Format**: `type(scope): description`
*   **Types**: `feat`, `fix`, `test`, `chore`, `docs`, `refactor`.
*   **Fixes**: Must describe the bug found.
*   **Todo**: If a feature is postponed, add `// @todo [Explanation]` in code and mention in commit msg.

## 5. File Header & Author Rules
Refer to `author.md` in the project root.
*   **New Files**: Must include the standard copyright header.
*   **Modifications**: Use `// @updated_by {date} {alias} {email}` for significant changes.
