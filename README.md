# JOJIハウス 入退出管理システム

# データベース

## PostgreSQL

### メンバー(members)

`members` テーブルでは、オフィスの利用メンバーの情報を管理します。

| フィールド名 | 説明 | 備考 |
|-------------|------|------|
| id | 内部処理用 | primary key |
| name | ニックネーム |  |
| description | 概要（任意） |  |
| barcode | カードに印刷するバーコード（EAN-13） | 13桁の数字 |
| contact | 連絡先情報（Xアカウントなど） |  |
| remaining_entries | 入場可能回数 | 入場ごとに減少（同日再入場では減らない） |
| registered_at | 登録日 |  |
| entry_count | 総入場回数 |  |

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

### ロール用中間テーブル(member_roles)

| フィールド名 | 説明 |
|-------------|------|
|member_id|メンバーID|
|role_id|ロールID|
|(member_id, role_id)|Primary Key|

## MongoDB

### 入退室ログ(access_log)

`access_log` コレクションでは、メンバーの入退室の記録を管理します。

| フィールド名 | 説明 | 備考 |
|-------------|------|------|
| log_id | 内部処理用 | primary key |
| member_id | members テーブルと紐づけ |  |
| time | 入退出の時刻 |  |
| access_type | 入退室の種類 | "entry" or "exit" |

### 入場可能回数の変更ログ(entry_count_log)
`entry_count_log` コレクションでは、入場可能回数の変更履歴を管理します。

| フィールド名 | 説明 | 備考 |
|-------------|------|------|
| log_id | 内部処理用 | primary key |
| member_id | members テーブルと紐づけ |  |
| previous_entries | 変更前の入場可能回数 |  |
| new_entries | 変更後の入場可能回数 |  |
| reason | 追加理由 |  |
| updated_by | 変更を行った管理者名 |  |
| updated_at | 変更日時 |  |
