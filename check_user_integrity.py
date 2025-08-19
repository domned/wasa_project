import sqlite3
import json
from collections import Counter

DB_PATH = 'service/database/app.db'

def check_users_and_conversations():
    conn = sqlite3.connect(DB_PATH)
    c = conn.cursor()

    # Get all users
    c.execute('SELECT id, username FROM users')
    users = c.fetchall()
    user_ids = set()
    username_counter = Counter()
    for uid, uname in users:
        user_ids.add(uid)
        username_counter[uname] += 1

    # Print duplicate usernames
    print('Duplicate usernames:')
    for uname, count in username_counter.items():
        if count > 1:
            print(f'  {uname} (count: {count})')

    # Get all user IDs referenced in conversations
    c.execute('SELECT id, participants FROM conversations')
    missing = set()
    for conv_id, participants_json in c.fetchall():
        try:
            participants = json.loads(participants_json)
        except Exception as e:
            print(f'  Conversation {conv_id} has invalid JSON')
            continue
        for pid in participants:
            if pid not in user_ids:
                print(f'  Conversation {conv_id} references missing user id: {pid}')
                missing.add(pid)

    if not missing:
        print('All conversation participants exist in users table.')
    else:
        print(f'Missing user IDs: {missing}')

    conn.close()

if __name__ == '__main__':
    check_users_and_conversations()
