# Readme

## 建表SQL
### 1 user表
```sql
CREATE TABLE `user` (
                        `user_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
                        `username` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名',
                        `password` varchar(32) DEFAULT '' COMMENT '密码',
                        `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        PRIMARY KEY (`user_id`),
                        UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
```