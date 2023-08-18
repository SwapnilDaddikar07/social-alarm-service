CREATE TABLE `alarm_media`
(
    `alarm_id` varchar(50) DEFAULT NULL,
    `media_id` varchar(50) DEFAULT NULL,
    KEY        `fk_am_alarm_id_idx` (`alarm_id`),
    KEY        `fk_am_media_id_idx` (`media_id`),
    CONSTRAINT `fk_am_alarm_id` FOREIGN KEY (`alarm_id`) REFERENCES `alarms` (`alarm_id`),
    CONSTRAINT `fk_am_media_id` FOREIGN KEY (`media_id`) REFERENCES `media` (`media_id`)
);