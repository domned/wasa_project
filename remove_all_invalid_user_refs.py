import sqlite3
import json

DB_PATH = 'service/database/app.db'
MISSING_IDS = {'9218e5cc-0b71-48b7-b70e-3566b2205390', '4c5c8cf1-a53f-4637-bb43-a1eab4f88e4e', 'dd34fba5-f367-4c50-afbb-27b6bd775b85', 'e7ea4543-9138-4508-bf02-e0d61370c274', 'df05e551-425e-4dde-b04f-528ef005eefc', '924b1428-58b0-47c5-b6d6-bdba93910c6e', '69b38843-dbb9-47a6-9529-e28419792c2f', '25d5545e-c9b3-4cc5-8d51-4deb58a39cc3', 'db2299d8-90a5-41fd-80ed-858c8ffc8417'}

def remove_invalid_refs():
    conn = sqlite3.connect(DB_PATH)
    c = conn.cursor()
    # Conversations
    c.execute('SELECT id, participants FROM conversations')
    for conv_id, participants_json in c.fetchall():
        try:
            participants = json.loads(participants_json)
        except Exception:
            continue
        if any(pid in MISSING_IDS for pid in participants):
            c.execute('DELETE FROM conversations WHERE id = ?', (conv_id,))
            print(f'Deleted conversation {conv_id}')
    # Messages
    c.execute('SELECT id, sender_id FROM messages')
    for msg_id, sender_id in c.fetchall():
        if sender_id in MISSING_IDS:
            c.execute('DELETE FROM messages WHERE id = ?', (msg_id,))
            print(f'Deleted message {msg_id}')
    # Contacts
    c.execute('SELECT id, user_id, contact_id FROM contacts')
    for cid, user_id, contact_id in c.fetchall():
        if user_id in MISSING_IDS or contact_id in MISSING_IDS:
            c.execute('DELETE FROM contacts WHERE id = ?', (cid,))
            print(f'Deleted contact {cid}')
    # Reactions
    c.execute('SELECT id, sender_id FROM reactions')
    for rid, sender_id in c.fetchall():
        if sender_id in MISSING_IDS:
            c.execute('DELETE FROM reactions WHERE id = ?', (rid,))
            print(f'Deleted reaction {rid}')
    # Comments
    c.execute('SELECT id, sender_id FROM comments')
    for cmid, sender_id in c.fetchall():
        if sender_id in MISSING_IDS:
            c.execute('DELETE FROM comments WHERE id = ?', (cmid,))
            print(f'Deleted comment {cmid}')
    conn.commit()
    conn.close()
    print('Done cleaning all invalid user references.')

if __name__ == '__main__':
    remove_invalid_refs()
