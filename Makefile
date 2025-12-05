.PHONY: scaffold

# Copy scaffold to a fresher's directory.
# Usage: make scaffold FRESHER=chien
# Will NOT overwrite existing files (uses cp -n).
scaffold:
ifndef FRESHER
	$(error FRESHER is required. Usage: make scaffold FRESHER=chien)
endif
	@echo "Copying scaffold to $(FRESHER)/..."
	@cp -Rn resources/scaffold/todos/ $(FRESHER)/todos/ 2>/dev/null; true
	@cp -Rn resources/scaffold/todos-bff/ $(FRESHER)/todos-bff/ 2>/dev/null; true
	@echo "Done. Scaffold copied to $(FRESHER)/todos/ and $(FRESHER)/todos-bff/"
	@echo "Existing files were NOT overwritten."
