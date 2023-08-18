CREATE TABLE `non_repeating_device_alarm_id`
(
    `alarm_id`        varchar(50) NOT NULL,
    `device_alarm_id` int DEFAULT NULL,
    PRIMARY KEY (`alarm_id`),
    KEY               `fk_nrda_alarm_id_idx` (`alarm_id`),
    CONSTRAINT `fk_nrda_alarm_id` FOREIGN KEY (`alarm_id`) REFERENCES `alarms` (`alarm_id`)
)