# ベースイメージ
FROM node:18-alpine as builder

# 作業ディレクトリ
WORKDIR /app

# package.json, package-lock.jsonをコピー
COPY package*.json ./
# 依存関係をインストール
RUN npm install

# ソースコードをコピー
COPY . .
# ビルド (ビルドステップは不要なので省略)
# 本番環境ではビルド処理を入れてください


# 最終イメージ
FROM nginx:alpine

# アプリケーションの静的ファイルをコピー
COPY --from=builder /app ./app

# Nginxの設定ファイルを作成
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d
# ポートの設定
EXPOSE 80
# 実行コマンド
CMD ["nginx", "-g", "daemon off;"]