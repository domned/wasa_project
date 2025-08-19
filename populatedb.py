import sqlite3
import json
import uuid
import os

DB_PATH = "service/database/app.db"  # Change if your DB is elsewhere
print("Using DB at:", os.path.abspath(DB_PATH))

def make_id():
    return str(uuid.uuid4())

def main():
    conn = sqlite3.connect(DB_PATH)
    c = conn.cursor()
    # CLEAR TABLES for a clean slate (after cursor is created)
    c.execute("DELETE FROM comments")
    c.execute("DELETE FROM reactions")
    c.execute("DELETE FROM messages")
    c.execute("DELETE FROM conversations")
    c.execute("DELETE FROM contacts")
    c.execute("DELETE FROM users")

    # USERS
    users = [
        {"id": "f2555a8a-2e66-4326-9588-20e7e298d615", "username": "Alice", "picture": None},
        {"id": "b1b2b3b4-1111-2222-3333-444455556666", "username": "Bob", "picture": None},
        {"id": "c1c2c3c4-1111-2222-3333-444455556666", "username": "Charlie", "picture": None},
        {"id": "d1d2d3d4-1111-2222-3333-444455556666", "username": "Diana", "picture": None},
        {"id": "e1e2e3e4-1111-2222-3333-444455556666", "username": "Eve", "picture": None},
        {"id": "f1f2f3f4-1111-2222-3333-444455556666", "username": "Frank", "picture": None},
        {"id": "g1g2g3g4-1111-2222-3333-444455556666", "username": "Grace", "picture": None},
        {"id": "h1h2h3h4-1111-2222-3333-444455556666", "username": "Heidi", "picture": None},
    ]
    for user in users:
        pic = user["picture"] if user["picture"] is not None else ""
        c.execute("INSERT OR IGNORE INTO users (id, username, picture) VALUES (?, ?, ?)", (user["id"], user["username"], pic))

    # CONTACTS (each user adds the next as a contact, circular)
    contacts = []
    for i, user in enumerate(users):
        contact = users[(i+1)%len(users)]
        contacts.append({
            "id": make_id(),
            "user_id": user["id"],
            "contact_id": contact["id"]
        })
    for contact in contacts:
        # If you ever add a picture field to contacts, ensure it's not None
        c.execute("INSERT OR IGNORE INTO contacts (id, user_id, contact_id) VALUES (?, ?, ?)", (contact["id"], contact["user_id"], contact["contact_id"]))

    # CONVERSATIONS (1:1 and group)
    conversations = [
        # 1:1 conversations (name=None)
        {"participants": [users[0], users[1]], "name": None, "picture": None},
        {"participants": [users[0], users[2]], "name": None, "picture": None},
        {"participants": [users[0], users[3]], "name": None, "picture": None},
        {"participants": [users[0], users[4]], "name": None, "picture": None},
        {"participants": [users[0], users[5]], "name": None, "picture": None},
        {"participants": [users[0], users[6]], "name": None, "picture": None},
        {"participants": [users[0], users[7]], "name": None, "picture": None},
        # Group conversations (name set)
        {"participants": [users[0], users[2], users[3]], "name": "Project Team", "picture": None},
        {"participants": [users[0], users[2], users[4]], "name": "Fun Group", "picture": None},
        {"participants": [users[0], users[1], users[4]], "name": "Secret Club", "picture": None},
        {"participants": [users[0], users[5], users[6]], "name": "Alice, Frank & Grace", "picture": None},
        {"participants": [users[0], users[5], users[7]], "name": "Alice, Frank & Heidi", "picture": None},
        {"participants": [users[0], users[1], users[5], users[6]], "name": "Big Group", "picture": None}
    ]
    # Use fixed UUIDs for conversations for reproducibility
    fixed_conv_ids = [
        "c0c0c0c0-0000-0000-0000-000000000001",
        "c0c0c0c0-0000-0000-0000-000000000002",
        "c0c0c0c0-0000-0000-0000-000000000003",
        "c0c0c0c0-0000-0000-0000-000000000004",
        "c0c0c0c0-0000-0000-0000-000000000005",
        "c0c0c0c0-0000-0000-0000-000000000006",
        "c0c0c0c0-0000-0000-0000-000000000007",
        "c0c0c0c0-0000-0000-0000-000000000008",
        "c0c0c0c0-0000-0000-0000-000000000009",
        "c0c0c0c0-0000-0000-0000-00000000000a",
        "c0c0c0c0-0000-0000-0000-00000000000b",
        "c0c0c0c0-0000-0000-0000-00000000000c",
        "c0c0c0c0-0000-0000-0000-00000000000d"
    ]
    conv_ids = []
    for idx, conv in enumerate(conversations):
        conv_id = fixed_conv_ids[idx]
        conv_ids.append(conv_id)
        # Store only user IDs in participants JSON
        participants_json = json.dumps([u["id"] for u in conv["participants"]])
        pic = conv["picture"] if conv["picture"] is not None else ""
        c.execute("INSERT INTO conversations (id, participants, name, picture) VALUES (?, ?, ?, ?)", (conv_id, participants_json, conv["name"], pic))

    # Print conversation index, participants, and ID mapping for debugging
    print("Conversation index mapping:")
    for idx, conv in enumerate(conversations):
        participant_names = ', '.join([u['username'] for u in conv['participants']])
        print(f"  idx={idx}: participants=[{participant_names}] id={conv_ids[idx]}")

    # MESSAGES
    messages = [
        # Existing conversations
        {"conversation_idx": 0, "sender": users[0], "message": "Hey Bob!"},
        {"conversation_idx": 0, "sender": users[1], "message": "Hi Alice!"},
        {"conversation_idx": 1, "sender": users[2], "message": "Hello team!"},
        {"conversation_idx": 1, "sender": users[3], "message": "Let's get started."},
        {"conversation_idx": 2, "sender": users[0], "message": "Hi Charlie!"},
        {"conversation_idx": 2, "sender": users[2], "message": "Hey Alice!"},
        {"conversation_idx": 3, "sender": users[0], "message": "Hi Diana!"},
        {"conversation_idx": 3, "sender": users[3], "message": "Hello Alice!"},
        {"conversation_idx": 4, "sender": users[0], "message": "Hi Eve!"},
        {"conversation_idx": 4, "sender": users[4], "message": "Hey Alice!"},
        {"conversation_idx": 5, "sender": users[0], "message": "Welcome to Fun Group!"},
        {"conversation_idx": 5, "sender": users[2], "message": "Thanks Alice!"},
        {"conversation_idx": 5, "sender": users[4], "message": "Glad to be here!"},
        {"conversation_idx": 6, "sender": users[0], "message": "Welcome to Secret Club!"},
        {"conversation_idx": 6, "sender": users[1], "message": "Shh! It's a secret."},
        {"conversation_idx": 6, "sender": users[4], "message": "Excited!"},
        # More 1:1 and group chats for Alice
        {"conversation_idx": 7, "sender": users[0], "message": "Hi Frank!"},
        {"conversation_idx": 7, "sender": users[5], "message": "Hey Alice!"},
        {"conversation_idx": 8, "sender": users[0], "message": "Hi Grace!"},
        {"conversation_idx": 8, "sender": users[6], "message": "Hello Alice!"},
        {"conversation_idx": 9, "sender": users[0], "message": "Hi Heidi!"},
        {"conversation_idx": 9, "sender": users[7], "message": "Hi Alice!"},
        {"conversation_idx": 10, "sender": users[0], "message": "Hey Frank and Grace!"},
        {"conversation_idx": 10, "sender": users[5], "message": "Hey all!"},
        {"conversation_idx": 10, "sender": users[6], "message": "Hi Alice!"},
        {"conversation_idx": 11, "sender": users[0], "message": "Hey Frank and Heidi!"},
        {"conversation_idx": 11, "sender": users[5], "message": "Hi Alice!"},
        {"conversation_idx": 11, "sender": users[7], "message": "Hi all!"},
        {"conversation_idx": 12, "sender": users[0], "message": "Welcome to the Big Group!"},
        {"conversation_idx": 12, "sender": users[1], "message": "Hi everyone!"},
        {"conversation_idx": 12, "sender": users[5], "message": "Glad to join!"},
        {"conversation_idx": 12, "sender": users[6], "message": "Hello!"}
    ]
    msg_ids = []
    for msg in messages:
        msg_id = make_id()
        msg_ids.append(msg_id)
        conv_id = conv_ids[msg["conversation_idx"]]
        c.execute("INSERT INTO messages (id, conversation_id, sender_id, message) VALUES (?, ?, ?, ?)", (msg_id, conv_id, msg["sender"]["id"], msg["message"]))

    # REACTIONS (each user reacts to the first message)
    emojis = ["👍", "😂", "🔥", "❤️", "🎉"]
    reactions = []
    for i, user in enumerate(users):
        reactions.append({
            "id": make_id(),
            "message_id": msg_ids[0],
            "sender_id": user["id"],
            "emoji": emojis[i % len(emojis)]
        })
    for reaction in reactions:
        c.execute("INSERT INTO reactions (id, message_id, sender_id, emoji) VALUES (?, ?, ?, ?)", (reaction["id"], reaction["message_id"], reaction["sender_id"], reaction["emoji"]))

    # COMMENTS (each user comments on the second message)
    comments = []
    for i, user in enumerate(users):
        comments.append({
            "id": make_id(),
            "message_id": msg_ids[1],
            "sender_id": user["id"],
            "comment": f"Comment {i+1} from {user['username']}"
        })
    for comment in comments:
        c.execute("INSERT INTO comments (id, message_id, sender_id, comment) VALUES (?, ?, ?, ?)", (comment["id"], comment["message_id"], comment["sender_id"], comment["comment"]))

    conn.commit()
    conn.close()
    print("Demo data inserted for users, contacts, conversations, messages, reactions, and comments.")

if __name__ == "__main__":
    main()