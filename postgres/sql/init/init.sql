CREATE TABLE IF NOT EXISTS home_location (
    id SERIAL PRIMARY KEY,
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL
);

-- テーブルが空の場合のみ初期データを挿入
INSERT INTO home_location (latitude, longitude)
SELECT 35.681236, 139.767125
WHERE NOT EXISTS (SELECT 1 FROM home_location);