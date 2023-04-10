include ../../scripts/cmd.mk

COMMAND = kpeering
VERSION = v0.0.0

# temp use of docker to run kpeering
docker_run:
	docker run \
		--network kpeering-net \
		--mount type=bind,source=$(CURDIR)/.kpeering,target=/.kpeering \
		kpeering start