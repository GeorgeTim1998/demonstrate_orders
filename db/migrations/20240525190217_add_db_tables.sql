-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "order" (
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
CREATE TABLE IF NOT EXISTS "delivery" (
    "order_uid" VARCHAR(255) PRIMARY KEY,
    "name" VARCHAR(255),
    "phone" VARCHAR(255),
    "zip" VARCHAR(255),
    "city" VARCHAR(255),
    "address" VARCHAR(255),
    "region" VARCHAR(255),
    "email" VARCHAR(255)
);

-- Миграция для создания таблицы "payment"
CREATE TABLE IF NOT EXISTS "payment" (
    "order_uid" VARCHAR(255) PRIMARY KEY,
    "transaction" VARCHAR(255),
    "request_id" VARCHAR(255),
    "currency" VARCHAR(255),
    "provider" VARCHAR(255),
    "amount" INTEGER,
    "payment_dt" INTEGER,
    "bank" VARCHAR(255),
    "delivery_cost" INTEGER,
    "goods_total" INTEGER,
    "custom_fee" INTEGER
);

-- Миграция для создания таблицы "item"
CREATE TABLE IF NOT EXISTS "item" (
    "chrt_id" INTEGER PRIMARY KEY,
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
    "status" INTEGER
);

-- Создание внешних ключей
ALTER TABLE "delivery" ADD FOREIGN KEY ("order_uid") REFERENCES "order" ("order_uid");
ALTER TABLE "payment" ADD FOREIGN KEY ("order_uid") REFERENCES "order" ("order_uid");
ALTER TABLE "item" ADD FOREIGN KEY ("order_uid") REFERENCES "order" ("order_uid");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Удаление внешних ключей
ALTER TABLE "item" DROP CONSTRAINT "item_order_uid_fkey";
ALTER TABLE "payment" DROP CONSTRAINT "payment_order_uid_fkey";
ALTER TABLE "delivery" DROP CONSTRAINT "delivery_order_uid_fkey";

-- Удаление таблиц
DROP TABLE IF EXISTS "item";
DROP TABLE IF EXISTS "payment";
DROP TABLE IF EXISTS "delivery";
DROP TABLE IF EXISTS "order";
;
-- +goose StatementEnd
