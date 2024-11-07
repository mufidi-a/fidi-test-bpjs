CREATE TABLE IF NOT EXISTS education(
profile_code int,
id SERIAL PRIMARY KEY NOT NULL,
school varchar,
degree varchar,
start_date DATE,
end_date DATE,
city varchar,
description varchar,
created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT education_fk FOREIGN KEY (profile_code) REFERENCES profile(profile_code));