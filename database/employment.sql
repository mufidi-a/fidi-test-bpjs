CREATE TABLE IF NOT EXISTS employment(
profile_code int,
id SERIAL PRIMARY KEY NOT NULL,
job_title varchar,
employer varchar,
start_date DATE,
end_date DATE,
city varchar,
description varchar,
created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT employment_fk FOREIGN KEY (profile_code) REFERENCES profile(profile_code));