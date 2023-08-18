CREATE TABLE `alarms`
(
    `alarm_id`             varchar(50) NOT NULL,
    `user_id`              varchar(50)  DEFAULT NULL,
    `visibility`           varchar(20)  DEFAULT NULL,
    `description`          varchar(100) DEFAULT NULL,
    `status`               varchar(10)  DEFAULT NULL,
    `created_at`           datetime     DEFAULT NULL,
    `alarm_start_datetime` datetime     DEFAULT NULL,
    PRIMARY KEY (`alarm_id`),
    KEY                    `alarm_usr_fk_key_idx` (`user_id`),
    CONSTRAINT `fk_alarm_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
);