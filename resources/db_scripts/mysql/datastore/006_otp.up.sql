--
-- otp verified, unverified
--
CREATE TABLE IF NOT EXISTS `wf_otps` (
  `otpid`               BIGINT       NOT NULL AUTO_INCREMENT,
  `tenantid`            INT,
  `device_mac`          VARCHAR(255) NOT NULL,
  `country_code`        VARCHAR(255) NOT NULL,
  `mobile_number`       VARCHAR(255) NOT NULL,
  `otp`                 VARCHAR(255) NOT NULL,
  `generate_count`      INT                   DEFAULT 1,
  `try_count`           INT                   DEFAULT 0,
  `status`              VARCHAR(255)          DEFAULT 'unverified',
  `last_generated_time` DATETIME     NOT NULL,
  `expiration_time`     DATETIME     NOT NULL,
  `created`             TIMESTAMP             DEFAULT now(),
  `updated`             TIMESTAMP             DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (`otpid`),
  FOREIGN KEY (tenantid) REFERENCES vs_tenants (tenantid)
    ON DELETE CASCADE,
  UNIQUE (tenantid, device_mac)
)
  ENGINE = InnoDB;