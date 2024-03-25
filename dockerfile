FROM golang:1.21-bookworm AS build

ARG GITHUB_USERNAME
ARG GITHUB_TOKEN
ARG GITHUB_REPO="github.com/elmodis/go-libs"

RUN git config --global \
    url."https://${GITHUB_USERNAME}:${GITHUB_TOKEN}@${GITHUB_REPO}".insteadOf "https://${GITHUB_REPO}"

RUN go install github.com/swaggo/swag/cmd/swag@latest
WORKDIR /app
COPY . .

WORKDIR /app/src
RUN swag i

RUN go mod download
RUN go build .

FROM debian:bookworm-slim AS run

WORKDIR /app
COPY --from=build /app/src/ .
CMD ["./scripts-api"]

