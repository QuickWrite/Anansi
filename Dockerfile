# Build step
FROM golang:1.26-alpine as build
WORKDIR /
COPY . ./
RUN go build

# Execution step
FROM busybox

LABEL org.opencontainers.image.title="Anansi" \
      org.opencontainers.image.description="A mischievous web application to generate an endless stream of senseless text." \
      org.opencontainers.image.licenses="MPL-2.0" \
      org.opencontainers.image.url="https://github.com/QuickWrite/Anansi" \
      org.opencontainers.image.source="https://github.com/QuickWrite/Anansi"

COPY --from=build /Anansi /

# The current application runs on the port 8080
EXPOSE 8080/tcp

CMD [ "/Anansi" ]
