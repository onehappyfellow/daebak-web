CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    parent_id INT REFERENCES tags(id),
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE article_tags (
    article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (article_id, tag_id)
);
