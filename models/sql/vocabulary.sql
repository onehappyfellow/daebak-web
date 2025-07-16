CREATE TABLE vocabulary (
    id SERIAL PRIMARY KEY,
    word TEXT UNIQUE NOT NULL,
    definition TEXT,
    examples TEXT,
    translation_en TEXT
);

CREATE TABLE article_vocabulary (
    vocabulary_id INT NOT NULL REFERENCES vocabulary(id) ON DELETE CASCADE,
    article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    article_location TEXT, -- example: "{block: 1, start: 5, end: 10}"
    PRIMARY KEY (vocabulary_id, article_id)
);