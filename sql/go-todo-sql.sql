CREATE TABLE IF NOT EXISTS user(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(20) NOT NULL,
    password_hash VARCHAR(128) NOT NULL,
    webhook VARCHAR(256)
);

CREATE TABLE IF NOT EXISTS message(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(20) NOT NULL,
    message VARCHAR(128) NOT NULL,
    create_time INTEGER NOT NULL ,
    finished_time INTEGER,
    notify_time INTEGER
);

CREATE INDEX IF NOT EXISTS idx_username ON message(username);