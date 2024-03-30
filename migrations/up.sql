CREATE TABLE commands (
    command_id SERIAL PRIMARY KEY,
    command_text VARCHAR NOT NULL,
    last_output VARCHAR NOT NULL DEFAULT ''
);
