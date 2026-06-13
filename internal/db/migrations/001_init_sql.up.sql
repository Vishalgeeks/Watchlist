CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS stocks (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    symbol        VARCHAR(20) UNIQUE NOT NULL,
    company_name  VARCHAR(255) NOT NULL,
    exchange      VARCHAR(50),
    current_price DECIMAL(12,2) DEFAULT 0.00,
    last_updated  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS watchlists (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_watchlists_user ON watchlists(user_id);

CREATE TABLE IF NOT EXISTS watchlist_items (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    watchlist_id UUID NOT NULL REFERENCES watchlists(id) ON DELETE CASCADE,
    stock_id     UUID NOT NULL REFERENCES stocks(id)     ON DELETE CASCADE,
    added_at     TIMESTAMP DEFAULT NOW(),
    UNIQUE(watchlist_id, stock_id)
);

CREATE INDEX IF NOT EXISTS idx_items_watchlist ON watchlist_items(watchlist_id);

INSERT INTO users (id, name) VALUES
    ('00000000-0000-0000-0000-000000000001', 'Vish'),
    ('00000000-0000-0000-0000-000000000002', 'Rahul')
ON CONFLICT DO NOTHING;

INSERT INTO stocks (symbol, company_name, exchange, current_price) VALUES
    ('RELIANCE',   'Reliance Industries Ltd',   'NSE', 2450.75),
    ('TCS',        'Tata Consultancy Services', 'NSE', 3820.50),
    ('INFY',       'Infosys Ltd',               'NSE', 1456.30),
    ('HDFCBANK',   'HDFC Bank Ltd',             'NSE', 1678.90),
    ('TATAMOTORS', 'Tata Motors Ltd',           'NSE',  875.40),
    ('WIPRO',      'Wipro Ltd',                 'NSE',  456.20),
    ('ICICIBANK',  'ICICI Bank Ltd',            'NSE', 1123.60),
    ('SBIN',       'State Bank of India',       'NSE',  734.85)
ON CONFLICT (symbol) DO NOTHING;