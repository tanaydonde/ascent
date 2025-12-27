CREATE TABLE IF NOT EXISTS problems (
    problem_id TEXT PRIMARY KEY,
    rating INT NOT NULL,
    tags TEXT[]
);

DO $$ BEGIN 
    CREATE TYPE solve_status as ENUM ('solved', 'partially_solved', 'failed');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS user_logs (
    id SERIAL PRIMARY KEY,
    handle TEXT NOT NULL,
    problem_id TEXT REFERENCES problems(problem_id),
    status solve_status NOT NULL,

    -- time_spent_minutes is NULL when data is unavailable
    time_spent_minutes INT,

    submission_count INT DEFAULT 1,
    is_api_synced BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS topic_dependencies (
    prereq_topic TEXT,
    target_topic TEXT,
    PRIMARY KEY (prereq_topic, target_topic)
);