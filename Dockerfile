# Start from a golang:1.6.1
# and a workspace (GOPATH) configured at /go.
FROM golang:1.6.1

# Create the directory for proper Go'ness
RUN mkdir -p $GOPATH/src/github.com/emccode

# Change the working directory
WORKDIR $GOPATH/src/github.com/emccode

# Clone the latest
RUN git clone https://github.com/emccode/polly

# Change the Working directory
WORKDIR $GOPATH/src/github.com/emccode/polly

# Build the go binary
RUN make

# Copy the binary to /usr/local and set permissions
RUN cp /go/bin/polly /usr/local/ && chmod u+x /usr/local/polly

# Start the Polly Service
ENTRYPOINT polly service start -f

EXPOSE 7978 7979
