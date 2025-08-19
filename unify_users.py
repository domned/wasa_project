import sqlite3
import json
from collections import defaultdict

DB_PATH = 'service/database/app.db'

def unify_users():
    conn = sqlite3.connect(DB_PATH)
    c = conn.cursor()

    # 1. Find all users grouped by username
    c.execute('SELECT id, username FROM users')
    users = c.fetchall()
    username_to_ids = defaultdict(list)
    for user_id, username in users:
        username_to_ids[username].append(user_id)

    # 2. For each username with duplicates, pick canonical and remap others
    for username, ids in username_to_ids.items():
        if len(ids) <= 1:
            continue
        canonical_id = ids[0]
        duplicate_ids = ids[1:]
        print(f"Unifying {username}: keeping {canonical_id}, removing {duplicate_ids}")

        # Update messages.sender_id
        c.execute(f"""
            UPDATE messages SET sender_id = ? WHERE sender_id IN ({','.join(['?']*len(duplicate_ids))})
        """, [canonical_id] + duplicate_ids)

        # Update contacts.user_id and contacts.contact_id
        c.execute(f"""
            UPDATE contacts SET user_id = ? WHERE user_id IN ({','.join(['?']*len(duplicate_ids))})
        """, [canonical_id] + duplicate_ids)
        c.execute(f"""
            UPDATE contacts SET contact_id = ? WHERE contact_id IN ({','.join(['?']*len(duplicate_ids))})
        """, [canonical_id] + duplicate_ids)

        # Update conversations.participants (JSON array of user IDs)
        c.execute('SELECT id, participants FROM conversations')
        for conv_id, participants_json in c.fetchall():
            participants = json.loads(participants_json)
            changed = False
            for i, pid in enumerate(participants):
                if pid in duplicate_ids:
                    participants[i] = canonical_id
                    changed = True
            if changed:
                new_json = json.dumps(participants)
                c.execute('UPDATE conversations SET participants = ? WHERE id = ?', (new_json, conv_id))

        # Remove duplicate users
        c.execute(f"DELETE FROM users WHERE id IN ({','.join(['?']*len(duplicate_ids))})", duplicate_ids)

    conn.commit()
    conn.close()
    print("User unification complete.")

if __name__ == '__main__':
    unify_users()
