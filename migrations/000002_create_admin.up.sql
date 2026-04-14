-- Создать первого администратора
-- Пароль: admin123
-- Hash создан с помощью bcrypt (cost 10)
INSERT INTO admins (username, password_hash) 
VALUES ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy');
