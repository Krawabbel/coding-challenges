DROP TABLE IF EXISTS urls;

CREATE TABLE urls (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    short VARCHAR(100) NOT NULL UNIQUE,
    uri TEXT NOT NULL,
    creator VARCHAR(100) NOT NULL,
    created DATETIME NOT NULL DEFAULT current_timestamp,
    expires DATETIME NOT NULL DEFAULT (datetime('now','+1 day'))
);

CREATE INDEX idx_short ON urls(short);

INSERT INTO urls (short, uri, creator) VALUES ('gidf', 'https://www.google.com', 'default setup');

SELECT * FROM urls;