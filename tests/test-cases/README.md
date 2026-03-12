# Test Case Template Guide

## Test Case Structure

Each test case must follow this JSON structure:

### Required Fields

- **id** (string): Unique identifier for the test case
  - Format: `<module-prefix>-<number>`
  - Example: `auth-001`, `dash-002`

- **module** (string): Module name this test belongs to
  - Must match the directory name
  - Example: `authentication`, `dashboard`

- **scenario** (string): Brief description of what is being tested
  - Should be clear and concise
  - Example: "User login with valid credentials"

- **steps** (array of strings): Detailed test steps
  - Each step should be actionable
  - Write in imperative mood
  - Example: ["Navigate to login page", "Enter username", "Click login"]

- **expectedResults** (array of strings): Expected outcomes
  - Should match the number of steps or key checkpoints
  - Be specific about what should happen
  - Example: ["User is redirected to dashboard", "Success message shown"]

- **actualResults** (array): Actual test results (populated during execution)
  - Leave as empty array `[]` when creating test cases
  - Will be filled automatically during test execution

- **status** (string): Test execution status
  - Valid values: `"not_run"`, `"pass"`, `"fail"`, `"blocked"`
  - Use `"not_run"` for new test cases

- **lastExecuted** (string|null): Timestamp of last execution
  - Set to `null` for new test cases
  - Will be updated automatically during execution

- **executionTime** (number|null): Execution time in milliseconds
  - Set to `null` for new test cases
  - Will be updated automatically during execution

### Optional Fields

- **tags** (array of strings): Test categorization tags
  - Priority: `"critical"`, `"high"`, `"medium"`, `"low"`
  - Type: `"smoke"`, `"regression"`, `"integration"`, `"negative"`
  - Example: `["critical", "smoke"]`

## Example Test Cases

### Example 1: Login Test

```json
{
  "id": "auth-001",
  "module": "authentication",
  "scenario": "User login with valid credentials",
  "steps": [
    "Navigate to login page",
    "Enter valid username",
    "Enter valid password",
    "Click login button"
  ],
  "expectedResults": [
    "User is redirected to dashboard",
    "User profile is displayed in header",
    "Authentication token is stored",
    "No error messages are shown"
  ],
  "actualResults": [],
  "status": "not_run",
  "lastExecuted": null,
  "executionTime": null,
  "tags": ["critical", "smoke"]
}
```

### Example 2: Create Record Test

```json
{
  "id": "menu-001",
  "module": "menu-manajemen",
  "scenario": "Create new menu item",
  "steps": [
    "Navigate to Menu Management",
    "Click add menu button",
    "Enter menu name and description",
    "Set price",
    "Save menu item"
  ],
  "expectedResults": [
    "Menu form is displayed",
    "Menu item is created successfully",
    "Success message is shown",
    "New menu item appears in list"
  ],
  "actualResults": [],
  "status": "not_run",
  "lastExecuted": null,
  "executionTime": null,
  "tags": ["critical"]
}
```

### Example 3: Negative Test

```json
{
  "id": "auth-002",
  "module": "authentication",
  "scenario": "User login with invalid credentials",
  "steps": [
    "Navigate to login page",
    "Enter invalid username",
    "Enter invalid password",
    "Click login button"
  ],
  "expectedResults": [
    "User remains on login page",
    "Error message 'Invalid credentials' is displayed",
    "No authentication token is stored",
    "Login form is still visible"
  ],
  "actualResults": [],
  "status": "not_run",
  "lastExecuted": null,
  "executionTime": null,
  "tags": ["high", "negative"]
}
```

## Best Practices

### Writing Good Test Cases

1. **Be Specific**: Clearly describe what action to take and what to expect
2. **Be Atomic**: Each test case should test one specific scenario
3. **Be Independent**: Tests should not depend on other tests
4. **Be Repeatable**: Tests should produce same results when run multiple times

### Naming Conventions

- **Test ID**: Use module prefix + sequential number
  - Good: `auth-001`, `dash-002`, `menu-003`
  - Bad: `test1`, `mytest`, `login`

- **Scenario**: Use present tense, be descriptive
  - Good: "User login with valid credentials"
  - Bad: "Login", "Test login"

- **Steps**: Use imperative mood
  - Good: "Click login button", "Enter username"
  - Bad: "Clicking login", "Username is entered"

### Tags Usage

- **Priority Tags**: Indicate test importance
  - `critical`: Must pass, blocks release
  - `high`: Important functionality
  - `medium`: Standard functionality
  - `low`: Nice to have

- **Type Tags**: Indicate test category
  - `smoke`: Quick sanity check
  - `regression`: Verify existing functionality
  - `integration`: Test module interactions
  - `negative`: Test error conditions
  - `validation`: Test input validation

## Adding New Test Cases

1. Navigate to module directory: `test-cases/<module>/`
2. Open `test-cases.json`
3. Add new test case object to the array
4. Follow the template structure
5. Ensure JSON is valid (no trailing commas, proper quotes)
6. Save the file

## Validation

Before running tests, validate your test cases:

```bash
node -e "console.log(JSON.parse(require('fs').readFileSync('test-cases/<module>/test-cases.json')))"
```

This will catch JSON syntax errors.

## Common Mistakes to Avoid

1. ❌ Missing required fields
2. ❌ Invalid status values
3. ❌ Empty steps or expectedResults arrays
4. ❌ Duplicate test IDs
5. ❌ Invalid JSON syntax (trailing commas, unquoted keys)
6. ❌ Module name doesn't match directory name
