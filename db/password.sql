CREATE TABLE IF NOT EXISTS `password`
(
    `id`       INTEGER PRIMARY KEY AUTO_INCREMENT,
    `user_id`  TEXT NOT NULL,
    `resource` TEXT NOT NULL,
    `login`    TEXT NOT NULL,
    `password` TEXT NOT NULL
);