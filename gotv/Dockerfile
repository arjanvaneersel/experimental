# Use the default Debian based golang image
FROM golang
MAINTAINER arjan@balkan.tech

# Copy all files to the container
ADD . /go/src/ponzu-deployment

# Install ponzu
RUN go get github.com/ponzu-cms/ponzu/...

# Build the image
RUN cd /go/src/ponzu-deployment
RUN ponzu build

# Run Ponzu on port 80 when the container starts
ENTRYPOINT /go/bin/ponzu run
EXPOSE 80
