
-- --------------------------------------------------------
--
-- Table structures for dashboard
--
CREATE TABLE IF NOT EXISTS `vs_tenants` (
  `tenantid`  INT PRIMARY KEY,
  `domain`    VARCHAR(255)  UNIQUE,
  `status`    VARCHAR(255)          DEFAULT 'active',
  `createdon` TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `vs_users` (
  `userid`          INT PRIMARY KEY ,
  `tenantid`        INT,
  `username`        VARCHAR(255) DEFAULT NULL,
  `password`        VARCHAR(255) DEFAULT NULL,
  `email`           VARCHAR(255) DEFAULT NULL,
  `status`          VARCHAR(255) DEFAULT NULL,
  `lastupdatedtime` TIMESTAMP,
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `vs_permissions` (
  `permissionid` INT PRIMARY KEY ,
  `tenantid`     INT,
  `name`         VARCHAR(255) DEFAULT NULL,
  `action`       VARCHAR(255) DEFAULT NULL,
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `vs_user_permissions` (
  `permissionid` INT,
  `userid`       INT,
  PRIMARY KEY (`permissionid`, `userid`),
  FOREIGN KEY (userid) REFERENCES vs_users (userid)
    ON DELETE CASCADE,
  FOREIGN KEY (permissionid) REFERENCES vs_permissions (permissionid)
    ON DELETE CASCADE
);

-- provisioned methods : socialauth, pinauth, emailauth
CREATE TABLE IF NOT EXISTS `wf_users` (
  `userid`           BIGINT       ,
  `tenantid`         INT,
  `username`         VARCHAR(255),
  `password`         VARCHAR(255),
  `email`            VARCHAR(255)  UNIQUE,
  `account_status`   VARCHAR(255)          DEFAULT NULL,
  `first_name`       VARCHAR(255)          DEFAULT NULL,
  `last_name`        VARCHAR(255)          DEFAULT NULL,
  `gender`           VARCHAR(255)          DEFAULT NULL,
  `birthday`         DATE,
  `age`              INT,
  `age_upper`        INT,
  `age_lower`        INT,
  `religion`         VARCHAR(255)          DEFAULT NULL,
  `occupation`        VARCHAR(255)          DEFAULT NULL,
  `marital_status`   VARCHAR(255)          DEFAULT NULL,
  `profile_image`    VARCHAR(255)          DEFAULT NULL,
  `mobile_number`    VARCHAR(255)          DEFAULT NULL,
  `admin_notes`      VARCHAR(255)          DEFAULT NULL,
  `last_updatedtime` TIMESTAMP,
  PRIMARY KEY (`userid`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `wf_user_devices` (
  `userid`          BIGINT       ,
  `mac`             VARCHAR(255)  UNIQUE,
  `password`        VARCHAR(255) DEFAULT NULL,
  `parentalcontrol` VARCHAR(255) DEFAULT 'OFF',
  `useragent`       VARCHAR(255),
  `browser`         VARCHAR(255),
  `os`              VARCHAR(255),
  `device`          VARCHAR(255),
  `creationdate`    TIMESTAMP,
  PRIMARY KEY (mac),
  FOREIGN KEY (userid) REFERENCES wf_users (userid)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `wf_groups` (
  `tenantid`  INT(10),
  `groupid`   INT(10),
  `groupname` VARCHAR(255) ,
  PRIMARY KEY (`groupid`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `wf_policies` (
  `tenantid`   INT(10),
  `policyid`   INT(10),
  `policyname` VARCHAR(255) ,
  PRIMARY KEY (`policyid`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `wf_pins` (
  `tenantid`          INT(10),
  `pinid`             INT(10),
  `pin`               INT(10),
  `createdby`         INT,
  `devicelimit`       INT,
  `activedevicecount` INT,
  `creationtime`      TIMESTAMP,
  `status`            VARCHAR(255) DEFAULT 'valid',
  PRIMARY KEY (`pinid`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE,
  FOREIGN KEY (createdby) REFERENCES vs_users (userid)
    ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS `wf_user_groups` (
  `userid`  BIGINT,
  `groupid` INT(10),
  FOREIGN KEY (userid) REFERENCES wf_users (userid)
    ON DELETE CASCADE,
  FOREIGN KEY (groupid) REFERENCES wf_groups (groupid)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `wf_group_policies` (
  `groupid`  INT(10),
  `policyid` INT(10),
  FOREIGN KEY (groupid) REFERENCES wf_groups (groupid)
    ON DELETE CASCADE,
  FOREIGN KEY (policyid) REFERENCES wf_policies (policyid)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `wf_pin_policies` (
  `pinid`    INT(10),
  `policyid` INT(10),
  FOREIGN KEY (pinid) REFERENCES wf_pins (pinid)
    ON DELETE CASCADE,
  FOREIGN KEY (policyid) REFERENCES wf_policies (policyid)
    ON DELETE CASCADE
);

--
-- Usage
--
CREATE TABLE IF NOT EXISTS `wf_daily_usage` (
  `userid`       BIGINT,
  `date`         DATE ,
  `inputoctets`  BIGINT(20) DEFAULT 0,
  `outputoctets` BIGINT(20) DEFAULT 0
);

--
-- initial data set
--
-- Inserting default data set
INSERT INTO vs_tenants (tenantid,domain, status)
VALUES (1,'super.com', 'active');

INSERT INTO vs_users (tenantid, username, password, email, status)
VALUES (1, 'admin', '$2a$10$FesfnIBKqhH2MuF1hmss0umXNrrx28AW1E4re9OCAwib3cIOKBz3C', 'admin@vedicsoft.com', 'active');

INSERT INTO vs_permissions (permissionid, tenantid, name, action)
VALUES
  (1, 1, 'wifi_location', 'read'),
  (2, 1, 'wifi_location', 'write'),
  (3, 1, 'wifi_location', 'execute'),
  (4, 1, 'wifi_users', 'read'),
  (5, 1, 'wifi_users', 'write'),
  (6, 1, 'wifi_users', 'execute'),
  (7, 1, 'dashboard_users', 'read'),
  (8, 1, 'dashboard_users', 'write'),
  (9, 1, 'dashboard_users', 'execute');

INSERT INTO vs_user_permissions (userid, permissionid)
VALUES (1, 1),
  (1, 2),
  (1, 3),
  (1, 4),
  (1, 5),
  (1, 6),
  (1, 7),
  (1, 8),
  (1, 9);