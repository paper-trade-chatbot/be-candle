
-- +migrate Up
CREATE TABLE IF NOT EXISTS `be-candle`.`candle`
(
    `product_id` INTEGER UNSIGNED NOT NULL COMMENT '產品id',
    `interval_type` TINYINT(4) UNSIGNED NOT NULL COMMENT '區間種類 21:1MI, 22: 2MI, 25:5MI, 210:10MI, 215:15MI, 230:30MI, 31:1HR, 41: 1DY, 45: 5DY, 51:1WK, 61:1MO, 71:1YR',
    `start` TIMESTAMP NOT NULL COMMENT '開始時間',
    `open` DOUBLE(20,10) NOT NULL COMMENT '開盤價',
    `close` DOUBLE(20,10) NOT NULL COMMENT '收盤價',
    `high` DOUBLE(20,10) NOT NULL COMMENT '最高價',
    `low` DOUBLE(20,10) NOT NULL COMMENT '最低價',
    `quantity` DOUBLE(20,10) NOT NULL COMMENT '成交量',
    PRIMARY KEY (`product_id`, `interval_type`, `start`)
) DEFAULT CHARSET=`utf8mb4` COLLATE=`utf8mb4_general_ci` COMMENT 'k線';


-- +migrate Down
SET FOREIGN_KEY_CHECKS=0;
DROP TABLE IF EXISTS `candle`;