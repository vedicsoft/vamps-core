-- provisioned methods : socialauth, pinauth, emailauth
CREATE TABLE IF NOT EXISTS `wf_subscribers` (
  `subscriberid`     BIGINT       NOT NULL AUTO_INCREMENT,
  `tenantid`         INT,
  `venueid`          INT,
  `username`         VARCHAR(255) CHARACTER SET utf8 DEFAULT NULL,
  `password`         VARCHAR(255) CHARACTER SET utf8 DEFAULT NULL,
  `email`            VARCHAR(255) CHARACTER SET utf8 NOT NULL,
  `account_status`   VARCHAR(255)          DEFAULT 'inactive',
  `social_id`        VARCHAR(255)          DEFAULT NULL,
  `first_name`       VARCHAR(255) CHARACTER SET utf8 DEFAULT NULL,
  `last_name`        VARCHAR(255)          DEFAULT NULL,
  `gender`           VARCHAR(255)          DEFAULT NULL,
  `birthday`         DATE,
  `age`              INT,
  `age_upper`        INT,
  `age_lower`        INT,
  `religion`         VARCHAR(255)          DEFAULT NULL,
  `occupation`       VARCHAR(255)          DEFAULT NULL,
  `marital_status`   VARCHAR(255)          DEFAULT NULL,
  `profile_image`    VARCHAR(255)          DEFAULT NULL,
  `address_line_1`   VARCHAR(255)          DEFAULT NULL,
  `address_line_2`   VARCHAR(255)          DEFAULT NULL,
  `city`             VARCHAR(255)          DEFAULT NULL,
  `state`            VARCHAR(255)          DEFAULT NULL,
  `zip`              VARCHAR(255)          DEFAULT NULL,
  `country`          VARCHAR(255)          DEFAULT NULL,
  `mobile_number`    VARCHAR(255)          DEFAULT NULL,
  `telephone_number` VARCHAR(255)          DEFAULT NULL,
  `admin_notes`      VARCHAR(255)          DEFAULT NULL,
  `created`          TIMESTAMP             DEFAULT now(),
  `updated`          TIMESTAMP             DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (`subscriberid`),
  UNIQUE KEY (`tenantid`, `username`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid) ON DELETE CASCADE,
  FOREIGN KEY (venueid) REFERENCES wf_venues (venueid) ON DELETE SET NULL
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_subscriber_devices` (
  `tenantid`         INT,
  `subscriberid`     BIGINT       NOT NULL,
  `mac`              VARCHAR(255) NOT NULL,
  `parental_control` VARCHAR(255) DEFAULT 'OFF',
  `device`           VARCHAR(255),
  `created`          TIMESTAMP    DEFAULT now(),
  `updated`          TIMESTAMP    DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (tenantid, subscriberid, mac),
  FOREIGN KEY (subscriberid) REFERENCES wf_subscribers (subscriberid)
    ON DELETE CASCADE,
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_subscriber_auth_log` (
  `tenantid`     INT,
  `subscriberid` BIGINT       NOT NULL,
  `venueid`      INT          NOT NULL,
  `ssid`         INT          NOT NULL,
  `apmac`        VARCHAR(255) NOT NULL,
  `mac`          VARCHAR(255) NOT NULL,
  `device`       VARCHAR(255),
  `browser`      VARCHAR(255),
  `os`           VARCHAR(255),
  `age`          INT,
  `gender`       VARCHAR(255) DEFAULT NULL,
  `auth_status`  BOOLEAN      DEFAULT FALSE,
  `created`      TIMESTAMP    DEFAULT now()
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_groups` (
  `tenantid`    INT(10),
  `groupid`     INT          NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(255) NOT NULL,
  `description` VARCHAR(255) NOT NULL,
  `created`     TIMESTAMP             DEFAULT now(),
  `updated`     TIMESTAMP             DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (`groupid`),
  UNIQUE KEY (`tenantid`, `name`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_roles` (
  `tenantid`    INT(10),
  `roleid`      INT          NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(255) NOT NULL,
  `type`        VARCHAR(255)          DEFAULT 'custom',
  `description` VARCHAR(255),
  `created`     TIMESTAMP             DEFAULT now(),
  `updated`     TIMESTAMP             DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (`roleid`),
  UNIQUE KEY (`tenantid`, `name`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_policies` (
  `tenantid` INT,
  `policyid` INT           NOT NULL AUTO_INCREMENT,
  `name`     VARCHAR(255)  NOT NULL,
  `type`     VARCHAR(255)           DEFAULT 'general',
  `policy`   VARCHAR(8000) NOT NULL,
  `created`  TIMESTAMP              DEFAULT now(),
  `updated`  TIMESTAMP              DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (`policyid`),
  UNIQUE KEY (`tenantid`, `name`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_subscriber_groups` (
  `subscriberid` BIGINT,
  `groupid`      INT(10),
  UNIQUE KEY (`subscriberid`, `groupid`),
  FOREIGN KEY (subscriberid) REFERENCES wf_subscribers (subscriberid)
    ON DELETE CASCADE,
  FOREIGN KEY (groupid) REFERENCES wf_groups (groupid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_group_policies` (
  `groupid`  INT(10),
  `policyid` INT(10),
  UNIQUE KEY (`groupid`, `policyid`),
  FOREIGN KEY (groupid) REFERENCES wf_groups (groupid)
    ON DELETE CASCADE,
  FOREIGN KEY (policyid) REFERENCES wf_policies (policyid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_role_policies` (
  `roleid`   INT(10),
  `policyid` INT(10),
  UNIQUE KEY (`roleid`, `policyid`),
  FOREIGN KEY (roleid) REFERENCES wf_roles (roleid)
    ON DELETE CASCADE,
  FOREIGN KEY (policyid) REFERENCES wf_policies (policyid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

-- Adding default super admin details ---
INSERT IGNORE INTO vs_roles (tenantid, name, description)
VALUES (1, 'super_admin', 'super admin role');

INSERT IGNORE INTO vs_policies (tenantid, name, policy)
VALUES (1, 'tenant_admin', '{\"name":\"tenant_admin\",\"type\":"system\",\"statements\":[{\"effect\":\"allow\",
\"actions\":[\"**.**\"],\"resources\":[\"**.**\"]}]}');

INSERT IGNORE INTO vs_role_policies (roleid, policyid)
VALUES (1, 1);

INSERT IGNORE INTO vs_user_roles (userid, roleid)
VALUES (1, 1);

-- adding default subscriber groups and policies --
INSERT IGNORE INTO wf_groups (tenantid, name, description)
VALUES (1, 'default', 'default subscriber group');

INSERT IGNORE INTO wf_roles (tenantid, name, description)
VALUES (1, 'default', 'default subscriber role');

INSERT IGNORE INTO wf_policies (tenantid, name, policy)
VALUES (1, 'default',
        '{\r\n  \"name\": \"default\",\r\n\"bandwidth\": {\r\n\"committedUpLink\": 234,\r\n\"committedDownLink\": 345,\r\n \"reducedDownLink\": 234,\r\n    \"reducedUpLink\": 3456\r\n  },\r\n  \"dataLimit\": {\r\n    \"tx\": 50000,\r\n    \"rx\": 50000\r\n  },\r\n  \"schedule\": {\r\n    \"startDate\": \"2016-10-01\",\r\n\"stopDate\": \"2016-11-30\",\r\n\"accessBy\": {\r\n \"dayOfWeek\": [\r\n1,\r\n 2,\r\n 3,\r\n  4,\r\n 5,\r\n  6,\r\n        7\r\n      ],\r\n      \"time\": {\r\n        \"startTime\": \"01h00h00\",\r\n        \"stopTime\": \"22h00m00s\"\r\n      }\r\n    },\r\n    \"accessDuration\": 3345\r\n  },\r\n  \"network\": {\r\n    \"ssids\": {\r\n      \"allow\": [\r\n        \"*\", \"office\"\r\n      ],\r\n      \"deny\": [\r\n        \"ssid3\"\r\n      ]\r\n    },\r\n    \"apmacs\": {\r\n      \"allow\": [\r\n        \"*\",\"11:22:33:44:55:66\"\r\n ],\r\n  \"deny\": [\r\n \"ssid3\"\r\n ]\r\n }\r\n  },\r\n  \"deviceLimit\": 2,\r\n \"sessionLimit\": 10,\r\n \"sessionTime\" : 30\r\n}');

INSERT IGNORE INTO wf_role_policies (roleid, policyid)
VALUES (1, 1);

INSERT IGNORE INTO wf_group_policies (groupid, policyid)
VALUES (1, 1);


CREATE TABLE IF NOT EXISTS `vs_modules` (
  `module_id` int(11) NOT NULL AUTO_INCREMENT,
  `module_name` varchar(100) NOT NULL,
  `description` varchar(500) DEFAULT NULL,
  `link` varchar(500) DEFAULT NULL,
  `icon` varchar(500) DEFAULT NULL,
  `btnid` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `vs_tenant_modules` (
  `tenant_id` int(11) NOT NULL,
  `module_id` int(11) NOT NULL,
  PRIMARY KEY (`tenant_id`,`module_id`),
  FOREIGN KEY (`tenant_id`) REFERENCES `vs_tenants` (`tenantid`) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (`module_id`) REFERENCES `vs_modules` (`module_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;