CREATE USER IF NOT EXISTS 'admin'@'%' IDENTIFIED WITH caching_sha2_password BY 'admin123';
CREATE DATABASE IF NOT EXISTS `futures_trading` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL ON futures_trading.* TO 'admin'@'%';
