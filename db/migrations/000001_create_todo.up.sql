CREATE TABLE IF NOT EXISTS "todo" (
  "ID" serial PRIMARY KEY,
  "User" serial NOT NULL,
  "Title" varchar NOT NULL,
  "Description" varchar NOT NULL,
  "CreatedAt" timestamp NOT NULL,
  "DeadlineAt" timestamp NOT NULL,
  "FinishedAt" timestamp,
  "Completed" boolean NOT NULL DEFAULT false,
  "Notificated" boolean NOT NULL DEFAULT false,
  "DeadlineNotify" bool,
  "NotifyPeriod" timestamp []
);

CREATE INDEX ON "todo" ("User");

CREATE TABLE IF NOT EXISTS "users" (
  "ID" serial PRIMARY KEY,
  "Name" varchar NOT NULL,
  "Email" varchar NOT NULL UNIQUE ,
  "Password" varchar NOT NULL,
  "Provider" varchar NOT NULL,
  "Verified" bool DEFAULT false,
  "CreatedAt" timestamp NOT NULL,
  "UpdatedAt" timestamp,
  "Notification" bool NOT NULL DEFAULT false,
  "Period" INT
);

ALTER TABLE "todo" ADD FOREIGN KEY ("User") REFERENCES "users" ("ID") ON DELETE CASCADE;
