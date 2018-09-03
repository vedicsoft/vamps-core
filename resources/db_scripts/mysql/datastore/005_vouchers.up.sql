CREATE TABLE IF NOT EXISTS `wf_batches` (
  `batchid`   BIGINT           NOT NULL AUTO_INCREMENT,
  `tenantid`  INT,
  `batch_num` INT(4) DEFAULT 1 NOT NULL,
  `created_by` INT,
  `created`   TIMESTAMP                 DEFAULT now(),
  `updated`   TIMESTAMP                 DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (`batchid`)
)
  ENGINE = InnoDB;

--
-- voucher status new, activated, deactivated, onhold
--
CREATE TABLE IF NOT EXISTS `wf_vouchers` (
  `voucherid`    BIGINT       NOT NULL AUTO_INCREMENT,
  `tenantid`     INT,
  `policyid`     INT(10),
  `batch_num`    INT(4),
  `pin`          VARCHAR(255) NOT NULL,
  `pan`          VARCHAR(255) NOT NULL,
  `status`       VARCHAR(255)          DEFAULT 'new',
  `activated_by` BIGINT,
  `activated`    TIMESTAMP,
  `created`      TIMESTAMP             DEFAULT now(),
  `updated`      TIMESTAMP             DEFAULT now() ON UPDATE now(),
  PRIMARY KEY (`voucherid`)
)
  ENGINE = InnoDB;

create index wf_vouchers_t_pin on wf_vouchers(tenantid, pin);