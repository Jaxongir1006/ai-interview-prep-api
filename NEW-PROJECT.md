# Creating a New Project from This Blueprint

## Quick Start

```bash
# 1. Clone the blueprint
git clone https://github.com/rise-and-shine/go-enterprise-blueprint.git my-new-project

# 2. Enter the project directory
cd my-new-project

# 3. Remove the blueprint's git history and start fresh
rm -rf .git
git init
git add .
git commit -m "Initial commit from go-enterprise-blueprint"

# 4. Rename the Go module
#    Replace "my-org/my-new-project" with your actual module path
OLD_MODULE="go-enterprise-blueprint"
NEW_MODULE="hire-ready-api"

# macOS
find . -type f -name '*.go' -exec sed -i '' "s|${OLD_MODULE}|${NEW_MODULE}|g" {} +
sed -i '' "s|${OLD_MODULE}|${NEW_MODULE}|g" go.mod

# Linux
find . -type f -name '*.go' -exec sed -i "s|${OLD_MODULE}|${NEW_MODULE}|g" {} +
sed -i "s|${OLD_MODULE}|${NEW_MODULE}|g" go.mod

# 5. Tidy dependencies
go mod tidy

# 6. Verify everything compiles
go build ./...

# 7. Point to your new remote
git remote add origin https://github.com/my-org/my-new-project.git
git push -u origin main
```

## What to Customize

After cloning, review and adapt the following to your project:

| Item            | Location                       | Action                                           |
| --------------- | ------------------------------ | ------------------------------------------------ |
| Module name     | `go.mod`, all `*.go` imports   | Replaced in step 4 above                         |
| App config      | `config/`                      | Update app name, ports, secrets                  |
| README          | `README.md`                    | Rewrite for your project                         |
| AGENTS.md       | `AGENTS.md`                    | Update if your conventions differ                |
| Docker          | `Dockerfile`, `dev-infra.yaml` | Update image names and services                  |
| Migrations      | `migrations/`                  | Keep or remove existing migrations               |
| Example modules | `internal/modules/`            | Remove modules you don't need, keep as reference |
| Specs           | `docs/specs/`                  | Replace with your own module specs               |
| i18n            | `i18n/`                        | Update translations for your domain              |
| Scripts         | `scripts/`                     | Review and adapt helper scripts                  |
| This file       | `NEW-PROJECT.md`               | Delete it once you're set up                     |
