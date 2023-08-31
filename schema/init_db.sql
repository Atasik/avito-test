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
    percentage real,
    name varchar(255) not null unique
);

CREATE TABLE users_segments
(
    user_id int references users (id) on delete cascade not null,
    seg_id int references segments (id) on delete cascade not null,
    expired_at timestamp,
    PRIMARY KEY(user_id, seg_id)
);

CREATE TABLE history
(
    id serial not null unique,
    user_id int references users (id) on delete cascade not null,
    segment varchar(255) not null,
    operation varchar(255) not null,
    created_at timestamp not null
);

CREATE OR REPLACE FUNCTION history_users_segments_insert()
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

CREATE OR REPLACE FUNCTION history_users_segments_delete()
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

CREATE OR REPLACE FUNCTION auto_users_segments_insert()
RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO users_segments (user_id, seg_id, expired_at)
    SELECT NEW.user_id, subquery.seg_id, NULL
    FROM (
        SELECT
            s.id AS seg_id,
            s.percentage AS segment_percentage,
            COALESCE(COUNT(us.seg_id), 0) * 100.0 / NULLIF((SELECT COUNT(*) FROM users_segments), 0) AS fill_percentage
        FROM segments s
        LEFT JOIN users_segments us ON us.seg_id = s.id
        GROUP BY s.id, s.percentage
        ORDER BY s.id
    ) AS subquery
    WHERE subquery.segment_percentage > subquery.fill_percentage;
    
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';


CREATE TRIGGER history_insert
    AFTER INSERT ON users_segments
    FOR EACH ROW
    EXECUTE FUNCTION history_users_segments_insert();

CREATE TRIGGER history_delete
    AFTER DELETE ON users_segments
    FOR EACH ROW
    EXECUTE FUNCTION history_users_segments_delete();

CREATE TRIGGER auto_insert
    AFTER INSERT ON users_segments
    FOR EACH ROW
    EXECUTE FUNCTION auto_users_segments_insert();
