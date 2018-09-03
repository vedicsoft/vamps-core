SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";

-- --------------------------------------------------------
--
-- Table structures for userstore
--
CREATE TABLE IF NOT EXISTS `vs_tenants` (
  `tenantid`     INT          NOT NULL AUTO_INCREMENT,
  `organization` VARCHAR(255)          DEFAULT NULL,
  `domain`       VARCHAR(255) NOT NULL UNIQUE,
  `status`       VARCHAR(255)          DEFAULT 'active',
  `email`        VARCHAR(255)          DEFAULT NULL,
  `contactno`    VARCHAR(255) ,
  `address`      VARCHAR(255)          DEFAULT NULL,
  `country`      VARCHAR(255)          DEFAULT NULL,
  `created`      TIMESTAMP             DEFAULT now(),
  `updated`      TIMESTAMP             DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (`tenantid`)
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_users` (
  `userid`          INT NOT NULL AUTO_INCREMENT,
  `tenantid`        INT,
  `username`        VARCHAR(255) DEFAULT NULL,
  `password`        VARCHAR(255) DEFAULT NULL,
  `confirmpassword` VARCHAR(255) DEFAULT NULL,
  `email`           VARCHAR(255) DEFAULT NULL,
  `contactno`       VARCHAR(255) DEFAULT NULL,
  `status`          VARCHAR(255) DEFAULT NULL,
  `created`         TIMESTAMP    DEFAULT now(),
  `updated`         TIMESTAMP    DEFAULT now() ON UPDATE now(),
  `privilege`       int(11)      DEFAULT '1',
  PRIMARY KEY (`userid`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_groups` (
  `groupid`  INT          NOT NULL AUTO_INCREMENT,
  `tenantid` INT,
  `name`     VARCHAR(255) NOT NULL,
  `created`  TIMESTAMP             DEFAULT now(),
  `updated`  TIMESTAMP             DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (`groupid`),
  UNIQUE (`tenantid`, `name`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_policies` (
  `policyid` INT           NOT NULL AUTO_INCREMENT,
  `tenantid` INT,
  `name`     VARCHAR(255)  NOT NULL,
  `type`     VARCHAR(255)  DEFAULT 'general',
  `policy`   VARCHAR(9000) NOT NULL,
  `created`  TIMESTAMP              DEFAULT now(),
  `updated`  TIMESTAMP              DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (`policyid`),
  UNIQUE (`tenantid`, `name`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_roles` (
  `tenantid`    INT(10),
  `roleid`      INT          NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(255) NOT NULL,
  `type`        VARCHAR(255) DEFAULT 'custom',
  `description` VARCHAR(255),
  `created`     TIMESTAMP             DEFAULT now(),
  `updated`     TIMESTAMP             DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (`roleid`),
  UNIQUE KEY (`tenantid`, `name`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_role_policies` (
  `roleid`   INT,
  `policyid` INT,
  FOREIGN KEY (roleid) REFERENCES vs_roles (roleid)
    ON DELETE CASCADE,
  FOREIGN KEY (policyid) REFERENCES vs_policies (policyid)
    ON DELETE CASCADE,
  PRIMARY KEY (`policyid`, `roleid`)
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_user_roles` (
  `userid` INT,
  `roleid` INT,
  PRIMARY KEY (`userid`, `roleid`),
  FOREIGN KEY (userid) REFERENCES vs_users (userid)
    ON DELETE CASCADE,
  FOREIGN KEY (roleid) REFERENCES vs_roles (roleid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_user_groups` (
  `userid`  INT,
  `groupid` INT,
  PRIMARY KEY (`userid`, `groupid`),
  FOREIGN KEY (userid) REFERENCES vs_users (userid)
    ON DELETE CASCADE,
  FOREIGN KEY (groupid) REFERENCES vs_groups (groupid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_user_policies` (
  `userid`   INT,
  `policyid` INT,
  PRIMARY KEY (`userid`, `policyid`),
  FOREIGN KEY (userid) REFERENCES vs_users (userid)
    ON DELETE CASCADE,
  FOREIGN KEY (policyid) REFERENCES vs_policies (policyid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

--
-- initial data set
--
-- Inserting system admin data data set
INSERT IGNORE INTO vs_tenants (domain, status)
VALUES ('super.com', 'active');

INSERT IGNORE INTO vs_users (tenantid, username, password, email, status)
VALUES (1, 'admin', '$2a$10$FesfnIBKqhH2MuF1hmss0umXNrrx28AW1E4re9OCAwib3cIOKBz3C', 'admin@vedicsoft.com', 'active');
