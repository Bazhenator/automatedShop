WORKDIR=$(shell pwd)
BINARY_PATH=$(WORKDIR)/third_party/bin

MIGRATE_DB_USER=user
MIGRATE_DB_PASS=pass
MIGRATE_DB_HOST=localhost
MIGRATE_DB_PORT=3306
MIGRATE_DB_NAME=visitors
DBDRIVER=mysql

.PHONY: install-local-goose
install-local-goose:
  	@GOBIN=$(BINARY_PATH) go install github.com/pressly/goose/v3/cmd/goose@v3.18.0
	@$(BINARY_PATH)/goose -version
	@if test -f "$(BINARY_PATH)/.goose.env"; then \
    	rm $(BINARY_PATH)/.goose.env; \
  	fi
	@touch $(BINARY_PATH)/.goose.env
	@echo "GOOSE_DRIVER=$(DBDRIVER)" >> $(BINARY_PATH)/.goose.env
	@echo "GOOSE_DBSTRING=$(MIGRATE_DB_USER):$(MIGRATE_DB_PASS)@tcp($(MIGRATE_DB_HOST):$(MIGRATE_DB_PORT))/$(MIGRATE_DB_NAME)" >> $(BINARY_PATH)/.goose.env


.PHONY: create-migration
create-migration:
	$(BINARY_PATH)/goose -dir $(MIGRATION_FILES_PATH) create ${name} sql

.PHONY: up-migration
up-migration:
	@$(eval include $(BINARY_PATH)/.goose.env)
	@$(eval export)
  	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" $(BINARY_PATH)/goose -dir $(MIGRATION_FILES_PATH) up

.PHONY: down-migration
down-migration:
	@$(eval include $(BINARY_PATH)/.goose.env)
	@$(eval export)
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" $(BINARY_PATH)/goose -dir $(MIGRATION_FILES_PATH) down-to 0

.PHONY: validate-migration
validate-migration:
	@$(eval include $(BINARY_PATH)/.goose.env)
	@$(eval export)
  	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" $(BINARY_PATH)/goose -dir $(MIGRATION_FILES_PATH) -v validate
