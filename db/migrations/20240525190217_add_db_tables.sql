-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "orders" (
    "order_uid" VARCHAR(255) PRIMARY KEY,
    "track_number" VARCHAR(255),
    "entry" VARCHAR(255),
    "locale" VARCHAR(255),
    "internal_signature" VARCHAR(255),
    "customer_id" VARCHAR(255),
    "delivery_service" VARCHAR(255),
    "shardkey" VARCHAR(255),
    "sm_id" INTEGER,
    "date_created" TIMESTAMP,
    "oof_shard" VARCHAR(255)
);

-- Миграция для создания таблицы "Delivery"
CREATE TABLE IF NOT EXISTS "deliveries" (
    "id" SERIAL PRIMARY KEY,
    "order_uid" VARCHAR(255),
    "name" VARCHAR(255),
    "phone" VARCHAR(255),
    "zip" VARCHAR(255),
    "city" VARCHAR(255),
    "address" VARCHAR(255),
    "region" VARCHAR(255),
    "email" VARCHAR(255),
    FOREIGN KEY ("order_uid") REFERENCES "orders" ("order_uid") ON DELETE CASCADE
);

-- Миграция для создания таблицы "payments"
CREATE TABLE IF NOT EXISTS "payments" (
    "id" SERIAL PRIMARY KEY,
    "order_uid" VARCHAR(255),
    "transaction" VARCHAR(255),
    "request_id" VARCHAR(255),
    "currency" VARCHAR(255),
    "provider" VARCHAR(255),
    "amount" INTEGER,
    "payment_dt" INTEGER,
    "bank" VARCHAR(255),
    "delivery_cost" INTEGER,
    "goods_total" INTEGER,
    "custom_fee" INTEGER,
    FOREIGN KEY ("order_uid") REFERENCES "orders" ("order_uid") ON DELETE CASCADE
);

-- Миграция для создания таблицы "items"
CREATE TABLE IF NOT EXISTS "items" (
    "id" SERIAL PRIMARY KEY,
    "chrt_id" INTEGER,
    "order_uid" VARCHAR(255),
    "track_number" VARCHAR(255),
    "price" INTEGER,
    "rid" VARCHAR(255),
    "name" VARCHAR(255),
    "sale" INTEGER,
    "size" VARCHAR(255),
    "total_price" INTEGER,
    "nm_id" INTEGER,
    "brand" VARCHAR(255),
    "status" INTEGER,
    FOREIGN KEY ("order_uid") REFERENCES "orders" ("order_uid") ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Удаление внешних ключей
ALTER TABLE "items" DROP CONSTRAINT IF EXISTS "items_order_uid_fkey";
ALTER TABLE "payments" DROP CONSTRAINT IF EXISTS "payments_order_uid_fkey";
ALTER TABLE "deliveries" DROP CONSTRAINT IF EXISTS "deliveries_order_uid_fkey";

-- Удаление таблиц
DROP TABLE IF EXISTS "items";
DROP TABLE IF EXISTS "payments";
DROP TABLE IF EXISTS "deliveries";
DROP TABLE IF EXISTS "orders";
;
-- +goose StatementEnd
