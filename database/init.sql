SET NAMES utf8mb4;

CREATE DATABASE IF NOT EXISTS gomall_lite
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_0900_ai_ci;

USE gomall_lite;

CREATE TABLE IF NOT EXISTS users (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  username varchar(64) NOT NULL,
  password_hash varchar(255) NOT NULL,
  nickname varchar(64) DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_users_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS products (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(128) NOT NULL,
  price bigint NOT NULL,
  stock bigint NOT NULL,
  image varchar(255) DEFAULT NULL,
  category varchar(64) DEFAULT NULL,
  description text,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_products_category (category),
  KEY idx_products_category_price (category, price),
  KEY idx_products_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS cart_items (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id bigint unsigned NOT NULL,
  product_id bigint unsigned NOT NULL,
  count bigint NOT NULL,
  checked tinyint(1) NOT NULL DEFAULT 1,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_cart_items_user_id (user_id),
  KEY idx_cart_items_product_id (product_id),
  UNIQUE KEY uk_cart_items_user_product (user_id, product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS addresses (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id bigint unsigned NOT NULL,
  receiver varchar(64) NOT NULL,
  phone varchar(32) NOT NULL,
  province varchar(64) NOT NULL,
  city varchar(64) NOT NULL,
  district varchar(64) NOT NULL,
  detail varchar(255) NOT NULL,
  is_default tinyint(1) NOT NULL DEFAULT 0,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_addresses_user_id (user_id),
  KEY idx_addresses_user_default (user_id, is_default)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS orders (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id bigint unsigned NOT NULL,
  order_no varchar(64) NOT NULL,
  total_amount bigint NOT NULL,
  status varchar(32) NOT NULL,
  address_snapshot text,
  remark varchar(255) DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_orders_order_no (order_no),
  KEY idx_orders_user_id (user_id),
  KEY idx_orders_user_status_created (user_id, status, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS order_items (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  order_id bigint unsigned NOT NULL,
  product_id bigint unsigned DEFAULT NULL,
  product_name varchar(128) NOT NULL,
  product_image varchar(255) DEFAULT NULL,
  price bigint NOT NULL,
  count bigint NOT NULL,
  subtotal bigint NOT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_order_items_order_id (order_id),
  KEY idx_order_items_product_id (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO users (id, username, password_hash, nickname, created_at, updated_at)
VALUES (1, 'admin', '$2a$10$VQev2ckoMLaTO6tHrCSz9eIN/2vXXZqbLEAeAtagnblfMwzORwOvm', '演示管理员', NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE updated_at = VALUES(updated_at);

INSERT INTO users (id, username, password_hash, nickname, created_at, updated_at)
WITH RECURSIVE seq(n) AS (
  SELECT 2
  UNION ALL
  SELECT n + 1 FROM seq WHERE n < 30
)
SELECT
  n,
  CONCAT('demo_user_', LPAD(n - 1, 2, '0')),
  '$2a$10$VQev2ckoMLaTO6tHrCSz9eIN/2vXXZqbLEAeAtagnblfMwzORwOvm',
  CONCAT('演示用户', LPAD(n - 1, 2, '0')),
  NOW(3),
  NOW(3)
FROM seq
ON DUPLICATE KEY UPDATE nickname = VALUES(nickname), updated_at = VALUES(updated_at);

INSERT INTO products (id, name, price, stock, image, category, description, created_at, updated_at)
WITH RECURSIVE seq(n) AS (
  SELECT 1
  UNION ALL
  SELECT n + 1 FROM seq WHERE n < 240
)
SELECT
  n,
  CASE ((n - 1) MOD 6)
    WHEN 0 THEN CONCAT('数码精选 ', LPAD(n, 3, '0'))
    WHEN 1 THEN CONCAT('办公装备 ', LPAD(n, 3, '0'))
    WHEN 2 THEN CONCAT('家电好物 ', LPAD(n, 3, '0'))
    WHEN 3 THEN CONCAT('穿搭单品 ', LPAD(n, 3, '0'))
    WHEN 4 THEN CONCAT('个护精选 ', LPAD(n, 3, '0'))
    ELSE CONCAT('食品礼包 ', LPAD(n, 3, '0'))
  END,
  59 + ((n * 37) MOD 9000),
  20 + ((n * 11) MOD 260),
  CASE ((n - 1) MOD 6)
    WHEN 0 THEN CASE (n MOD 5)
      WHEN 0 THEN '/images/iphone.svg'
      WHEN 1 THEN '/images/mi14.svg'
      WHEN 2 THEN '/images/headphone.svg'
      WHEN 3 THEN '/images/airpods.svg'
      ELSE '/images/watch.svg'
    END
    WHEN 1 THEN '/images/macbook.svg'
    WHEN 2 THEN '/images/airfryer.svg'
    WHEN 3 THEN '/images/shoe.svg'
    WHEN 4 THEN '/images/sk2.svg'
    ELSE '/images/snack.svg'
  END,
  CASE ((n - 1) MOD 6)
    WHEN 0 THEN '手机数码'
    WHEN 1 THEN '电脑办公'
    WHEN 2 THEN '家用电器'
    WHEN 3 THEN '服饰鞋包'
    WHEN 4 THEN '美妆个护'
    ELSE '食品生鲜'
  END,
  CASE ((n - 1) MOD 6)
    WHEN 0 THEN '适合日常学习、办公、娱乐和通勤使用的热门数码单品。'
    WHEN 1 THEN '覆盖移动办公、内容创作和高效学习场景的电脑办公产品。'
    WHEN 2 THEN '面向厨房、清洁和居家生活的小型家用电器。'
    WHEN 3 THEN '适合通勤、运动和休闲场景的服饰鞋包商品。'
    WHEN 4 THEN '覆盖基础护肤、清洁护理和日常个护的展示商品。'
    ELSE '适合家庭休闲、办公室分享和节日囤货的食品生鲜商品。'
  END,
  NOW(3),
  NOW(3)
FROM seq
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  price = VALUES(price),
  stock = VALUES(stock),
  image = VALUES(image),
  category = VALUES(category),
  description = VALUES(description),
  updated_at = VALUES(updated_at);

UPDATE products SET name = 'iPhone 15 Pro 256GB', price = 7999, stock = 100, image = '/images/iphone.svg', category = '手机数码', description = '搭载 A17 Pro 芯片，钛金属边框，适合日常学习、办公和娱乐。', updated_at = NOW(3) WHERE id = 1;
UPDATE products SET name = 'MacBook Air M2', price = 8999, stock = 50, image = '/images/macbook.svg', category = '电脑办公', description = '轻薄便携的 M2 芯片笔记本，适合开发、办公、学习和内容创作。', updated_at = NOW(3) WHERE id = 2;
UPDATE products SET name = '索尼 WH-1000XM5', price = 2599, stock = 80, image = '/images/headphone.svg', category = '手机数码', description = '旗舰级主动降噪耳机，拥有舒适佩戴体验和优秀音质表现。', updated_at = NOW(3) WHERE id = 3;
UPDATE products SET name = '小米 14', price = 4299, stock = 120, image = '/images/mi14.svg', category = '手机数码', description = '高性能安卓旗舰手机，轻薄机身，日常拍照和游戏体验优秀。', updated_at = NOW(3) WHERE id = 4;
UPDATE products SET name = 'SK-II 神仙水 230ml', price = 1230, stock = 90, image = '/images/sk2.svg', category = '美妆个护', description = '经典护肤精华水，适合日常护肤使用。', updated_at = NOW(3) WHERE id = 5;
UPDATE products SET name = 'Nike Air Force 1', price = 699, stock = 60, image = '/images/shoe.svg', category = '服饰鞋包', description = '经典百搭休闲鞋，适合多场景穿搭。', updated_at = NOW(3) WHERE id = 6;
UPDATE products SET name = 'Apple Watch Series 9', price = 2999, stock = 70, image = '/images/watch.svg', category = '手机数码', description = '健康记录、运动追踪和消息提醒的智能手表。', updated_at = NOW(3) WHERE id = 7;
UPDATE products SET name = 'AirPods Pro 第二代', price = 1899, stock = 85, image = '/images/airpods.svg', category = '手机数码', description = '主动降噪真无线耳机，便携易用。', updated_at = NOW(3) WHERE id = 8;
UPDATE products SET name = '三只松鼠零食礼包', price = 129, stock = 200, image = '/images/snack.svg', category = '食品生鲜', description = '多口味零食组合，适合家庭休闲和办公室分享。', updated_at = NOW(3) WHERE id = 9;
UPDATE products SET name = '美的空气炸锅', price = 399, stock = 45, image = '/images/airfryer.svg', category = '家用电器', description = '便捷厨房小家电，轻松制作低油美食。', updated_at = NOW(3) WHERE id = 10;

INSERT INTO addresses (id, user_id, receiver, phone, province, city, district, detail, is_default, created_at, updated_at)
WITH RECURSIVE seq(n) AS (
  SELECT 1
  UNION ALL
  SELECT n + 1 FROM seq WHERE n < 30
)
SELECT
  1000 + n,
  n,
  IF(n = 1, '张三', CONCAT('演示用户', LPAD(n - 1, 2, '0'))),
  CONCAT('138', LPAD(10000000 + n, 8, '0')),
  '北京市',
  '北京市',
  CASE (n MOD 4) WHEN 0 THEN '朝阳区' WHEN 1 THEN '海淀区' WHEN 2 THEN '西城区' ELSE '丰台区' END,
  CONCAT('演示路 ', 80 + n, ' 号 ', 1 + (n MOD 20), ' 单元'),
  1,
  NOW(3),
  NOW(3)
FROM seq
ON DUPLICATE KEY UPDATE
  receiver = VALUES(receiver),
  phone = VALUES(phone),
  province = VALUES(province),
  city = VALUES(city),
  district = VALUES(district),
  detail = VALUES(detail),
  is_default = VALUES(is_default),
  updated_at = VALUES(updated_at);

INSERT INTO cart_items (id, user_id, product_id, count, checked, created_at, updated_at)
WITH RECURSIVE seq(n) AS (
  SELECT 1
  UNION ALL
  SELECT n + 1 FROM seq WHERE n < 20
)
SELECT
  1000 + n,
  1 + (n MOD 10),
  1 + ((n * 5) MOD 240),
  1 + (n MOD 3),
  IF(n MOD 4 = 0, 0, 1),
  NOW(3),
  NOW(3)
FROM seq
ON DUPLICATE KEY UPDATE count = VALUES(count), checked = VALUES(checked), updated_at = VALUES(updated_at);

INSERT INTO orders (id, user_id, order_no, total_amount, status, address_snapshot, remark, created_at, updated_at)
WITH RECURSIVE seq(n) AS (
  SELECT 1
  UNION ALL
  SELECT n + 1 FROM seq WHERE n < 40
)
SELECT
  1000 + n,
  1 + ((n - 1) MOD 20),
  CONCAT('DEMO', DATE_FORMAT(NOW(), '%Y%m%d'), LPAD(n, 5, '0')),
  0,
  CASE (n MOD 4) WHEN 0 THEN '待支付' WHEN 1 THEN '已支付' WHEN 2 THEN '已完成' ELSE '已取消' END,
  CONCAT('演示用户 ', 1 + ((n - 1) MOD 20), ' 13800000000 北京市 北京市 朝阳区 演示路 ', n, ' 号'),
  IF(n MOD 5 = 0, '请尽快发货', ''),
  DATE_SUB(NOW(3), INTERVAL n DAY),
  NOW(3)
FROM seq
ON DUPLICATE KEY UPDATE status = VALUES(status), remark = VALUES(remark), updated_at = VALUES(updated_at);

INSERT INTO order_items (id, order_id, product_id, product_name, product_image, price, count, subtotal, created_at, updated_at)
WITH RECURSIVE seq(n) AS (
  SELECT 1
  UNION ALL
  SELECT n + 1 FROM seq WHERE n < 40
)
SELECT
  2000 + n * 2 - 1,
  1000 + n,
  p.id,
  p.name,
  p.image,
  p.price,
  1 + (n MOD 3),
  p.price * (1 + (n MOD 3)),
  NOW(3),
  NOW(3)
FROM seq
JOIN products p ON p.id = 1 + ((n * 7) MOD 240)
ON DUPLICATE KEY UPDATE
  product_name = VALUES(product_name),
  product_image = VALUES(product_image),
  price = VALUES(price),
  count = VALUES(count),
  subtotal = VALUES(subtotal),
  updated_at = VALUES(updated_at);

INSERT INTO order_items (id, order_id, product_id, product_name, product_image, price, count, subtotal, created_at, updated_at)
WITH RECURSIVE seq(n) AS (
  SELECT 1
  UNION ALL
  SELECT n + 1 FROM seq WHERE n < 40
)
SELECT
  2000 + n * 2,
  1000 + n,
  p.id,
  p.name,
  p.image,
  p.price,
  1 + (n MOD 2),
  p.price * (1 + (n MOD 2)),
  NOW(3),
  NOW(3)
FROM seq
JOIN products p ON p.id = 1 + ((n * 13) MOD 240)
ON DUPLICATE KEY UPDATE
  product_name = VALUES(product_name),
  product_image = VALUES(product_image),
  price = VALUES(price),
  count = VALUES(count),
  subtotal = VALUES(subtotal),
  updated_at = VALUES(updated_at);

UPDATE orders o
JOIN (
  SELECT order_id, SUM(subtotal) AS total_amount
  FROM order_items
  WHERE order_id BETWEEN 1001 AND 1040
  GROUP BY order_id
) s ON s.order_id = o.id
SET o.total_amount = s.total_amount,
    o.updated_at = NOW(3);
