CREATE TABLE IF NOT EXISTS users (
  "ID" serial PRIMARY KEY NOT NULL,
  "Nickname" varchar NOT NULL,
  "Email" varchar NOT NULL,
  "CreatedAt" timestamp NOT NULL,
  "EncryptedPassword" varchar NOT NULL
);