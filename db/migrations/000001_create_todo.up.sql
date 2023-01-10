CREATE TABLE IF NOT EXISTS todo (
    "Id" serial not null unique ,
    "Description" varchar NOT NULL,
    "CreatedAt" timestamp NOT NULL,
    "DeadlineAt" timestamp NOT NULL,
    "FinishedAt" timestamp,
    "Completed" boolean default false
);