TAG = 20170404
VOLUME := $(shell pwd)

REPO = krausm/$(shell basename `pwd`)
IMAGE=$(REPO):$(TAG)

.PHONY: echo build run

echo:
	echo $(IMAGE)
build:
	docker build -t $(IMAGE) .
	docker images | grep '$(REPO)'
run:
	-docker run -it --rm -v $(VOLUME):/go/src/app $(IMAGE)
