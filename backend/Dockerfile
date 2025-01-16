# Choose whatever you want, version >= 1.16
FROM golang:1.23-bullseye

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY . .
RUN go mod download

# 環境変数を設定
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASS=postgres
ENV DB_NAME=location_app

# ポートの設定
EXPOSE 8080

CMD ["air", "-c", ".air.toml"]