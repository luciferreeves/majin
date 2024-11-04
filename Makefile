# Variables
BINARY_NAME=majin
MAIN_PATH=cmd/majin/main.go
ENV_FILE=.env

.PHONY: all build clean dev setup templ tailwind tailwind-watch

setup:
	@echo "Setting up Majin Web Mail..."
	@if [ ! -f $(ENV_FILE) ]; then \
		cp .env.example $(ENV_FILE); \
		echo "Created .env file from example"; \
	fi
	go mod tidy
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/air-verse/air@latest
	npm install
	@echo "Setup complete! Edit .env file with your configuration."

build: templ tailwind
	@echo "Building Majin Web Mail..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/*
	rm -rf templates/*_templ.go

dev:
	@echo "Starting air for Go live reload..."
	air

templ:
	@echo "Generating templates..."
	templ generate

tailwind-watch:
	@echo "Watching Tailwind CSS..."
	npm run dev

tailwind:
	@echo "Building Tailwind CSS..."
	npm run build

prod: build
	@echo "Starting production server..."
	ENV=production ./bin/$(BINARY_NAME)
