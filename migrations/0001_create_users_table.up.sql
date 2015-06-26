CREATE TABLE users(
  id            VARCHAR(255) PRIMARY KEY NOT NULL,
  encoded_id    VARCHAR(255) NOT NULL,
  gender        VARCHAR(255) NOT NULL,
  date_of_birth DATE NOT NULL,
  access_token  TEXT NOT NULL,
  refresh_token VARCHAR(255) NOT NULL,
  expiry        VARCHAR(255) NOT NULL,
  token_type    VARCHAR(255) NOT NULL,
  created_at    TIMESTAMP    NOT NULL,
  updated_at    TIMESTAMP    NOT NULL
);

CREATE INDEX users_encoded_id_idx ON users USING btree (encoded_id);
