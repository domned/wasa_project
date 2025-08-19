import sqlite3
import json
from collections import Counter

DB_PATH = 'service/database/app.db'

def check_all_refs():
    conn = sqlite3.connect(DB_PATH)
    c = conn.cursor()
    c.execute('SELECT id FROM users')
    user_ids = set(row[0] for row in c.fetchall())
    missing = set()
    # Check conversations
    c.execute('SELECT id, participants FROM conversations')
    for conv_id, participants_json in c.fetchall():
        try:
            participants = json.loads(participants_json)
        except Exception:
            print(f'Conversation {conv_id} has invalid JSON')
            continue
        for pid in participants:
            if pid not in user_ids:
                print(f'Conversation {conv_id} references missing user id: {pid}')
                missing.add(pid)
    # Check messages
    c.execute('SELECT id, sender_id FROM messages')
    for msg_id, sender_id in c.fetchall():
        if sender_id not in user_ids:
            print(f'Message {msg_id} references missing sender_id: {sender_id}')
            missing.add(sender_id)
    # Check contacts
    c.execute('SELECT id, user_id, contact_id FROM contacts')
    for cid, user_id, contact_id in c.fetchall():
        if user_id not in user_ids:
            print(f'Contact {cid} references missing user_id: {user_id}')
            missing.add(user_id)
        if contact_id not in user_ids:
            print(f'Contact {cid} references missing contact_id: {contact_id}')
            missing.add(contact_id)
    # Check reactions
    c.execute('SELECT id, sender_id FROM reactions')
    for rid, sender_id in c.fetchall():
        if sender_id not in user_ids:
            print(f'Reaction {rid} references missing sender_id: {sender_id}')
            missing.add(sender_id)
    # Check comments
    c.execute('SELECT id, sender_id FROM comments')
    for cmid, sender_id in c.fetchall():
        if sender_id not in user_ids:
            print(f'Comment {cmid} references missing sender_id: {sender_id}')
            missing.add(sender_id)
    if not missing:
        print('All user references are valid.')
    else:
        print(f'Missing user IDs referenced: {missing}')
    conn.close()

if __name__ == '__main__':
    check_all_refs()
