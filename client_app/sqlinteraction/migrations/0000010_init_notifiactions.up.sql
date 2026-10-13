CREATE TABLE IF NOT EXISTS notifications (
    id TEXT PRIMARY KEY,
    recipient_id TEXT NOT NULL,      
    sender_id TEXT,                
    notification_category TEXT NOT NULL, 
    target_reference_id TEXT,       
    target_reference_type TEXT,     
    display_message TEXT NOT NULL,  
    read_status INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(recipient_id, sender_id, notification_category, target_reference_id)
    FOREIGN KEY (recipient_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_notifications_lookup 
ON notifications(recipient_id, read_status, created_at);