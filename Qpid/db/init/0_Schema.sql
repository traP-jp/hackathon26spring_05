CREATE TABLE users (
  username VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  major VARCHAR(255) NOT NULL,
  hometown VARCHAR(255),
  like_topic VARCHAR(255),
  like_value VARCHAR(255),
  dislike_topic VARCHAR(255),
  dislike_value VARCHAR(255),
  tool VARCHAR(255),
  usual_situation VARCHAR(255),
  bio TEXT,
  PRIMARY KEY (username)
);

CREATE INDEX idx_like_topic ON users (like_topic);
CREATE INDEX idx_like_value ON users (like_value);
CREATE INDEX idx_dislike_topic ON users (dislike_topic);
CREATE INDEX idx_dislike_value ON users (dislike_value);
CREATE INDEX idx_tool ON users (tool);

CREATE TABLE tags (
  username VARCHAR(36) NOT NULL,
  name VARCHAR(255) NOT NULL,
  PRIMARY KEY (username, name)
);

CREATE INDEX idx_tag_name ON tags (name);

CREATE TABLE actions (
  id VARCHAR(36) NOT NULL,
  from_username VARCHAR(36) NOT NULL,
  to_username VARCHAR(36) NOT NULL,
  from_userstatus TINYINT NOT NULL,
  to_userstatus TINYINT NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT uq_action UNIQUE (from_username, to_username)
);

CREATE INDEX idx_userstatus ON actions (from_userstatus,to_userstatus);

CREATE TABLE icons(
  username VARCHAR(36) NOT NULL,
  icon LONGBLOB NOT NULL,
  mime_type VARCHAR(50) NOT NULL,
  PRIMARY KEY (username)
);