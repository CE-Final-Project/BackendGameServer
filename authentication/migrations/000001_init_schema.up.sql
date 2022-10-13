CREATE TABLE "accounts" (
    "id" SERIAL PRIMARY KEY,
    "player_id" varchar(15) UNIQUE,
    "username" varchar(50) UNIQUE,
    "email" varchar(320) UNIQUE,
    "password_hashed" varchar(255),
    "created_at" timestamptz,
    "updated_at" timestamptz
);

CREATE TABLE "roles" (
    "id" SERIAL PRIMARY KEY,
    "name" varchar,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "created_by" int,
    "updated_by" int
);

CREATE TABLE "account_roles" (
     "account_id" int,
     "role_id" int,
     "created_at" timestamptz
);

CREATE TABLE "friendship" (
    "account_id" int,
    "friend_id" int,
    "created_at" timestamptz
);

ALTER TABLE "roles" ADD FOREIGN KEY ("created_by") REFERENCES "accounts" ("id");

ALTER TABLE "roles" ADD FOREIGN KEY ("updated_by") REFERENCES "accounts" ("id");

ALTER TABLE "account_roles" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "account_roles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "friendship" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "friendship" ADD FOREIGN KEY ("friend_id") REFERENCES "accounts" ("id");
