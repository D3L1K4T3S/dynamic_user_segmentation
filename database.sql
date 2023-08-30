DROP TABLE IF EXISTS USERS;
CREATE TABLE USERS (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

DROP TABLE IF EXISTS SEGMENTS;
CREATE TABLE SEGMENTS (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    percent NUMERIC
);

DROP TABLE IF EXISTS CONSUMERS_SEGMENTS;
CREATE TABLE test (
    id SERIAL PRIMARY KEY,
    segment_id INT NOT NULL ,
    ttl TIMESTAMP,
    FOREIGN KEY (segment_id) REFERENCES SEGMENTS(id)
);

DROP TABLE IF EXISTS CONSUMERS;
CREATE TABLE CONSUMERS (
    id SERIAL PRIMARY KEY,
    consumer_id INT NOT NULL,
    segment_id INT,
    FOREIGN KEY (segment_id) REFERENCES CONSUMERS_SEGMENTS(id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS ACTIONS;
CREATE TABLE ACTIONS (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

DROP TABLE IF EXISTS OPERATIONS;
CREATE TABLE OPERATIONS (
    id SERIAL PRIMARY KEY,
    consumer_id INT NOT NULL,
    segment_id INT NOT NULL,
    action_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (segment_id) REFERENCES SEGMENTS(id) ON DELETE CASCADE,
    FOREIGN KEY (action_id) REFERENCES ACTIONS(id) ON DELETE CASCADE
);

CREATE FUNCTION delete_expired_ttl(c_id int) RETURNS void LANGUAGE plpgsql as
    $$
    begin
        DELETE FROM consumers_segments USING consumers
        WHERE consumers_segments.id = consumers.segment_id
        AND consumers.consumer_id = c_id
        AND consumers_segments.ttl IS NOT NULL
        AND consumers_segments.ttl < NOW();
    end
    $$
;

