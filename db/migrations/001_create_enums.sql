-- migrate:up
CREATE TYPE member_role AS ENUM ('owner', 'admin', 'member');
CREATE TYPE account_type AS ENUM ('cash', 'bank', 'credit_card', 'e_money', 'investment');
CREATE TYPE transaction_type AS ENUM ('income', 'expense', 'transfer');
CREATE TYPE category_type AS ENUM ('income', 'expense');

-- migrate:down
DROP TYPE IF EXISTS category_type;
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS account_type;
DROP TYPE IF EXISTS member_role;
