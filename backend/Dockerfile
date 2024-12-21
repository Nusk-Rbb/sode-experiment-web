# ベースイメージ
FROM golang:1.23-alpine AS builder

# 作業ディレクトリ
WORKDIR /app

# プロジェクトの依存関係を取得
COPY go.mod go.sum ./
RUN go mod tidy

# ソースコードをコピー
COPY . .

# ビルド
RUN go build -o main .


# 最終イメージ
FROM alpine:latest

# 作業ディレクトリ
WORKDIR /app

# ビルドされた実行ファイルをコピー
COPY --from=builder /app/main .

# 環境変数を設定
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASS=postgres
ENV DB_NAME=location_app

# ポートの設定
EXPOSE 8080

# 実行コマンド
CMD ["./main"]