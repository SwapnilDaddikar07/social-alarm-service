CREATE TABLE `repeating_device_alarm_id`
(
    `alarm_id`            varchar(50) DEFAULT NULL,
    `mon_device_alarm_id` int         DEFAULT NULL,
    `tue_device_alarm_id` int         DEFAULT NULL,
    `wed_device_alarm_id` int         DEFAULT NULL,
    `thu_device_alarm_id` int         DEFAULT NULL,
    `fri_device_alarm_id` int         DEFAULT NULL,
    `sat_device_alarm_id` int         DEFAULT NULL,
    `sun_device_alarm_id` int         DEFAULT NULL,
    KEY                   `fk_rda_alarm_id_idx` (`alarm_id`),
    CONSTRAINT `fk_rda_alarm_id` FOREIGN KEY (`alarm_id`) REFERENCES `alarms` (`alarm_id`)
)