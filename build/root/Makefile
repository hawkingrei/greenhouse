# Don't allow an implicit 'all' rule. This is not a user-facing file.
ifeq ($(MAKECMDGOALS),)
    $(error This Makefile requires an explicit rule to be specified)
endif

ifeq ($(DBG_MAKEFILE),1)
    $(warning ***** starting Makefile.generated_files for goal(s) "$(MAKECMDGOALS)")
    $(warning ***** $(shell date))
endif

# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /bin/bash
ARCH      := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"
# We don't need make's built-in rules.
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

# Constants used throughout.
.EXPORT_ALL_VARIABLES:
OUT_DIR ?= _output
BIN_DIR := $(OUT_DIR)/bin

.PHONY: build update clean
all: check bazel-update 

build: init check bazel-build 

simple-build:
	bazel build --watchfs -- //...
bazel-build:
	bazel build --watchfs //...
clean:
	bazel clean --expunge 
update: init bazel-update

bazel-update:
	bazel run //:gazelle
check:
	@./build/check.sh
init:
	@if [ ! -f .git/hooks/pre-commit ] ; \
	then \
	echo "make all" >> .git/hooks/pre-commit; \
	sudo chmod +x .git/hooks/pre-commit; \
	fi
