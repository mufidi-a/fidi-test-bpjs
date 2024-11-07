CREATE DATABASE IF NOT EXISTS working_experiences
id int PRIMARY KEY NOT NULL AUTO_INCREMENT,
profile_code int,
working_experiences varchar,
CONSTRAINT working_experiences_fk FOREIGN KEY (profile_code) REFERENCES profile(profile_code),
CONSTRAINT working_experiences_un UNIQUE (profile_code);