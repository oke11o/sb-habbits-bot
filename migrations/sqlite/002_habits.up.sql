CREATE TABLE IF NOT EXISTS habits
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL,
    name        TEXT    NOT NULL,
    type        TEXT    NOT NULL,
    target      INTEGER   DEFAULT 0,
    target_time TEXT      DEFAULT NULL,
    max_time    TEXT      DEFAULT NULL,
    unit        TEXT      DEFAULT NULL,
    points      INTEGER   DEFAULT 0,
    points_mode TEXT      DEFAULT 'fixed',
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_habits_user_id_name ON habits (user_id, name);

CREATE TABLE IF NOT EXISTS records
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    habit_id  INTEGER NOT NULL,
    user_id   INTEGER NOT NULL,
    value     INTEGER   DEFAULT 0,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    points    INTEGER   DEFAULT 0,
    FOREIGN KEY (habit_id) REFERENCES habits (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS reminders
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    habit_id INTEGER NOT NULL,
    user_id  INTEGER NOT NULL,
    time     TEXT,
    days     TEXT,
    FOREIGN KEY (habit_id) REFERENCES habits (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS progress
(
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    habit_id          INTEGER NOT NULL,
    user_id           INTEGER NOT NULL,
    accumulated_value INTEGER   DEFAULT 0,
    target            INTEGER   DEFAULT 0,
    last_updated      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (habit_id) REFERENCES habits (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);