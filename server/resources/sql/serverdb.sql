SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";

-- --------------------------------------------------------
--
-- Table structures for dashboard
--
CREATE TABLE IF NOT EXISTS `tenants` (
  `tenantid`  INT NOT NULL AUTO_INCREMENT,
  `domain`    VARCHAR(255) DEFAULT NULL UNIQUE,
  `status`    VARCHAR(255) DEFAULT NULL,
  `createdon` TIMESTAMP,
  PRIMARY KEY (`tenantid`)
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `users` (
  `userid`          BIGINT NOT NULL AUTO_INCREMENT,
  `tenantid`        INT,
  `username`        VARCHAR(255)    DEFAULT NULL,
  `password`        VARCHAR(255)    DEFAULT NULL,
  `email`           VARCHAR(255)    DEFAULT NULL,
  `status`          VARCHAR(255)    DEFAULT NULL,
  `lastupdatedtime` TIMESTAMP,
  PRIMARY KEY (`userid`),
  FOREIGN KEY (tenantid) REFERENCES tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `permissions` (
  `permissionid` BIGINT NOT NULL AUTO_INCREMENT,
  `tenantid`     INT,
  `name`         VARCHAR(255)    DEFAULT NULL,
  `action`       VARCHAR(255)    DEFAULT NULL,
  PRIMARY KEY (`permissionid`),
  FOREIGN KEY (tenantid) REFERENCES tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `userpermissions` (
  `permissionid` BIGINT,
  `userid`       BIGINT,
  PRIMARY KEY (`permissionid`, `userid`),
  FOREIGN KEY (userid) REFERENCES users (userid)
    ON DELETE CASCADE,
  FOREIGN KEY (permissionid) REFERENCES permissions (permissionid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;


