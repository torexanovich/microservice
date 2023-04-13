CREATE TABLE IF NOT EXISTS "posts" (
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR(60),
    "description" VARCHAR(60),
    "likes" INTEGER DEFAULT 0,
    "user_id" UUID NOT NULL,
    "created_at" TIME DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIME DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIME
)