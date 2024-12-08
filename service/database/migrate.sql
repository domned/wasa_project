PRAGMA foreign_keys = ON;


CREATE TABLE IF NOT EXISTS users{
    ID string PRIMARY KEY NOT NULL,
    name string UNIQUE NOT NULL
};

CREATE TABLE IF NOT EXISTS chats{
    ID string NOT NULL
};

CREATE TABLE IF NOT EXISTS groups{
    members string NOT NULL, 
};

