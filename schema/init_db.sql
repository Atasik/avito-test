DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS segments;
DROP TABLE IF EXISTS users_segments;
DROP TABLE IF EXISTS history;

CREATE TABLE users
(
    id int not null unique
);

CREATE TABLE segments
(
    id serial not null unique,
    name varchar(255) not null unique
);

CREATE TABLE users_segments
(
    -- id serial not null unique,
    user_id int references users (id) on delete cascade not null,
    seg_id int references segments (id) on delete cascade not null,
    expired_at timestamp,
    PRIMARY KEY(user_id, seg_id)
);

CREATE TABLE history
(
    id serial not null unique,
    user_id int references users (id) on delete cascade not null,
    -- seg_id int references segments (id) on delete cascade not null,
    segment varchar(255) not null,
    operation varchar(255) not null,
    created_at timestamp not null
);

CREATE OR REPLACE FUNCTION update_insert_history()
RETURNS TRIGGER AS 
$$
BEGIN 
    INSERT INTO history(user_id, segment, operation, created_at)
    SELECT NEW.user_id, s.name, 'add', NOW()
    FROM segments s
    WHERE s.id = NEW.seg_id;
    
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION update_delete_history()
RETURNS TRIGGER AS 
$$
BEGIN
    INSERT INTO history(user_id, segment, operation, created_at)
    SELECT OLD.user_id, s.name, 'delete', NOW()
    FROM segments s
    WHERE s.id = OLD.seg_id;
    
    RETURN OLD;
END;
$$
LANGUAGE 'plpgsql';

CREATE TRIGGER insert_entry
    AFTER INSERT ON users_segments
    FOR EACH ROW
    EXECUTE FUNCTION update_insert_history();

CREATE TRIGGER delete_entry
    AFTER DELETE ON users_segments
    FOR EACH ROW
    EXECUTE FUNCTION update_delete_history();


INSERT INTO segments(name) VALUES('AVITO_VOICE_MESSAGES');
INSERT INTO segments(name) VALUES('AVITO_PERFORMANCE_VAS');
INSERT INTO segments(name) VALUES('AVITO_DISCOUNT_30');
INSERT INTO segments(name) VALUES('AVITO_DISCOUNT_50');

-- INSERT INTO users(ID) VALUES(1000);
-- INSERT INTO users(ID) VALUES(1002);
-- INSERT INTO users(ID) VALUES(1004);

-- INSERT INTO users_segments(user_id, seg_id) VALUES(1000, 1);
-- INSERT INTO users_segments(user_id, seg_id) VALUES(1000, 2);
-- INSERT INTO users_segments(user_id, seg_id) VALUES(1000, 3);
-- INSERT INTO users_segments(user_id, seg_id) VALUES(1002, 1);
-- INSERT INTO users_segments(user_id, seg_id) VALUES(1002, 4);