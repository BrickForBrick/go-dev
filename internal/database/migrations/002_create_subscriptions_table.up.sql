-- Создание таблицы подписок
CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    user_id UUID NOT NULL,
    start_date VARCHAR(7) NOT NULL CHECK (start_date ~ '^\d{2}-\d{4}$'),
    end_date VARCHAR(7) CHECK (end_date IS NULL OR end_date ~ '^\d{2}-\d{4}$'),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Внешний ключ на таблицу пользователей
    CONSTRAINT fk_subscriptions_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Создание индексов для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_service_name ON subscriptions(service_name);
CREATE INDEX IF NOT EXISTS idx_subscriptions_start_date ON subscriptions(start_date);
CREATE INDEX IF NOT EXISTS idx_subscriptions_end_date ON subscriptions(end_date) WHERE end_date IS NOT NULL;

-- Комментарии для документации
COMMENT ON TABLE subscriptions IS 'Таблица подписок пользователей на сервисы';
COMMENT ON COLUMN subscriptions.price IS 'Цена подписки в копейках/центах';
COMMENT ON COLUMN subscriptions.start_date IS 'Дата начала подписки в формате MM-YYYY';
COMMENT ON COLUMN subscriptions.end_date IS 'Дата окончания подписки в формате MM-YYYY (NULL = бессрочная)';