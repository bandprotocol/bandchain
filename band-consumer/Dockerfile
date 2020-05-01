# # Simple usage with a mounted data directory:
# # > docker build -t gaia .
# # > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.bcd:/root/.gaiad -v ~/.gaiacli:/root/.gaiacli gaia gaiad init
# # > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.bcd:/root/.gaiad -v ~/.gaiacli:/root/.gaiacli gaia gaiad start
# FROM golang:1.14-alpine AS build-env

# # Set up dependencies
# ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python

# # Set working directory for the build
# WORKDIR /go/src/github.com/bandprotocol/band-consumer

# # Add source files
# COPY . ./band-consumer

# COPY ../chain chain

# # Install minimum necessary dependencies, build Cosmos SDK, remove packages
# RUN apk add --no-cache $PACKAGES && \
#     make tools && \
#     make install

# # Final image
# FROM alpine:edge

# # Install ca-certificates
# RUN apk add --update ca-certificates
# WORKDIR /root

# # Copy over binaries from the build-env
# COPY --from=build-env /go/bin/bcd /usr/bin/bcd
# COPY --from=build-env /go/bin/bccli /usr/bin/bccli

# # Run bcd by default, omit entrypoint to ease using container with bccli
# CMD ["bcd"]
