help:
	@echo ""
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  install  Install vendors"
	@echo "  build    Builds for production"
	@echo "  lint     Runs stylus and es lint"
	@echo ""

install:
	@echo '> Installing Node.js packages...'
	@npm install
	@echo 'Node.js packages successfully installed'

build:
	@echo '> Building app for production...'
	@./node_modules/.bin/webpack
	@echo 'App successfully built!'

lint:
	@./node_modules/.bin/eslint src/*
	@./node_modules/.bin/stylint src/styles


.PHONY: install build lint
