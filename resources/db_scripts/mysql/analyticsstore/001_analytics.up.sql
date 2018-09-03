--
-- Analytics
--


CREATE TABLE IF NOT EXISTS `vs_apps` (
  `tenantid`    INT,
  `appid`       INT          NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(255) NOT NULL,
  `createon`    TIMESTAMP,
  `appicon`    VARCHAR(255),
  PRIMARY KEY (appid),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_app_venues` (
  `appid`  INT,
  `venueid` INT,
  FOREIGN KEY (appid) REFERENCES vs_apps (appid)
    ON DELETE CASCADE,
  FOREIGN KEY (venueid) REFERENCES wf_venues (venueid)
    ON DELETE CASCADE,
  PRIMARY KEY (appid, venueid)
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_apps_ssids` (
  `appid`  INT,
  `ssid`   INT,
  FOREIGN KEY (appid) REFERENCES vs_apps (appid)
    ON DELETE CASCADE,
  FOREIGN KEY (ssid) REFERENCES wf_ssids (ssid)
    ON DELETE CASCADE,
  PRIMARY KEY (appid, ssid)
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_apps_aps` (
  `appid`  INT,
  `mac` VARCHAR(225),
  FOREIGN KEY (appid) REFERENCES vs_apps (appid)
    ON DELETE CASCADE,
  FOREIGN KEY (mac) REFERENCES wf_aps (mac)
    ON DELETE CASCADE,
  PRIMARY KEY (appid, mac)
)
  ENGINE = InnoDB;


CREATE TABLE IF NOT EXISTS `vs_apps_users` (
  `appid`  INT(11),
  `userid` INT(11),
  FOREIGN KEY (appid) REFERENCES vs_apps (appid)
    ON DELETE CASCADE,
  FOREIGN KEY (userid) REFERENCES vs_users (userid)
    ON DELETE CASCADE,
  PRIMARY KEY (appid, userid)
)
  ENGINE = InnoDB;