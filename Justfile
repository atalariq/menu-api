set dotenv-load

default:
  @just --list

# Run the application
run:
  GIN_MODE=debug go run cmd/server/main.go

# Generate Swagger Docs
docs:
  swag init -g cmd/server/main.go --output docs

# Run Tests
test:
  go test ./test/... -v

# Clean build files
clean:
  rm -f menu.db
  rm -rf docs

# Run clean > docs > run in single command 
lazy: clean docs run
