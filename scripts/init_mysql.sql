-- Create a database
CREATE DATABASE IF NOT EXISTS `gin-admin` DEFAULT CHARACTER SET = `utf8mb4`;
-- 初始化数据库 username root password root
INSERT INTO `g_user` (`created_at`,`updated_at`,`user_name`,`real_name`,`password`,`email`,`phone`,`status`,`creator`,`id`) VALUES ('2022-01-06 14:56:55.542','2022-01-06 14:56:55.542','root','root','dc76e9f0c0006e8f919e0c515c66dbba3982f785','z@gmail.com','1234578932326',1,9,23524702785896783)