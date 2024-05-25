-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "Order" (
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
CREATE TABLE IF NOT EXISTS "Delivery" (
    "order_uid" VARCHAR(255) PRIMARY KEY,
    "name" VARCHAR(255),
    "phone" VARCHAR(255),
    "zip" VARCHAR(255),
    "city" VARCHAR(255),
    "address" VARCHAR(255),
    "region" VARCHAR(255),
    "email" VARCHAR(255)
);

-- Миграция для создания таблицы "Payment"
CREATE TABLE IF NOT EXISTS "Payment" (
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

-- Миграция для создания таблицы "Item"
CREATE TABLE IF NOT EXISTS "Item" (
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
ALTER TABLE "Delivery" ADD FOREIGN KEY ("order_uid") REFERENCES "Order" ("order_uid");
ALTER TABLE "Payment" ADD FOREIGN KEY ("order_uid") REFERENCES "Order" ("order_uid");
ALTER TABLE "Item" ADD FOREIGN KEY ("order_uid") REFERENCES "Order" ("order_uid");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Удаление внешних ключей
ALTER TABLE "Item" DROP CONSTRAINT "Item_order_uid_fkey";
ALTER TABLE "Payment" DROP CONSTRAINT "Payment_order_uid_fkey";
ALTER TABLE "Delivery" DROP CONSTRAINT "Delivery_order_uid_fkey";

-- Удаление таблиц
DROP TABLE IF EXISTS "Item";
DROP TABLE IF EXISTS "Payment";
DROP TABLE IF EXISTS "Delivery";
DROP TABLE IF EXISTS "Order";
;
-- +goose StatementEnd
