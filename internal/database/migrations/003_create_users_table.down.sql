-- Удаление внешнего ключа
ALTER TABLE subscriptions DROP CONSTRAINT IF EXISTS fk_subscriptions_user_id;

-- Возвращение типа user_id обратно на VARCHAR
ALTER TABLE subscriptions ALTER COLUMN user_id TYPE VARCHAR(255);

-- Удаление таблицы пользователей
DROP TABLE IF EXISTS users;