FROM golang:1.11 as build

ENV GO111MODULE=on
ARG repo

WORKDIR ${repo}/app
ADD ./app .

RUN CGO_ENABLE=0 GOOS=linux go build -o /tmp/app

FROM gcr.io/distroless/base
ADD ./runtime /runtime
COPY --from=build /tmp/app /
ENV ROOT_DIR=/runtime
ENTRYPOINT ["/app"]
