FROM golang:1.24 AS golang-builder

# in the image, make and navigate into a folder called /app
WORKDIR /app

# the commands to copy go.mod and go.sum and to download dependencies are run first, to leverage caching during the build of the image
# copy the go.mod and go.sum files into the directory /app
COPY go.mod go.sum ./

# download required dependencies
RUN go mod download 

#COPY all the rest of the files within the project directory, into the app directory of the golang-builder image
COPY . .

# now build the required binary, remember that within the project, to build this binary you would still need to navigate to the folder
# ./cmd/server_cmd to be able to access the main package for building out of the server
# the resulting binary that is built will be availble in /app/server_cmd in the golang-builder image

RUN go build -o server_cmd ./cmd/server_cmd 

# now start from another smaller image, make and navigate into a folder called /app
FROM debian:stable-slim

WORKDIR /app

# copy over the built binary from the golang-builder image, into the /app folder in this new debian based image
COPY --from=golang-builder /app/server_cmd .

# expost the port 8080
EXPOSE 8080

CMD ["./server_cmd"]