-- Race table to store horse racing events
CREATE TABLE IF NOT EXISTS races (
    id INTEGER PRIMARY KEY,
    meeting_id INTEGER,
    name TEXT,
    number INTEGER,
    visible INTEGER,
    advertised_start_time DATETIME
);

-- Add an index on meeting_id to optimize query filtering by this column
CREATE INDEX IF NOT EXISTS idx_races_meeting_id ON races(meeting_id);
