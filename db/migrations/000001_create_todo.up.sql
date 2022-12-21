CREATE TABLE IF NOT EXISTS todo (
    "Id" serial primary key ,
    "Description" varchar NOT NULL,
    "CreatedAt" timestamptz NOT NULL,
    "DeadlineAt" timestamptz NOT NULL,
    "FinishedAt" timestamptz,
    "Completed" boolean default false
);