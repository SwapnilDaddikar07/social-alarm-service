CREATE TABLE `users`
(
    `user_id`      varchar(50) NOT NULL,
    `display_name` varchar(100) DEFAULT NULL,
    `phone_number` varchar(45)  DEFAULT NULL,
    PRIMARY KEY (`user_id`)
);