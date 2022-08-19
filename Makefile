NAME := get-started
MOCKGEN_BUILD_FLAGS ?= -mod=mod

deploy:
	ibmcloud cf push

go.list:
	go list -json -m all > go.list

mockgen:
	mockgen -build_flags=$(MOCKGEN_BUILD_FLAGS) -destination=./services/mock_services/mock_dictionary.go github.com/hblab-ngocnd/get-started/services DictionaryService
	mockgen -build_flags=$(MOCKGEN_BUILD_FLAGS) -destination=./services/mock_services/mock_translate.go github.com/hblab-ngocnd/get-started/services TranslateService


