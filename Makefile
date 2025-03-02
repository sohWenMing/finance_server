.PHONY: docker-server-build
# PHONY being used to ensure that docker does not associate the command docker-build with any files, will always run the build

# command builds out out the server_cmd binary, because it is referring to the Dockerfile defined in the path with the -f tag
# still needs to be run from the root of the project - note that the build context is defined by a path or url, so in this case should be .

# Usage:  docker buildx build [OPTIONS] PATH | URL | -

# go.sum still needs to be created in case it is not yet created, this is for the build of the container
# in the event that go.sum already exists, then nothing will happen regarding the go.sum file
docker-server-build:
	touch go.sum
	docker build -t finance_server -f cmd/server_cmd/Dockerfile .

.PHONY: docker-server-start
docker-server-start:
	docker run -d -p 8080:8080 --name finance_server finance_server

.PHONY: docker-server-stop
docker-server-stop:
	docker stop finance_server
	docker rm finance_server