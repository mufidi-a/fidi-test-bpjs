CREATE TABLE IF NOT EXISTS skill(
profile_code int,
id SERIAL PRIMARY KEY NOT NULL,
skill varchar,
level varchar,
created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT education_fk FOREIGN KEY (profile_code) REFERENCES profile(profile_code));