USE `emm_moneybox`;

-- -------------------
-- Create table `cash`
-- -------------------
DROP TABLE IF EXISTS cash_flow;
CREATE TABLE `cash_flow`
(
    `id`           VARCHAR(24)  NOT NULL,
    `category_id`  VARCHAR(24)  NOT NULL,
    `belongs_date` TIMESTAMP    NOT NULL,
    `flow_type`    VARCHAR(10)  NOT NULL COMMENT 'INCOME/OUTCOME',
    `amount`       DECIMAL      NOT NULL,
    `description`  VARCHAR(200) NOT NULL,
    `remark`       VARCHAR(200)          DEFAULT NULL COMMENT 'KEEP EMPTY',
    `create_time`  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    `modify_time`  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = UTF8MB4
    COMMENT ='Cash Flow Table';

CREATE INDEX cash_flow_category_id_index ON cash_flow (category_id);
CREATE INDEX cash_flow_belongs_date_index ON cash_flow (belongs_date);
CREATE INDEX cash_flow_flow_type_index ON cash_flow (flow_type);
