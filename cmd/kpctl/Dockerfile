FROM golang:1.20-alpine as Builder

ARG COMMAND
ARG VERSION

WORKDIR /build

ADD go.mod ./

RUN go mod download

COPY . .

WORKDIR /build/cmd/${COMMAND}

RUN go build -o /bin/main main.go

FROM scratch

COPY --from=Builder /bin/main main

ENTRYPOINT ["/main"]
CMD ["--help"]