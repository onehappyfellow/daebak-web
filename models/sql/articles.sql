CREATE TABLE articles(
    id SERIAL PRIMARY KEY,
    slug TEXT UNIQUE NOT NULL,
    headline TEXT,
    content TEXT,
    date TIMESTAMPTZ DEFAULT Now(),
    published BOOLEAN DEFAULT false,
    author TEXT   
);
