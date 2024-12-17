FROM public.ecr.aws/docker/library/golang:1.23 AS builder

WORKDIR /app
COPY . .

RUN make install

FROM public.ecr.aws/docker/library/alpine:3.16
RUN apk update && \
    apk add ca-certificates libc6-compat curl && \
    rm -rf /var/cache/apk/*
WORKDIR /app

COPY --from=builder /go/bin/validation /app/service
ENTRYPOINT ["/app/service"]
