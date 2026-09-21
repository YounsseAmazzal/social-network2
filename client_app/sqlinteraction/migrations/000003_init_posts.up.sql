CREATE TABLE IF NOT EXISTS post (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    authorID TEXT NOT NULL,
    Image TEXT,
    Privecytype TEXT,
    createDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (authorID) REFERENCES user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS post_shared_users (
    post_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS post_vote (
    user_id TEXT NOT NULL,
    post_id TEXT NOT NULL,
    vote INTEGER NOT NULL,
    FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE,
    UNIQUE(user_id, post_id)
);
