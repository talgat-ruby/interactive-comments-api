FROM golang:1.21-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED 1
RUN apk --update add alpine-sdk
COPY . .
RUN go build -o api .

FROM alpine AS runner

# Add Maintainer Info
LABEL maintainer="Talgat Saribayev <talgat.s@protonmail.com>"

WORKDIR /app
COPY --from=builder /app/api api
COPY --from=builder /app/docs docs
COPY --from=builder /app/database.db database.db

ENV PORT 80
EXPOSE $PORT
CMD ["/app/api"]
