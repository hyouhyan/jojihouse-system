# postgresql
import psycopg2
import dotenv
import os

dotenv.load_dotenv()


POSTGRES_USER = os.getenv("POSTGRES_USER")
POSTGRES_PASSWORD = os.getenv("POSTGRES_PASSWORD")
POSTGRES_DB = os.getenv("POSTGRES_DB")
POSTGRES_HOST = os.getenv("POSTGRES_HOST")
POSTGRES_PORT = os.getenv("POSTGRES_PORT")

conn = psycopg2.connect(
    dbname=POSTGRES_DB,
    user=POSTGRES_USER,
    password=POSTGRES_PASSWORD,
    host=POSTGRES_HOST,
    port=POSTGRES_PORT
)

cursor = conn.cursor()

# ハウスへの入室処理
def userEntry(user_id):
    cursor.execute("SELECT name FROM users WHERE user_id = %s", (user_id,))
    user = cursor.fetchone()
    if user is None:
        print("ユーザーが存在しません")
        return False
    else:
        print(f"{user[0]}さんが入室しました")

    # 入場可能回数(remaining_entries)を減らす
    # まず取得
    cursor.execute("SELECT remaining_entries FROM users WHERE user_id = %s", (user_id,))
    remaining_entries = cursor.fetchone()[0]
    # 減らす
    remaining_entries -= 1
    # 更新
    cursor.execute("UPDATE users SET remaining_entries = %s WHERE user_id = %s", (remaining_entries, user_id))
    conn.commit()

    return True

def userExit(user_id):
    cursor.execute("SELECT name FROM users WHERE user_id = %s", (user_id,))
    user = cursor.fetchone()
    if user is None:
        print("ユーザーが存在しません")
        return False
    else:
        print(f"{user[0]}さんが退室しました")

    return True