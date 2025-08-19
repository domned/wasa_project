import sqlite3
import json

DB_PATH = 'service/database/app.db'

def fix_participants():
    conn = sqlite3.connect(DB_PATH)
    c = conn.cursor()
    c.execute('SELECT id, participants FROM conversations')
    rows = c.fetchall()
    for conv_id, participants_json in rows:
        try:
            participants = json.loads(participants_json)
        except Exception as e:
            print(f"Skipping {conv_id}: invalid JSON")
            continue
        # If already a list of strings, skip
        if participants and isinstance(participants[0], str):
            continue
        # If list of dicts, convert to list of ids
        if participants and isinstance(participants[0], dict) and 'id' in participants[0]:
            new_participants = [p['id'] for p in participants]
            new_json = json.dumps(new_participants)
            c.execute('UPDATE conversations SET participants = ? WHERE id = ?', (new_json, conv_id))
            print(f"Fixed participants for conversation {conv_id}")
    conn.commit()
    conn.close()
    print("All conversations checked.")

if __name__ == '__main__':
    fix_participants()
