# 袖実験用のWebアプリ

## Webアプリ概要

esp32から検知情報を受信しユーザーが外にいればユーザーに警報を出すアプリ

## やること

- [x] 環境構築
  - [x] Nginx + Go + PostgreSQLのDocker構築
- [ ] コード
  - [x] esp32受信用したとき
    - [x] esp32用のPOST処理
  - [x] 位置情報取得
  - [x] 家とユーザーの距離計算
  - [ ] ユーザー警告(LINE bot or Mail)