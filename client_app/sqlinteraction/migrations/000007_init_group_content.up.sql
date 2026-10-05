CREATE TABLE IF NOT EXISTS posts_of_group (
    id TEXT PRIMARY KEY,
    group_id TEXT,
    member_id TEXT,
    title TEXT NOT NULL, 
    content TEXT NOT NULL,
    image TEXT, 
    comments INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (member_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE TABLE IF NOT EXISTS post_group_rating (
    user_id TEXT NOT NULL,
    post_id TEXT NOT NULL,
    rate INTEGER NOT NULL,
    FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE,
    UNIQUE(user_id, post_id)
);
CREATE TABLE IF NOT EXISTS group_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id TEXT NOT NULL,
    group_id TEXT NOT NULL,
    content TEXT NOT NULL,
    sent_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS group_comments (
    id TEXT PRIMARY KEY,
    group_post_id TEXT NOT NULL,
    member_id TEXT NOT NULL,
    content TEXT ,
    image TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_post_id) REFERENCES group_posts(id) ON DELETE CASCADE,
    FOREIGN KEY (member_id) REFERENCES users(id) ON DELETE CASCADE
);