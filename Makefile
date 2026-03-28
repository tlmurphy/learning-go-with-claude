# Run lesson demos for a module
# Usage: make lesson 01
#        make lesson 14
lesson:
	@module=$$(ls -d $(word 2,$(MAKECMDGOALS))-* 2>/dev/null | head -1); \
	if [ -z "$$module" ]; then \
		echo "Module $(word 2,$(MAKECMDGOALS)) not found"; \
		exit 1; \
	fi; \
	echo "Running demos for $$module"; \
	go test -v -run "TestDemo" ./$$module/ 2>&1 | sed -E '/^=== RUN   Test/{ s/^=== RUN   Test/\n--- /; }; /^--- PASS|^PASS|^ok /d'

# Run exercises tests for a module
# Usage: make test 01
#        make test 14
test:
	@module=$$(ls -d $(word 2,$(MAKECMDGOALS))-* 2>/dev/null | head -1); \
	if [ -z "$$module" ]; then \
		echo "Module $(word 2,$(MAKECMDGOALS)) not found"; \
		exit 1; \
	fi; \
	echo "Running tests for $$module"; \
	go test -v ./$$module/

# Catch the module number argument so make doesn't complain
%:
	@:

.PHONY: lesson test
