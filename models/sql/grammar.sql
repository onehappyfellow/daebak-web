CREATE TABLE grammar (
    id SERIAL PRIMARY KEY,
    published BOOLEAN DEFAULT false,
    title TEXT UNIQUE NOT NULL,
    explanation TEXT,
    explanation_short TEXT,
    examples TEXT,
    practice TEXT
);

CREATE TABLE article_grammar (
    grammar_id INT NOT NULL REFERENCES grammar(id) ON DELETE CASCADE,
    article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    article_example TEXT,
    PRIMARY KEY (grammar_id, article_id)
);