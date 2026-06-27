CREATE TABLE users (
  id VARCHAR(36) NOT NULL,
  username VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  major VARCHAR(255) NOT NULL,
  hometown VARCHAR(255),
  like_topic VARCHAR(255),
  like_value VARCHAR(255),
  dislike_topic VARCHAR(255),
  dislike_value VARCHAR(255),
  bio TEXT,
  PRIMARY KEY (id)
);

CREATE TABLE tags (
  user_id VARCHAR(36) NOT NULL,
  name VARCHAR(255) NOT NULL,
  affinity TINYINT NOT NULL,
  strength FLOAT NOT NULL,
  PRIMARY KEY (user_id, name)
);

CREATE TABLE affiliations (
  user_id VARCHAR(36) NOT NULL,
  team VARCHAR(255) NOT NULL,
  PRIMARY KEY (user_id, team)
);

CREATE TABLE actions (
  id VARCHAR(36) NOT NULL,
  from_userid VARCHAR(36) NOT NULL,
  to_userid VARCHAR(36) NOT NULL,
  status TINYINT NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE icons(
  id VARCHAR(36) NOT NULL,
  icon LONGBLOB NOT NULL
  mime_type VARCHAR(50) NOT NULL
  PRIMARY KEY (user_id)
);