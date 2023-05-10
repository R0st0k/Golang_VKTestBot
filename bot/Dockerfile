# syntax=docker/dockerfile:1

FROM golang:1.19-buster AS builder

ENV USER=appuser
ENV UID=1001

# https://stackoverflow.com/a/55757473
RUN adduser \
    --disabled-password \
    --gecos "" \
    --no-create-home \
    --uid "$UID" \
    "$USER"

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

WORKDIR /usr/src/app/cmd/vk-bot-test-task

RUN CGO_ENABLED=0 go build -v -o /usr/local/bin/app

FROM gcr.io/distroless/static-debian11:nonroot

WORKDIR /

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /usr/local/bin/app /app

# Use an unprivileged user.
USER $USER:$USER

CMD ["/app"]