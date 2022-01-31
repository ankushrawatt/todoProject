CREATE TABLE session(
    token TEXT NOT NULL,
    userid TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);