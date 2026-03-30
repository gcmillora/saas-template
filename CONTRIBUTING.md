# Contributing

Thanks for your interest in contributing to this project!

## Getting Started

1. Fork the repo
2. Clone your fork
3. Run `./setup.sh` to get the dev environment running
4. Make your changes

## Bug Reports

Open an issue with:
- What you expected to happen
- What actually happened
- Steps to reproduce

## Feature Requests

Open an issue to discuss before building. This keeps scope manageable and avoids wasted effort.

## Pull Requests

- Bug fixes are always welcome
- Keep changes small and focused
- Follow existing patterns in the codebase
- Make sure `bash commands.sh lint` and `bash commands.sh test` pass in the backend
- Make sure `bun run build` passes in the frontend

I review PRs when I can — please be patient.

## Development Commands

```bash
# One-command setup
./setup.sh

# Backend (from backend/)
bash commands.sh webserver          # Start dev server
bash commands.sh lint               # Lint
bash commands.sh test               # Test
bash commands.sh openapi:codegen    # Regenerate API types

# Frontend (from frontend/)
bun run dev                         # Start dev server
bun run build                       # Type check + build
bun run lint                        # Lint
```
