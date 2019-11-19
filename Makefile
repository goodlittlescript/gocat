PACKAGE=goodlittlescript.com/gocat
COMMAND=$(notdir $(PACKAGE))
IMAGE_NAME=$(COMMAND)
TARGETS=shell app

#
# develop
#

images:
	@for target in $(TARGETS); \
	do docker build --build-arg PACKAGE=$(PACKAGE) --target $$target -t $(IMAGE_NAME):$$target .; \
	done

run:
	docker run --rm -i $(IMAGE_NAME):app /app/bin/$(COMMAND)

test:
	docker run --rm -v $(PWD):/go/src/$(PACKAGE) $(IMAGE_NAME):shell ./test/suite

lint:
	docker run --rm -v $(PWD):/go/src/$(PACKAGE) $(IMAGE_NAME):shell test -z "$$(gofmt -s -e -d . | tee /dev/stderr)"

fix:
	docker run --rm -v $(PWD):/go/src/$(PACKAGE) $(IMAGE_NAME):shell gofmt -s -w .

shell:
	docker run -it --rm -v $(PWD):/go/src/$(PACKAGE) $(IMAGE_NAME):shell /bin/bash

#
# deploy
#

artifacts:
	make images test lint >&2 && echo $(IMAGE_NAME):app

#
# utilities
#

phony:
	@sed -ne 's/^\([[:alnum:]_-]\{1,\}\):.*/	\1 \\/p' Makefile | sed -e '$$s/ \\//'

.PHONY: \
	images \
	run \
	test \
	lint \
	fix \
	shell \
	artifacts \
	phony
