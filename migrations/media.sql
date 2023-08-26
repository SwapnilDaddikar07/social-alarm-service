CREATE TABLE `media`
(
    `media_id`     varchar(50) NOT NULL,
    `sender_id`    varchar(50)  DEFAULT NULL,
    `resource_url` varchar(100) DEFAULT NULL,
    `created_at`   datetime     DEFAULT NULL,
    PRIMARY KEY (`media_id`),
    KEY            `_idx` (`sender_id`),
    CONSTRAINT `` FOREIGN KEY (`sender_id`) REFERENCES `users` (`user_id`)
);
