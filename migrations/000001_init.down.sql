DROP INDEX IF EXISTS bot.idx_users_telegram_id;
DROP INDEX IF EXISTS bot.idx_cards_status;
DROP INDEX IF EXISTS bot.idx_transactions_created_at;
DROP INDEX IF EXISTS bot.idx_transactions_card_id;
DROP INDEX IF EXISTS bot.idx_transactions_user_id;

DROP TABLE IF EXISTS bot.transactions CASCADE;
DROP TABLE IF EXISTS bot.cards CASCADE;
DROP TABLE IF EXISTS bot.users CASCADE;

DROP TYPE IF EXISTS bot.type_enum CASCADE;
DROP TYPE IF EXISTS bot.card_status CASCADE;

DROP SCHEMA IF EXISTS bot CASCADE;