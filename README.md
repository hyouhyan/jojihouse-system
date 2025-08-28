# JOJIハウス管理システム(jojihouse-system)

本システムは姫雛ガジェット研究所の実作業場である**JOJIハウス**の諸々を管理するためのシステムです。

# API

## API エンドポイント一覧

### ユーザー関連
| メソッド | エンドポイント | 説明 | リクエスト | レスポンス | 権限 |
|----------|--------------|------|-----|-----|-----|
| `POST`   | `/users/`  | ユーザーを新規作成 | CreateUser | | SystemAdmin |
| `GET`    | `/users/`  | すべてのユーザーを取得 | | | Member |
| `GET`    | `/users/:user_id` | 指定したユーザーの情報を取得 | | | Member |
| `PATCH`  | `/users/:user_id` | 指定したユーザーの情報を部分更新 | UpdateUser | | HouseAdmin |
| `DELETE` | `/users/:user_id` | 指定したユーザーを削除 | | | SystemAdmin |
| `GET`    | `/users/:user_id/logs` | 該当ユーザに関連するログを取得 | | | HouseAdmin |

### ロール管理
| メソッド | エンドポイント | 説明 | リクエスト | レスポンス | 権限 |
|----------|--------------|------|-----|-----|-----|
| `GET`    | `/roles` | 全てのロールを取得 | | | Member |
| `GET`    | `/users/:user_id/roles` | 指定したユーザーのロールを取得 | | | Member |
| `POST`   | `/users/:user_id/roles` | 指定したユーザーにロールを追加 | AddRole | | SystemAdmin |
| `DELETE` | `/users/:user_id/roles/:role_id` | 指定したユーザーからロールを削除 | | | SystemAdmin |

---

### 入退室管理
| メソッド | エンドポイント | 説明 | リクエスト | レスポンス | 権限 |
|----------|--------------|------|-----|-----|-----|
| `POST`   | `/entrance/` | 入退室を記録 | Entrance | | ??? |
| `GET`    | `/entrance/current` | 在室ユーザ一覧を取得 | | | Member |
| `GET`    | `/entrance/logs` | すべての入退室ログを取得 | | | HouseAdmin |
| `GET`    | `/entrance/logs/:user_id` | 指定したユーザーの入退室ログを取得 | | | HouseAdmin |

---

### 支払い管理

| メソッド | エンドポイント | 説明 | リクエスト | レスポンス | 権限 |
|----------|--------------|------|-----|-----|-----|
| `GET`   | `/payment/` | すべての支払いログを取得 | | | HouseAdmin |
| `POST`   | `/payment/` | 支払いを記録 | Payment | | HouseAdmin |
| `GET`   | `/payment/monthly?year=:year&month=:month` | 指定月の支払いログを取得 | | | HouseAdmin |
| `POST`   | `/kaisuken/` | 入場料支払いを記録 | BuyKaisuken | | HouseAdmin |

---


# データベース

## PostgreSQL

### ユーザー(users)

`users` テーブルでは、オフィスの利用メンバーの情報を管理します。

| フィールド名 | データタイプ | 説明 | 備考 |
|-------------|------|------|------|
| id | SERIAL | 内部処理用 | primary key |
| name | VARCHAR(255) | 名前 | NOT NULL |
| description | TEXT | 概要（任意） |  |
| barcode | VARCHAR(64) | カードに印刷するバーコード（EAN-8） | 8桁の数字 NOT NULL |
| discord_id | VARCHAR(64) | DiscordアカウントのID |  |
| remaining_entries | INT | 入場可能回数 | 入場ごとに減少 |
| registered_at | TIMESTAMP WITH TIME ZONE | 登録日 |  |
| total_entries | INT | 総入場回数 |  |
| allergy | VARCHAT(255) | アレルギー | |
| number | INT | 会員番号 | UNIQUE |

### ロール(roles)

| フィールド名 | データタイプ | 説明 | 備考 |
|-------------|------|------|------|
| id | SERIAL | ロールid | |
| name | VARCHAR(255) | ロール名 | |

初期データ  
|name|説明|
|---|---|
|member|一般Labメンバー|
|student|学生|
|system-admin|システム管理者|
|house-admin|家賃負担者|
|guest|ゲスト(ラボメン以外, 使うことある？)|

### ロール用中間テーブル(user_roles)

| フィールド名 | 説明 | 備考 |
|-------------|------|------|
|user_id|メンバーID| REFERENCES users(id) |
|role_id|ロールID| REFERENCES roles(id) |
|(user_id, role_id)| | Primary Key | 

### 在室ユーザー(current_users)

| フィールド名 | データタイプ | 説明 | 備考 |
|-------------|------|------|------|
| user_id | INT | ユーザーID | REFERENCES users(id) ON DELETE CASCADE |
| entered_at | TIMESTAMP WITH TIME ZONE | 入場時間 | NOT NULL |


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

### 支払いログ(payment_log)

| フィールド名 | 説明 | 備考 |
|-------------|------|------|
| log_id | 内部処理用 | primary key |
| user_id | users テーブルと紐づけ |  |
| time | 支払い日時 |  |
| description | 説明 |  |
| amount | 金額 |  |
| payway | 支払い方法 | "cash" or "olive" |
