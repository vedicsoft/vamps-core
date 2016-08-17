SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";

-- --------------------------------------------------------
--
-- Table structures for userstore
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
  `userid`          INT NOT NULL AUTO_INCREMENT,
  `tenantid`        INT,
  `username`        VARCHAR(255) DEFAULT NULL,
  `password`        VARCHAR(255) DEFAULT NULL,
  `email`           VARCHAR(255) DEFAULT NULL,
  `status`          VARCHAR(255) DEFAULT NULL,
  `lastupdatedtime` TIMESTAMP,
  PRIMARY KEY (`userid`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_permissions` (
  `permissionid` INT NOT NULL AUTO_INCREMENT,
  `tenantid`     INT,
  `name`         VARCHAR(255) DEFAULT NULL,
  `action`       VARCHAR(255) DEFAULT NULL,
  PRIMARY KEY (`permissionid`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_user_permissions` (
  `permissionid` INT,
  `userid`       INT,
  PRIMARY KEY (`permissionid`, `userid`),
  FOREIGN KEY (userid) REFERENCES vs_users (userid)
    ON DELETE CASCADE,
  FOREIGN KEY (permissionid) REFERENCES vs_permissions (permissionid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

--
-- initial data set
--
-- Inserting default data set
INSERT IGNORE INTO vs_tenants (domain, status)
VALUES ('super.com', 'active');

INSERT IGNORE INTO vs_users (tenantid, username, password, email, status)
VALUES (1, 'admin', '$2a$10$FesfnIBKqhH2MuF1hmss0umXNrrx28AW1E4re9OCAwib3cIOKBz3C', 'admin@vedicsoft.com', 'active');

INSERT IGNORE INTO vs_permissions (permissionid, tenantid, name, action)
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

INSERT IGNORE INTO vs_user_permissions (userid, permissionid)
VALUES (1, 1),
  (1, 2),
  (1, 3),
  (1, 4),
  (1, 5),
  (1, 6),
  (1, 7),
  (1, 8),
  (1, 9);