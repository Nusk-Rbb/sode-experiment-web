DROP TABLE IF EXISTS home_location;
DROP TABLE IF EXISTS user_location;
DROP TABLE IF EXISTS users;

-- テーブル作成
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS home_location (
    id SERIAL PRIMARY KEY,
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_location (
    id SERIAL PRIMARY KEY,
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL
);

-- テーブルが空の場合のみ初期データを挿入
INSERT INTO home_location (latitude, longitude)
SELECT 35.681236, 139.767125
WHERE NOT EXISTS (SELECT 1 FROM home_location);