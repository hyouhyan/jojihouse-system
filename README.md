# JOJIハウス管理システム(jojihouse-system)

本システムは姫雛ガジェット研究所の実作業場である**JOJIハウス**の諸々を管理するためのシステムです。

# API

## API エンドポイント一覧

### ユーザー関連
| メソッド | エンドポイント | 説明 | リクエスト | レスポンス |
|----------|--------------|------|-----|-----|
| `POST`   | `/users/`  | ユーザーを新規作成 |
| `GET`    | `/users/`  | すべてのユーザーを取得 |
| `GET`    | `/users/:user_id` | 指定したユーザーの情報を取得 |
| `PATCH`  | `/users/:user_id` | 指定したユーザーの情報を部分更新 |
| `DELETE` | `/users/:user_id` | 指定したユーザーを削除 |
| `GET`    | `/users/:user_id/logs` | 該当ユーザに関連するログを取得 |

### ロール管理
| メソッド | エンドポイント | 説明 | リクエスト | レスポンス |
|----------|--------------|------|-----|-----|
| `GET`    | `/roles` | 全てのロールを取得 |
| `GET`    | `/users/:user_id/roles` | 指定したユーザーのロールを取得 |
| `POST`   | `/users/:user_id/roles` | 指定したユーザーにロールを追加 |
| `DELETE` | `/users/:user_id/roles/:role_id` | 指定したユーザーからロールを削除 |

---

### 入退室管理
| メソッド | エンドポイント | 説明 | リクエスト | レスポンス |
|----------|--------------|------|-----|-----|
| `POST`   | `/entrance/` | 入退室を記録 |
| `GET`    | `/entrance/current` | 在室ユーザ一覧を取得 |
| `GET`    | `/entrance/logs` | すべての入退室ログを取得 |
| `GET`    | `/entrance/logs/:user_id` | 指定したユーザーの入退室ログを取得 |

---

### 支払い管理

| メソッド | エンドポイント | 説明 | リクエスト | レスポンス |
|----------|--------------|------|-----|-----|
| `GET`   | `/payment/` | すべての支払いログを取得 |
| `POST`   | `/payment/` | 支払いを記録 |
| `GET`   | `/payment/monthly?year=${year}&month=${month}` | 指定月の支払いログを取得 |
| `POST`   | `/kaisuken/` | 入場料支払いを記録 |

---


# データベース

## PostgreSQL

### ユーザー(users)

`users` テーブルでは、オフィスの利用メンバーの情報を管理します。

| フィールド名 | 説明 | 備考 |
|-------------|------|------|
| id | 内部処理用 | primary key |
| name | ニックネーム |  |
| description | 概要（任意） |  |
| barcode | カードに印刷するバーコード（EAN-13） | 13桁の数字 |
| contact | 連絡先情報（Xアカウントなど） |  |
| remaining_entries | 入場可能回数 | 入場ごとに減少（同日再入場では減らない） |
| registered_at | 登録日 |  |
| total_entries | 総入場回数 |  |

### ロール(roles)

| フィールド名 | 説明 | 備考 |
|-------------|------|------|
|id|役割のid||
|name|役割名||

初期データ  
|name|説明|
|---|---|
|member|一般Labメンバー|
|student|学生|
|system-admin|システム管理者|
|house-admin|ハウス管理者(月額出資者)|
|guest|ゲスト(ラボメン以外, 使うことある？)|

### ロール用中間テーブル(user_roles)

| フィールド名 | 説明 |
|-------------|------|
|user_id|メンバーID|
|role_id|ロールID|
|(user_id, role_id)|Primary Key|

### 在室ユーザー(current_users)

| フィールド名 | 説明 |
|-------------|------|
|user_id|ユーザーID|
|entered_at|入場時間|


## MongoDB

### 入退室ログ(access_log)

`access_log` コレクションでは、メンバーの入退室の記録を管理します。

| フィールド名 | 説明 | 備考 |
|-------------|------|------|
| log_id | 内部処理用 | primary key |
| user_id | users テーブルと紐づけ |  |
| time | 入退出の時刻 |  |
| access_type | 入退室の種類 | "entry" or "exit" |

### 入場可能回数ログ(remaining_entries_log)
`remaining_entries_log` コレクションでは、入場可能回数の変更履歴を管理します。

| フィールド名 | 説明 | 備考 |
|-------------|------|------|
| log_id | 内部処理用 | primary key |
| user_id | users テーブルと紐づけ |  |
| previous_entries | 変更前の入場可能回数 |  |
| new_entries | 変更後の入場可能回数 |  |
| reason | 追加理由 |  |
| updated_by | 変更を行った管理者名 |  |
| updated_at | 変更日時 |  |
