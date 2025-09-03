-- Migration to add unique constraint to reactions table
-- This prevents duplicate reactions from the same user for the same message

-- First, remove any existing duplicate reactions
DELETE FROM reactions 
WHERE id NOT IN (
    SELECT MIN(id) 
    FROM reactions 
    GROUP BY message_id, sender_id, emoji
);

-- Create new table with unique constraint
CREATE TABLE reactions_new (
    id TEXT PRIMARY KEY,
    message_id TEXT NOT NULL,
    sender_id TEXT NOT NULL,
    emoji TEXT NOT NULL,
    FOREIGN KEY(message_id) REFERENCES messages(id) ON DELETE CASCADE,
    UNIQUE(message_id, sender_id, emoji)
);

-- Copy data from old table
INSERT INTO reactions_new SELECT * FROM reactions;

-- Drop old table and rename new one
DROP TABLE reactions;
ALTER TABLE reactions_new RENAME TO reactions;
