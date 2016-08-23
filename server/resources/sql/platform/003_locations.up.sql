
--
-- Location
--
CREATE TABLE IF NOT EXISTS `wf_zones` (
  `userid`      INT,
  `zoneid`      INT          NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(255) NOT NULL,
  `description` VARCHAR(255),
  PRIMARY KEY (zoneid),
  FOREIGN KEY (userid) REFERENCES vs_users (userid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_venues` (
  `userid`      INT,
  `venueid`     INT          NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(255) NOT NULL,
  `address`     VARCHAR(255),
  `zipcode`     VARCHAR(255),
  `telephone`   VARCHAR(255),
  `email`       VARCHAR(255),
  `longitude`   FLOAT,
  `latitude`    FLOAT,
  `description` VARCHAR(255),
  PRIMARY KEY (venueid),
  FOREIGN KEY (userid) REFERENCES vs_users (userid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_zone_venues` (
  `zoneid`  INT,
  `venueid` INT,
  FOREIGN KEY (zoneid) REFERENCES wf_zones (zoneid)
    ON DELETE CASCADE,
  FOREIGN KEY (venueid) REFERENCES wf_venues (venueid)
    ON DELETE CASCADE,
  PRIMARY KEY (zoneid, venueid)
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_aps` (
  `userid`      INT,
  `apid`        INT NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(255) NOT NULL,
  `mac`         VARCHAR(255) UNIQUE,
  `type`        VARCHAR(255),
  `description` VARCHAR(255),
  PRIMARY KEY (apid),
  FOREIGN KEY (userid) REFERENCES vs_users (userid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_venue_aps` (
  `venueid` INT,
  `apid`    INT,
  FOREIGN KEY (venueid) REFERENCES wf_venues (venueid)
    ON DELETE CASCADE,
  FOREIGN KEY (apid) REFERENCES wf_aps (apid)
    ON DELETE CASCADE,
  PRIMARY KEY (venueid, apid)
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `wf_ssids` (
  `apid`        INT,
  `ssid`        VARCHAR(255) NOT NULL,
  `description` VARCHAR(255),
  FOREIGN KEY (apid) REFERENCES wf_aps (apid)
    ON DELETE CASCADE,
  PRIMARY KEY (`apid`, `ssid`)
)
  ENGINE = InnoDB;