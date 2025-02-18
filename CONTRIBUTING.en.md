# Contributing to K8s Resource Analyzer

[🇧🇷 Portuguese Version](CONTRIBUTING.md)

## Table of Contents
1. [How to Contribute](#how-to-contribute)
2. [Development Process](#development-process)
3. [Release Process](#release-process)
4. [Code Standards](#code-standards)

## How to Contribute

1. Fork the project
2. Create a branch for your feature (`git checkout -b feature/feature-name`)
3. Commit your changes (`git commit -m 'type: description'`)
4. Push to the branch (`git push origin feature/feature-name`)
5. Open a Pull Request

## Development Process

### Branches
- `main`: Production code
- `feature/*`: New features
- `bugfix/*`: Bug fixes
- `hotfix/*`: Urgent production fixes

### Commits
We use Conventional Commits:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting
- `refactor`: Refactoring
- `test`: Tests
- `chore`: Maintenance

## Release Process

### Prerequisites
- Be a project maintainer
- Have access to repository secrets
- Have write permission on Docker Hub

### Release Checklist
1. Ensure all tests pass
2. Verify documentation is up to date
3. Review CHANGELOG.md
4. Check if all dependencies are updated

### Creating a Release

1. Update main:
   ```bash
   git checkout main
   git pull origin main
   ```

2. Create and publish tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

3. Monitor release workflow in GitHub Actions:
   - Binary build
   - Docker image publication
   - GitHub release creation

### Versioning
MAJOR.MINOR.PATCH:
- MAJOR: Incompatible changes
- MINOR: New features
- PATCH: Bug fixes

### Post-Release
1. Verify Docker image was published
2. Validate release documentation on GitHub
3. Communicate the new version to the team

## Code Standards

### Go
- Use `gofmt` for formatting
- Follow Go conventions
- Maintain 80% test coverage
- Document public functions

### Documentation
- Keep README.md up to date
- Document API changes
- Update Swagger when necessary

### Quality
- All tests must pass
- No linter warnings
- Keep cyclomatic complexity low 