CREATE TABLE users (
  username VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  major VARCHAR(255) NOT NULL,
  hometown VARCHAR(255),
  like_topic VARCHAR(255),
  like_value VARCHAR(255),
  dislike_topic VARCHAR(255),
  dislike_value VARCHAR(255),
  bio TEXT,
  PRIMARY KEY (username)
);

CREATE TABLE tags (
  username VARCHAR(36) NOT NULL,
  name VARCHAR(255) NOT NULL,
  affinity TINYINT NOT NULL,
  strength FLOAT NOT NULL,
  PRIMARY KEY (username, name)
);

CREATE TABLE affiliations (
  username VARCHAR(36) NOT NULL,
  team VARCHAR(255) NOT NULL,
  PRIMARY KEY (username, team)
);

CREATE TABLE actions (
  id VARCHAR(36) NOT NULL,
  from_username VARCHAR(36) NOT NULL,
  to_username VARCHAR(36) NOT NULL,
  status TINYINT NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE icons(
  username VARCHAR(36) NOT NULL,
  icon LONGBLOB NOT NULL
  mime_type VARCHAR(50) NOT NULL
  PRIMARY KEY (username)
);