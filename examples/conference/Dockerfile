FROM golang:1.22-alpine3.18 as build-env
RUN apk --no-cache add git 2>&1 >/dev/null
RUN apk --no-cache add ca-certificates && \
    update-ca-certificates
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /build/conference

FROM busybox:1.35.0
ENV APP_ENTRYPOINT=/app/conference-go
RUN adduser -D entando
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env --chown=entando:entando /build/conference-go $APP_ENTRYPOINT
USER entando
ENTRYPOINT ["sh", "-c", "${APP_ENTRYPOINT}"]
EXPOSE 8081