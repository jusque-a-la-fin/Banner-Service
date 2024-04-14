-- название БД: banner_service. Она хранит всю информацию о баннерах
BEGIN;
-- Фичи
CREATE TABLE features (
    id SERIAL PRIMARY KEY,
    -- ID фичи
    feature_id INT UNIQUE NOT NULL
);

-- Баннеры
CREATE TABLE banners (
    id SERIAL PRIMARY KEY,
    feature_id INT NOT NULL,
    -- информация баннера
    content JSONB UNIQUE NOT NULL,
    -- время создания
    created_at VARCHAR(160) NOT NULL,
    -- время обновления
    updated_at VARCHAR(160) NOT NULL,
     -- флаг(включен или выключен баннер)
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    FOREIGN KEY (feature_id) REFERENCES features(id)
);

-- Теги
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    -- ID тега
    tag_id INT UNIQUE NOT NULL
);

-- Связи баннеров с тегами
CREATE TABLE banners_tags (
    banner_id INT,
    tag_id INT,
    PRIMARY KEY (banner_id, tag_id),
    FOREIGN KEY (banner_id) REFERENCES banners(id),
    FOREIGN KEY (tag_id) REFERENCES tags(tag_id)
);


INSERT INTO tags (tag_id) VALUES
    (1),
    (2),
    (3),
    (4),
    (5),
    (6),
    (7),
    (8),
    (9),
    (10);


INSERT INTO features (feature_id) VALUES
    (1),
    (2),
    (3),
    (4),
    (5),
    (6),
    (7),
    (8);


INSERT INTO banners (feature_id, content, created_at, updated_at, is_active) VALUES
    (1, '{"title": "some_title_1", "text": "some_text_1", "url": "some_url_1"}', '2024-04-11 16:42:30.946058783 +0300 MSK', '2024-04-11 16:42:30.946058783 +0300 MSK', TRUE),
    (2, '{"title": "some_title_2", "text": "some_text_2", "url": "some_url_2"}', '2024-04-11 16:42:30.946058783 +0300 MSK', '2024-04-11 16:42:30.946058783 +0300 MSK', TRUE),
    (3, '{"title": "some_title_3", "text": "some_text_3", "url": "some_url_3"}', '2024-04-11 16:42:30.946058783 +0300 MSK', '2024-04-11 16:42:30.946058783 +0300 MSK', TRUE),
    (4, '{"title": "some_title_4", "text": "some_text_4", "url": "some_url_4"}', '2024-04-11 16:42:30.946058783 +0300 MSK', '2024-04-11 16:42:30.946058783 +0300 MSK', FALSE),
    (5, '{"title": "some_title_5", "text": "some_text_5", "url": "some_url_5"}', '2024-04-11 16:42:30.946058783 +0300 MSK', '2024-04-11 16:42:30.946058783 +0300 MSK', TRUE),
    (6, '{"title": "some_title_6", "text": "some_text_6", "url": "some_url_6"}', '2024-04-11 16:42:30.946058783 +0300 MSK', '2024-04-11 16:42:30.946058783 +0300 MSK', TRUE),
    (7, '{"title": "some_title_7", "text": "some_text_7", "url": "some_url_7"}', '2024-04-11 16:42:30.946058783 +0300 MSK', '2024-04-11 16:42:30.946058783 +0300 MSK', FALSE),
    (8, '{"title": "some_title_8", "text": "some_text_8", "url": "some_url_8"}', '2024-04-11 16:42:30.946058783 +0300 MSK', '2024-04-11 16:42:30.946058783 +0300 MSK', TRUE),
    (3, '{"title": "some_title_9", "text": "some_text_9", "url": "some_url_9"}', '2024-04-11 16:42:30.946058783 +0300 MSK', '2024-04-11 16:42:30.946058783 +0300 MSK', TRUE),
    (6, '{"title": "some_title_10", "text": "some_text_10", "url": "some_url_10"}', '2024-04-11 16:42:30.946058783 +0300 MSK', '2024-04-11 16:42:30.946058783 +0300 MSK', TRUE);


INSERT INTO banners_tags (banner_id, tag_id) VALUES
    (1, 1), (1, 2), (1, 5),
    (2, 3),
    (3, 2), (3, 3), (3, 4), (3, 5),
    (4, 7), (4, 8),
    (5, 9), (5, 10),
    (6, 1), (6, 2), (6, 3), (6, 4), (6, 5),
    (7, 3), (7, 4),
    (8, 5), (8, 6),
    (9, 7),
    (10, 9), (10, 10);

COMMIT;
