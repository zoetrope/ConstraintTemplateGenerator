XDG_CONFIG_HOME = ~/.config
API_VERSION = kustomize.cybozu.com/v1
PLUGIN_KIND = constrainttemplategenerator
BIN = ConstraintTemplateGenerator
SRCS := $(shell find . -type f -name '*.go')

$(BIN): $(SRCS)
	go build ./

install: $(BIN)
	mkdir -p ${XDG_CONFIG_HOME}/kustomize/plugin/${API_VERSION}/${PLUGIN_KIND}/
	cp $(BIN) ${XDG_CONFIG_HOME}/kustomize/plugin/${API_VERSION}/${PLUGIN_KIND}/

.PHONY: install
