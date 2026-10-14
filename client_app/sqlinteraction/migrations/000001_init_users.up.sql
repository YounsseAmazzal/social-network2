CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    nickName TEXT UNIQUE,
    firstName TEXT,
    lastName TEXT,
    age INTEGER,
    gender TEXT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    avatarURL TEXT,
    about_me TEXT DEFAULT '',
    is_private INTEGER NOT NULL DEFAULT 0
);

