CREATE TABLE articles(
    id SERIAL PRIMARY KEY,
    slug TEXT UNIQUE NOT NULL,
    
);




CREATE TABLE views (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    date TIMESTAMPTZ DEFAULT Now(),
    user INT,
    device TEXT,
    ip TEXT
);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT,
    level INT
);

CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    is_published BOOLEAN,
    date TIMESTAMPTZ DEFAULT Now(),
    headline TEXT,
    content TEXT,
    context TEXT,
    lesson TEXT,
    questions TEXT,
    reference INT
);

CREATE TABLE article_translations (
    id SERIAL PRIMARY KEY,
    article_id INT NOT NULL,
    language TEXT NOT NULL,
    headline TEXT,
    content TEXT,
    context TEXT,
    lesson TEXT
);

-- many to many relationship
CREATE TABLE article_tags (
    id SERIAL PRIMARY KEY,
    article_id INT NOT NULL,
    tag_id INT NOT NULL
);

CREATE TABLE article_relations (
    id SERIAL PRIMARY KEY,
    a_id INT NOT NULL,
    b_id INT NOT NULL
);

-- do generated articles have a reference?
CREATE TABLE references (
    id SERIAL PRIMARY KEY,
    author TEXT,
    publication TEXT,
    url TEXT,
    access_date TIMESTAMPTZ DEFAULT Now()
);

-- many to many
CREATE TABLE related_vocab (
    id SERIAL PRIMARY KEY,
    a_id INT NOT NULL,
    b_id INT NOT NULL
);
