CREATE TABLE "account" (
                           "id" BigSerial PRIMARY KEY,
                           "owner" varchar NOT NULL,
                           "balance" bigint NOT NULL,
                           "currency" varchar NOT NULL,
                           "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
                           "id" BigSerial PRIMARY KEY,
                           "amount" bigint NOT NULL,
                           "account_id" bigint,
                           "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
                             "id" BigSerial PRIMARY KEY,
                             "amount" bigint NOT NULL,
                             "sender_id" bigint NOT NULL,
                             "receiver_id" bigint NOT NULL,
                             "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "account" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("sender_id");

CREATE INDEX ON "transfers" ("receiver_id");

CREATE INDEX ON "transfers" ("sender_id", "receiver_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be positive or negative';

COMMENT ON COLUMN "transfers"."amount" IS 'should be only positive';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "account" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("sender_id") REFERENCES "account" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("receiver_id") REFERENCES "account" ("id");
