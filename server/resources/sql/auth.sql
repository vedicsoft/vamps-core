SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";

-- --------------------------------------------------------
--
-- Table structures for dashboard
--
CREATE TABLE IF NOT EXISTS `vs_tenants` (
  `tenantid`  INT          NOT NULL AUTO_INCREMENT,
  `domain`    VARCHAR(255) NOT NULL UNIQUE,
  `status`    VARCHAR(255)          DEFAULT 'active',
  `createdon` TIMESTAMP,
  PRIMARY KEY (`tenantid`)
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_users` (
  `userid`          BIGINT NOT NULL AUTO_INCREMENT,
  `tenantid`        INT,
  `username`        VARCHAR(255)    DEFAULT NULL,
  `password`        VARCHAR(255)    DEFAULT NULL,
  `email`           VARCHAR(255)    DEFAULT NULL,
  `status`          VARCHAR(255)    DEFAULT NULL,
  `lastupdatedtime` TIMESTAMP,
  PRIMARY KEY (`userid`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_permissions` (
  `permissionid` BIGINT NOT NULL AUTO_INCREMENT,
  `tenantid`     INT,
  `name`         VARCHAR(255)    DEFAULT NULL,
  `action`       VARCHAR(255)    DEFAULT NULL,
  PRIMARY KEY (`permissionid`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_user_permissions` (
  `permissionid` BIGINT,
  `userid`       BIGINT,
  PRIMARY KEY (`permissionid`, `userid`),
  FOREIGN KEY (userid) REFERENCES vs_users (userid)
    ON DELETE CASCADE,
  FOREIGN KEY (permissionid) REFERENCES vs_permissions (permissionid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;


-- provisioned methods : socialauth, pinauth, emailauth
CREATE TABLE IF NOT EXISTS `wf_users` (
  `tenantid`              INT(10),
  `userid`                INT          NOT NULL,
  `uuid`                  CHAR(16) BINARY NOT NULL,
  `email`                 VARCHAR(255) DEFAULT NULL,
  `mobilenumber`          VARCHAR(255) DEFAULT NULL,
  `gender`                VARCHAR(255) DEFAULT NULL,
  `agegroup`              VARCHAR(255) DEFAULT NULL,
  `creationdate`          TIMESTAMP,
  `numberofactivedevices` INT(10),
  PRIMARY KEY (`userid`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_user_devices` (
  `userid`          INT       NOT NULL,
  `mac`             VARCHAR(255) NOT NULL UNIQUE,
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
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_groups` (
  `tenantid`  INT(10),
  `groupid`   INT(10),
  `groupname` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`groupid`),
  UNIQUE KEY (`tenantid`, `groupid`, `groupname`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_policies` (
  `tenantid`   INT(10),
  `policyid`   INT(10),
  `policyname` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`policyid`),
  UNIQUE KEY (`tenantid`, `policyid`, `policyname`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

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
  UNIQUE KEY (`pin`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE,
  FOREIGN KEY (createdby) REFERENCES wf_users (userid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;


CREATE TABLE IF NOT EXISTS `wf_user_groups` (
  `userid`  INT,
  `groupid` INT(10),
  UNIQUE KEY (`userid`, `groupid`),
  FOREIGN KEY (userid) REFERENCES wf_users (userid)
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

CREATE TABLE IF NOT EXISTS `wf_pin_policies` (
  `pinid`    INT(10),
  `policyid` INT(10),
  UNIQUE KEY (`pinid`, `policyid`),
  FOREIGN KEY (pinid) REFERENCES wf_pins (pinid)
    ON DELETE CASCADE,
  FOREIGN KEY (policyid) REFERENCES wf_policies (policyid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

--
-- Usage
--
CREATE TABLE IF NOT EXISTS `wf_daily_usage` (
  `userid`       INT NOT NULL,
  `date`         DATE   NOT NULL,
  `inputoctets`  BIGINT(20) DEFAULT 0,
  `outputoctets` BIGINT(20) DEFAULT 0
)