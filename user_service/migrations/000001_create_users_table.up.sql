CREATE TABLE IF NOT EXISTS "users" (
    "id" UUID,
    "first_name" VARCHAR(30) NOT NULL,
    "last_name" VARCHAR(30) NOT NULL,
    "email" VARCHAR(30) NOT NULL UNIQUE,
    "user_type" VARCHAR(50) NOT NULL,
    "password" TEXT NOT NULL,
    "refresh_token" TEXT NOT NULL,
    "created_at" TIME DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIME DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIME
);