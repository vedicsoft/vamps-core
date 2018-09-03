--
-- captive portals
--

CREATE TABLE IF NOT EXISTS `vs_portals` (
  `tenantid`    INT,
  `portalid`    INT          NOT NULL AUTO_INCREMENT,
  `portalname`  VARCHAR(255) NOT NULL,
  `type`        VARCHAR(255),
  `category`    VARCHAR(255),
  `description` VARCHAR(255)          DEFAULT NULL,
  `portalicon`  VARCHAR(255)          DEFAULT NULL,
  `prod_status` VARCHAR(255)          DEFAULT NULL,
  `prod_url`    VARCHAR(255)          DEFAULT NULL,
  `created`     TIMESTAMP             DEFAULT now(),
  `updated`     TIMESTAMP             DEFAULT now() ON UPDATE now(),

  PRIMARY KEY (portalid),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vs_portal_users` (
  `portalid` INT,
  `userid`   INT,
  FOREIGN KEY (portalid) REFERENCES vs_portals (portalid)
    ON DELETE CASCADE,
  FOREIGN KEY (userid) REFERENCES vs_users (userid)
    ON DELETE CASCADE,
  PRIMARY KEY (portalid, userid)
)
  ENGINE = InnoDB;


CREATE TABLE IF NOT EXISTS `vs_portal_templates` (
  `tenantid`    INT,
  `templateid`  INT          NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(255) NOT NULL,
  `category`    VARCHAR(255),
  `type`        VARCHAR(255),
  `templateURL` VARCHAR(255),
  `description` VARCHAR(255),
  `created`     TIMESTAMP             DEFAULT now(),
  `updated`     TIMESTAMP             DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (templateid),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE
)
  ENGINE = InnoDB;

INSERT INTO `vs_portal_templates` (`tenantid`, `templateid`, `name`, `category`)
VALUES
  (1, 2, 'Retail', 'retail demo'),
  (1, 1, 'Public', 'public demo'),
  (1, 3, 'Hospitality', 'Hospitality demo');