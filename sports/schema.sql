-- events table to store sports events
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY,
    name TEXT,
    category TEXT,
    competition TEXT,
    visible INTEGER,
    advertised_start_time DATETIME
);

-- Add an index on category to optimize query filtering by this column
CREATE INDEX IF NOT EXISTS idx_events_category ON events(category);

-- Add an index on visible to optimize query filtering by this column
CREATE INDEX IF NOT EXISTS idx_events_visible ON events(visible);
