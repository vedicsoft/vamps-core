--
-- Location
--

CREATE TABLE IF NOT EXISTS `wf_nodes` (
  `tenantid` INT,
  `nodeid`   INT(11)      NOT NULL AUTO_INCREMENT,
  `text`     VARCHAR(200) NOT NULL,
  `type`     VARCHAR(200) NOT NULL DEFAULT 'zone',
  `data`     VARCHAR(200),
  `parentid` INT                   DEFAULT NULL,
  `created`  TIMESTAMP             DEFAULT now(),
  `updated`  TIMESTAMP             DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (nodeid),
  FOREIGN KEY (parentid) REFERENCES wf_nodes (nodeid)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

INSERT INTO `wf_nodes` (`tenantid`, `nodeid`, `text`, `type`) VALUES
  (1, 1, 'Zones', 'zone');

CREATE TABLE IF NOT EXISTS `wf_venues` (
  `tenantid`      INT,
  `venueid`       INT          NOT NULL AUTO_INCREMENT,
  `nodeid`        INT,
  `name`          VARCHAR(255) NOT NULL,
  `street_number` VARCHAR(255),
  `route`         VARCHAR(255),
  `city`          VARCHAR(255),
  `state`         VARCHAR(255),
  `country`       VARCHAR(255),
  `zipcode`       VARCHAR(255),
  `telephone`     VARCHAR(255),
  `email`         VARCHAR(255),
  `longitude`     FLOAT,
  `latitude`      FLOAT,
  `description`   VARCHAR(255),
  `created`       TIMESTAMP             DEFAULT now(),
  `updated`       TIMESTAMP             DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (venueid),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE,
  FOREIGN KEY (nodeid) REFERENCES wf_nodes (nodeid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_ap_vendors` (
  `tenantid`    INT,
  `vendorid`    INT          NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(255) NOT NULL,
  `description` VARCHAR(255),
  `created`     TIMESTAMP             DEFAULT now(),
  `updated`     TIMESTAMP             DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (vendorid),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE,
  UNIQUE KEY `vendor_name`(`tenantid`, `name`)
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_ap_types` (
  `tenantid` INT,
  `typeid`   INT          NOT NULL AUTO_INCREMENT,
  `vendorid` INT,
  `name`     VARCHAR(255) NOT NULL,
  `version`  VARCHAR(255) NOT NULL,
  `created`  TIMESTAMP             DEFAULT now(),
  `updated`  TIMESTAMP             DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (typeid),
  FOREIGN KEY (vendorid) REFERENCES wf_ap_vendors (vendorid),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_aps` (
  `tenantid`    INT,
  `apid`        INT          NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(255) NOT NULL,
  `publicip`    VARCHAR(15),
  `mac`         VARCHAR(255) UNIQUE,
  `typeid`      INT,
  `venueid`     INT,
  `description` VARCHAR(255),
  `created`     TIMESTAMP             DEFAULT now(),
  `updated`     TIMESTAMP             DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (apid),
  FOREIGN KEY (typeid) REFERENCES wf_ap_types (typeid)
    ON DELETE SET NULL,
  FOREIGN KEY (venueid) REFERENCES wf_venues (venueid)
    ON DELETE SET NULL,
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_ssids` (
  `tenantid`    INT,
  `ssid`        INT          NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(255) NOT NULL,
  `description` VARCHAR(255),
  `created`     TIMESTAMP             DEFAULT now(),
  `updated`     TIMESTAMP             DEFAULT now() ON UPDATE now(),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE,
  PRIMARY KEY (`ssid`),
  UNIQUE KEY `ssid_name`(`tenantid`, `name`)
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_ap_ssids` (
  `apid` INT,
  `ssid` INT,
  FOREIGN KEY (apid) REFERENCES wf_aps (apid)
    ON DELETE CASCADE,
  FOREIGN KEY (ssid) REFERENCES wf_ssids (ssid)
    ON DELETE CASCADE,
  PRIMARY KEY (`apid`)
)
  ENGINE = InnoDB;

INSERT IGNORE INTO wf_ap_vendors (tenantid, vendorid, name, description)
VALUES
  (1, 1, 'Ruckus', 'ruckus vendor'),
  (1, 2, 'Zebra', 'Zebra vendor'),
  (1, 3, 'Mikrotik', 'MKT vendor'),
  (1, 4, 'Aruba', 'MKT vendor'),
  (1, 5, 'Cisco', 'MKT vendor');

INSERT IGNORE INTO wf_ap_types (tenantid, vendorid, name, version)
VALUES
  (1, 1, 'Ruckus 456', '1.5.3'),
  (1, 2, 'Zebra 7502', '3.4.2'),
  (1, 3, 'Mikrotik E', '1.2'),
  (1, 4, 'Aruba G400', '2.1'),
  (1, 5, 'Cisco A234', '2.2');

CREATE TABLE IF NOT EXISTS `vs_user_nodes` (
  `userid`  INT,
  `nodeid`  INT,
  PRIMARY KEY (userid, nodeid),
  FOREIGN KEY (userid) REFERENCES vs_users (userid)
    ON DELETE CASCADE,
  FOREIGN KEY (nodeid) REFERENCES wf_nodes (nodeid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;