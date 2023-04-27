CREATE TABLE IF NOT EXISTS "todo" (
  "ID" serial PRIMARY KEY,
  "User" varchar NOT NULL,
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

CREATE TABLE IF NOT EXISTS "users_configs" (
  "ID" varchar NOT NULL PRIMARY KEY,
  "Notification" bool NOT NULL DEFAULT false,
  "Period" INT,
  "CreatedAt" timestamp NOT NULL,
  "UpdatedAt" timestamp
);

ALTER TABLE "todo" ADD FOREIGN KEY ("User") REFERENCES "users_configs" ("ID");
