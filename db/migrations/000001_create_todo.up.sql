CREATE TABLE IF NOT EXISTS todo (
  "ID" serial PRIMARY KEY NOT NULL,
  "User" bigserial NOT NULL,
  "Description" varchar NOT NULL,
  "CreatedAt" timestamp NOT NULL,
  "DeadlineAt" timestamp NOT NULL,
  "FinishedAt" timestamp,
  "Completed" boolean DEFAULT false
);

CREATE INDEX ON "todo" ("User");

ALTER TABLE "todo" ADD FOREIGN KEY ("User") REFERENCES "users" ("ID");