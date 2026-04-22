CREATE TABLE IF NOT EXISTS "user" (
  "id" SERIAL PRIMARY KEY,
  "public_id" UUID UNIQUE,
  "email" VARCHAR,
  "passwordHash" VARCHAR,
  "createAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "changeAt" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "contests" (
  "id" SERIAL PRIMARY KEY,
  "public_id" UUID UNIQUE,
  "user_id" INTEGER,
  "name" varchar NOT NULL,
  "description" text,
  "is_up" boolean DEFAULT false,
  "max_votes_user" INTEGER DEFAULT 0,
  "create_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "change_at" TIMESTAMP,
  FOREIGN KEY ("user_id") REFERENCES "user" ("id") DEFERRABLE INITIALLY IMMEDIATE
);

CREATE TABLE IF NOT EXISTS "participants" (
  "id" SERIAL PRIMARY KEY,
  "public_id" UUID UNIQUE,
  "contest_id" INTEGER NOT NULL,
  "email" varchar,
  "nome" varchar NOT NULL,
  "telefone" varchar,
  "title" varchar NOT NULL,
  "description" text,
  "image_path" varchar,
  FOREIGN KEY ("contest_id") REFERENCES "contests" ("id") DEFERRABLE INITIALLY IMMEDIATE
);

CREATE TABLE IF NOT EXISTS "votes" (
  "id" SERIAL PRIMARY KEY,
  "contest_id" INTEGER NOT NULL,
  "participant_id" INTEGER NOT NULL,
  "voter_hash" varchar NULL,
  "voter_ip" varchar NULL,
  "voter_random_cookie" varchar NULL,
  "voted_at" timestamp,
  FOREIGN KEY ("contest_id") REFERENCES "contests" ("id") DEFERRABLE INITIALLY IMMEDIATE,
  FOREIGN KEY ("participant_id") REFERENCES "participants" ("id") DEFERRABLE INITIALLY IMMEDIATE
);

CREATE TABLE IF NOT EXISTS "vote_options" (
  "id" SERIAL PRIMARY KEY,
  "contest_id" INTEGER NOT NULL,
  "vote_id" INTEGER NOT NULL,
  "counter_option_id" UUID NOT NULL,
  FOREIGN KEY ("contest_id") REFERENCES "contests" ("id") DEFERRABLE INITIALLY IMMEDIATE,
  FOREIGN KEY ("vote_id") REFERENCES "votes" ("id") DEFERRABLE INITIALLY IMMEDIATE,
  FOREIGN KEY ("counter_option_id") REFERENCES "participants" ("public_id") DEFERRABLE INITIALLY IMMEDIATE
);

