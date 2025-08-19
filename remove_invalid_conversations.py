import sqlite3
import json

DB_PATH = 'service/database/app.db'
MISSING_IDS = {'924b1428-58b0-47c5-b6d6-bdba93910c6e', '4c5c8cf1-a53f-4637-bb43-a1eab4f88e4e', 'dd34fba5-f367-4c50-afbb-27b6bd775b85', 'e7ea4543-9138-4508-bf02-e0d61370c274'}

def remove_invalid_conversations():
    conn = sqlite3.connect(DB_PATH)
    c = conn.cursor()
    c.execute('SELECT id, participants FROM conversations')
    to_delete = []
    for conv_id, participants_json in c.fetchall():
        try:
            participants = json.loads(participants_json)
        except Exception as e:
            print(f"Skipping {conv_id}: invalid JSON")
            continue
        for pid in participants:
            if pid in MISSING_IDS:
                to_delete.append(conv_id)
                print(f"Will delete conversation {conv_id} (references missing user {pid})")
                break
    for conv_id in to_delete:
        c.execute('DELETE FROM conversations WHERE id = ?', (conv_id,))
        print(f"Deleted conversation {conv_id}")
    conn.commit()
    conn.close()
    print("Done.")

if __name__ == '__main__':
    remove_invalid_conversations()
