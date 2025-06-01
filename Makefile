SHELL := /bin/bash

.PHONY: env

env:
	@set -a; \
	source .env; \
	set +a; \
	echo "✅ Environment loaded. You can now run Goose commands:"; \
	echo ""; \
	echo "  goose status"; \
	echo "  goose up"; \
	echo "  goose down"; \
	echo ""; \
	$$SHELL
