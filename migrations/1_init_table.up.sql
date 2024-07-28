CREATE TABLE IF NOT EXISTS messages (
  id serial PRIMARY KEY,
  content text NOT NULL,
  processed boolean DEFAULT FALSE, 
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  processed_at timestamp
);
