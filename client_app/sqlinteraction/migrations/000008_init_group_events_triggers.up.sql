CREATE TABLE IF NOT EXISTS group_events (
    id TEXT PRIMARY KEY,
    group_id TEXT,
    member_id TEXT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    event_date DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    total_going INTEGER DEFAULT 0,
    total_not_going INTEGER DEFAULT 0,
    status TEXT,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (member_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS group_events_rating (
    id TEXT PRIMARY KEY,
    event_id TEXT,
    member_id TEXT,
    status TEXT CHECK(status IN ('going', 'not going')),
    FOREIGN KEY (event_id) REFERENCES group_events(id) ON DELETE CASCADE,
    FOREIGN KEY (member_id) REFERENCES users(id) ON DELETE CASCADE
);

-- TRIGGERS
CREATE TRIGGER IF NOT EXISTS upsert_vote
BEFORE INSERT ON group_events_rating
FOR EACH ROW
BEGIN
    DELETE FROM group_events_rating
    WHERE event_id = NEW.event_id
    AND member_id = NEW.member_id;
END;

CREATE TRIGGER IF NOT EXISTS update_vote_counts_after_insert
AFTER INSERT ON group_events_rating
FOR EACH ROW
BEGIN
    UPDATE group_events
    SET
        total_going = (SELECT COUNT(*) FROM group_events_rating WHERE event_id = NEW.event_id AND status = 'going'),
        total_not_going = (SELECT COUNT(*) FROM group_events_rating WHERE event_id = NEW.event_id AND status = 'not going')
    WHERE id = NEW.event_id;
END;

CREATE TRIGGER IF NOT EXISTS update_vote_counts_after_delete
AFTER DELETE ON group_events_rating
FOR EACH ROW
BEGIN
    UPDATE group_events
    SET
        total_going = (SELECT COUNT(*) FROM group_events_rating WHERE event_id = OLD.event_id AND status = 'going'),
        total_not_going = (SELECT COUNT(*) FROM group_events_rating WHERE event_id = OLD.event_id AND status = 'not going')
    WHERE id = OLD.event_id;
END;