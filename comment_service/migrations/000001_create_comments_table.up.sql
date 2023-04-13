CREATE TABLE IF NOT EXISTS "comments" (
    "id" SERIAL PRIMARY KEY,
    "post_id" INTEGER,
    "user_id" UUID NOT NULL,
    "text" TEXT,
    "created_at" TIME DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIME DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIME
)