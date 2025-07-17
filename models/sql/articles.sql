CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    uuid TEXT UNIQUE NOT NULL,
    published BOOLEAN DEFAULT false,
    source_published TIMESTAMP,
    source_accessed TIMESTAMP DEFAULT now(),
    source_url TEXT,
    source_publication TEXT,
    source_author TEXT,
    headline TEXT NOT NULL,
    headline_en TEXT,
    content JSONB,
    summary TEXT,
    context TEXT,
    topik_level INT,
    topik_level_explanation TEXT,
    comprehension_questions TEXT
);
