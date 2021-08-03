ARG DATABASE_FILE=/test.db

FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /gatekeeper ./cmd/server.go

ARG DATABASE_FILE

RUN touch $DATABASE_FILE && chmod 755 $DATABASE_FILE && ls -al /

FROM gcr.io/distroless/base-debian10
#FROM golang:1.16-buster

WORKDIR /

ARG DATABASE_FILE
COPY --from=build /gatekeeper /gatekeeper

COPY --from=build $DATABASE_FILE $DATABASE_FILE

EXPOSE 8080


ENTRYPOINT ["/gatekeeper"]
CMD ["--jwt-secret", "$GATEKEEPER_SHARED_SECRET", "--sqllite", "$DATABASE_FILE"]
