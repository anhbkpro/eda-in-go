# Contributing

## Commit Message Convention

This project uses [Conventional Commits](https://conventionalcommits.org/) for commit messages. Commit messages are automatically linted using [commitlint](https://commitlint.js.org/).

### Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code (white-space, formatting, etc.)
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools and libraries
- `perf`: A code change that improves performance
- `ci`: Changes to CI configuration files and scripts
- `build`: Changes that affect the build system or external dependencies
- `revert`: Reverts a previous commit

### Examples

```
feat: add payment service commands to grpc-requests.sh

- Add authorize-payment, confirm-payment, create-invoice, adjust-invoice, pay-invoice, and cancel-invoice commands
- Update script header and help documentation

fix: update database connection configuration

- Change PG_CONN environment variable to use correct database URL
- Fix connection issues when running locally

docs: update API documentation

- Add missing parameter descriptions
- Fix typos in endpoint documentation
```

### Commit Linting

Commits are automatically validated against the conventional commit format. If your commit message doesn't follow the format, the commit will be rejected with an error message explaining what's wrong.

To fix a commit message after a failed commit:

```bash
git commit --amend
```

## Development Workflow

### Setup

1. Clone the repository
2. Install Go dependencies and generate code:
   ```bash
   make install-tools
   make generate
   ```
3. Install Node.js dependencies and setup git hooks:
   ```bash
   npm install
   ```
   This installs commitlint and configures git to use the husky hooks in this directory.

### Making Changes

1. Create a feature branch from `main`
2. Make your changes
3. Test your changes
4. Commit your changes with a conventional commit message
5. Push your branch and create a pull request

### Git Hooks

This project uses [Husky](https://typicode.github.com/husky/) to run git hooks:

- `commit-msg`: Validates commit messages against conventional commit format

To bypass hooks (not recommended):
```bash
git commit --no-verify
```

## Code Style

- Follow Go conventions and best practices
- Use `go fmt` to format code
- Run `go vet` to check for common issues
- **Import Organization**: Keep imports organized with proper grouping:
  ```go
  import (
      // std libs
      "context"
      "fmt"

      // external libs
      "github.com/rs/zerolog"

      // alias
      pg "eda-in-golang/internal/postgres"

      // internal libs
      "eda-in-golang/internal/ddd"
  )
  ```
- **Import Checking**: Use `make check-imports-go` (preferred, faster) or `make check-imports` to validate import organization
- Ensure all tests pass before committing
