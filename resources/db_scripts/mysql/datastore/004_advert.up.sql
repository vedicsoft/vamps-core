CREATE TABLE IF NOT EXISTS `ad_campaign` (
  `tenantid`      INT(11) DEFAULT NULL,
  `campaign_id`   INT(11)      NOT NULL                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                AUTO_INCREMENT,
  `name`          VARCHAR(200) NOT NULL,
  `campaign_type` VARCHAR(30)  NOT NULL                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               DEFAULT 'Banner Ad',
  `campaign_tier` VARCHAR(30)  NOT NULL                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               DEFAULT 'Normal',
  `start_date`    DATE         NOT NULL,
  `end_date`      DATE         NOT NULL,
  `status`        VARCHAR(10)  NOT NULL                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               DEFAULT 'Active',
  `weight`        TINYINT(2)   NOT NULL                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               DEFAULT '5',
  `visibility`    VARCHAR(6)   NOT NULL                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               DEFAULT 'true',
  PRIMARY KEY (`campaign_id`),
  FOREIGN KEY (`tenantid`) REFERENCES `vs_tenants` (`tenantid`)
    ON DELETE CASCADE
)
  ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

CREATE TABLE IF NOT EXISTS `ad_zone` (
  `tenantid`             INT(11) DEFAULT NULL,
  `zone_id`              INT(11)      NOT NULL                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     AUTO_INCREMENT,
  `name`                 VARCHAR(200) NOT NULL,
  `zone_type`            VARCHAR(200) NOT NULL,
  `status`               VARCHAR(45)  NOT NULL,
  `auto_play_video`      VARCHAR(20)  NOT NULL                                                                                                                                                                                                                                                                                                                                                                                              DEFAULT 'Deactivate',
  `loop_video`           VARCHAR(20)  NOT NULL                                                                                                                                                                                                                                                                                                                                                                                              DEFAULT 'Deactivate',
  `mute`                 VARCHAR(20)  NOT NULL                                                                                                                                                                                                                                                                                                                                                                                              DEFAULT 'Deactivate',
  `show_video_contoller` VARCHAR(20)  NOT NULL                                                                                                                                                                                                                                                                                                                                                                                              DEFAULT 'Deactivate',
  PRIMARY KEY (`zone_id`),
  FOREIGN KEY (`tenantid`) REFERENCES `vs_tenants` (`tenantid`)
    ON DELETE CASCADE
)
  ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

CREATE TABLE IF NOT EXISTS `ad_creative` (
  `tenantid`     INT(11),
  `creative_id`  INT(11)      NOT NULL AUTO_INCREMENT,
  `campaign_id`  INT(11)      NOT NULL,
  `zone_type`    VARCHAR(100) NOT NULL,
  `name`         VARCHAR(200) NOT NULL,
  `Item_url`     MEDIUMTEXT   NOT NULL,
  `width`        SMALLINT(5)           DEFAULT NULL,
  `height`       SMALLINT(5)           DEFAULT NULL,
  `alter_text`   VARCHAR(200)          DEFAULT 'ad not found',
  `status`       VARCHAR(50)           DEFAULT NULL,
  `landing_page` MEDIUMTEXT   NOT NULL,
  PRIMARY KEY (`creative_id`, `campaign_id`),
  FOREIGN KEY (`campaign_id`) REFERENCES `ad_campaign` (`campaign_id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  FOREIGN KEY (`tenantid`) REFERENCES `vs_tenants` (`tenantid`)
    ON DELETE CASCADE
)
  ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

CREATE TABLE IF NOT EXISTS `ad_zone_creative` (
  `creative_id` INT(11) NOT NULL,
  `zone_id`     INT(11) NOT NULL,
  PRIMARY KEY (`creative_id`,`zone_id`),
  FOREIGN KEY (`zone_id`) REFERENCES `ad_zone` (`zone_id`)
    ON DELETE CASCADE,
  FOREIGN KEY (`creative_id`) REFERENCES `ad_creative` (`creative_id`)
    ON DELETE CASCADE
)
  ENGINE = InnoDB
  DEFAULT CHARSET = latin1;


CREATE TABLE IF NOT EXISTS `ad_insight` (
  `insight_id` int(11) NOT NULL AUTO_INCREMENT,
  `tenant_id` int(11) NOT NULL,
  `visitor_id` varchar(100) NOT NULL,
  `subscriberID` varchar(45) DEFAULT NULL,
  `username` varchar(100) DEFAULT NULL,
  `campaign_id` int(11) NOT NULL,
  `campaign_name` varchar(100) NOT NULL,
  `creative_id` int(11) NOT NULL,
  `creative_name` varchar(100) NOT NULL,
  `banner_name` tinytext NOT NULL,
  `browser_name` varchar(200) NOT NULL,
  `language` varchar(100) NOT NULL,
  `os` varchar(100) NOT NULL,
  `device` varchar(100) NOT NULL,
  `traffic` varchar(250) NOT NULL,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`insight_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `ad_banners` (
  `tenant_id` int(11) NOT NULL,
  `ad_banner_id` int(11) NOT NULL AUTO_INCREMENT,
  `ad_banner_name` varchar(500) NOT NULL,
  PRIMARY KEY (`ad_banner_id`,`tenant_id`),
  FOREIGN KEY (`tenant_id`) REFERENCES `vs_tenants` (`tenantid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `ad_target_list` (
  `target_list_id` int(11) NOT NULL AUTO_INCREMENT,
  `tenant_id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `description` varchar(500) NOT NULL,
  PRIMARY KEY (`target_list_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `ad_target_elements` (
  `element_id` int(11) NOT NULL,
  `target_list_id` int(11) NOT NULL,
  PRIMARY KEY (`element_id`,`target_list_id`),
  FOREIGN KEY (`target_list_id`) REFERENCES `ad_target_list` (`target_list_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



CREATE TABLE IF NOT EXISTS `ad_target_list_campaign` (
  `list_id` int(11) NOT NULL,
  `campaign_id` int(11) NOT NULL,
  PRIMARY KEY (`list_id`,`campaign_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `ad_user_interaction` (
  `interaction_id` int(11) NOT NULL AUTO_INCREMENT,
  `tenant_id` int(11) NOT NULL,
  `visitor_id` varchar(100) NOT NULL,
  `campaign_id` int(11) NOT NULL,
  `creative_id` int(11) NOT NULL,
  `username` varchar(100) NOT NULL,
  `campaign_name` varchar(100) NOT NULL,
  `creative_name` varchar(100) NOT NULL,
  `banner_name` varchar(100) NOT NULL,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`interaction_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;