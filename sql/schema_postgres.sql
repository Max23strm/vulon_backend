-- =====================================================================
-- pitz_backend - PostgreSQL schema
-- =====================================================================
-- Equivalente al esquema MySQL original.
-- Tipos de columna basados en los modelos Go (models/*.go) y las queries
-- SQL usadas en routes/*.go.
-- =====================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================================================================
-- USERS
-- =====================================================================
CREATE TABLE IF NOT EXISTS "users" (
    "user_uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "username" VARCHAR(100) NOT NULL UNIQUE,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "hashed_password" VARCHAR(255) NOT NULL,
    "first_name" VARCHAR(100) NOT NULL,
    "last_name" VARCHAR(100) NOT NULL,
    "created_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag" SMALLINT NOT NULL DEFAULT 0
);

-- =====================================================================
-- PLAYERS
-- =====================================================================
CREATE TABLE IF NOT EXISTS "players" (
    "player_uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "first_name" VARCHAR(100) NOT NULL,
    "last_name" VARCHAR(100) NOT NULL,
    "email" VARCHAR(255),
    "entity_uid" UUID NOT NULL,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "positions" JSONB,
    "phone_number" VARCHAR(30),
    "emergency_phone" VARCHAR(30),
    "birth_dt" TIMESTAMP,
    "blood_type" VARCHAR(10),
    "comments" TEXT,
    "credential" VARCHAR(100),
    "address" TEXT,
    "afiliation" VARCHAR(100),
    "sex" VARCHAR(20),
    "curp" VARCHAR(50),
    "enfermedad" TEXT,
    "insurance" BOOLEAN NOT NULL DEFAULT FALSE,
    "insurance_name" VARCHAR(100),
    "created_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag" SMALLINT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS "idx_players_first_name" ON "players" ("first_name");
CREATE INDEX IF NOT EXISTS "idx_players_entity" ON "players" ("entity_uid");

-- =====================================================================
-- EVENT TYPES
-- =====================================================================
CREATE TABLE IF NOT EXISTS "event_types" (
    "event_type_uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "type_name" VARCHAR(100) NOT NULL
);

-- =====================================================================
-- EVENT STATES
-- =====================================================================
CREATE TABLE IF NOT EXISTS "events_state" (
    "event_state_uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "event_state" VARCHAR(100) NOT NULL,
    "event_state_type" INTEGER NOT NULL,
);

-- =====================================================================
-- EVENTS
-- =====================================================================
CREATE TABLE IF NOT EXISTS "events" (
    "event_uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "event_type" UUID NOT NULL,
    "event_state_uid" UUID NOT NULL,
    "event_name" VARCHAR(255) NOT NULL,
    "date" TIMESTAMP NOT NULL,
    "address" TEXT,
    "coordinates" VARCHAR(255),
    "created_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag" SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT "fk_events_event_type"
        FOREIGN KEY ("event_type") REFERENCES "event_types" ("event_type_uid"),
    CONSTRAINT "fk_events_event_state"
        FOREIGN KEY ("event_state_uid") REFERENCES "events_state" ("event_state_uid")
    
);

-- =====================================================================
-- ASISTANCE TYPES
-- =====================================================================
CREATE TABLE IF NOT EXISTS "asistance_types" (
    "asistance_type_uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR(100) NOT NULL
);

-- =====================================================================
-- ASISTANCE
-- =====================================================================
CREATE TABLE IF NOT EXISTS "asistance" (
    "player_uid" UUID NOT NULL,
    "event_uid" UUID NOT NULL,
    "asistance_type_uid" UUID NOT NULL,
    "created_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("player_uid", "event_uid"),
    CONSTRAINT "fk_asistance_player"
        FOREIGN KEY ("player_uid") REFERENCES "players" ("player_uid"),
    CONSTRAINT "fk_asistance_event"
        FOREIGN KEY ("event_uid") REFERENCES "events" ("event_uid"),
    CONSTRAINT "fk_asistance_type"
        FOREIGN KEY ("asistance_type_uid") REFERENCES "asistance_types" ("asistance_type_uid")
);

-- =====================================================================
-- PAYMENT TYPES
-- =====================================================================
CREATE TABLE IF NOT EXISTS "payment_type" (
    "payment_type_uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "payment_name" VARCHAR(100) NOT NULL,
    "system" SMALLINT NOT NULL DEFAULT 0,
    "entity_uid" UUID NOT NULL,
    CONSTRAINT "fk_payment_type_entity"
        FOREIGN KEY ("entity_uid") REFERENCES "entities" ("entity_uid"),
);

-- =====================================================================
-- PAYMENT CATEGORIES
-- =====================================================================
CREATE TABLE IF NOT EXISTS "payment_category" (
    "payment_cat_uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "payment_type_uid" UUID NOT NULL,
    "amount" NUMERIC(12, 2) NOT NULL,
    "name" VARCHAR(100) NOT NULL,
    "delete_flag" SMALLINT NOT NULL DEFAULT 0,
    "entity_uid" UUID NOT NULL,
    CONSTRAINT "fk_payment_cat_entity"
        FOREIGN KEY ("entity_uid") REFERENCES "entities" ("entity_uid"),
    CONSTRAINT "fk_payment_cat_type"
        FOREIGN KEY ("payment_type_uid") REFERENCES "payment_type" ("payment_type_uid"),
);

-- =====================================================================
-- PAYMENTS
-- =====================================================================
CREATE TABLE IF NOT EXISTS "payments" (
    "payment_uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "player_uid" UUID NOT NULL,
    "payment_reference" VARCHAR(255),
    "amount" NUMERIC(12, 2) NOT NULL,
    "comment" TEXT,
    "date" TIMESTAMP NOT NULL,
    "payment_type_uid" UUID NOT NULL,
    "registered_by_uid" UUID NOT NULL,
    "created_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag" SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT "fk_payments_player"
        FOREIGN KEY ("player_uid") REFERENCES "players" ("player_uid"),
    CONSTRAINT "fk_payments_payment_type"
        FOREIGN KEY ("payment_type_uid") REFERENCES "payment_type" ("payment_type_uid"),
    CONSTRAINT "fk_payments_registered_by"
        FOREIGN KEY ("registered_by_uid") REFERENCES "users" ("user_uid")
);

CREATE INDEX IF NOT EXISTS "idx_payments_date" ON "payments" ("date");
CREATE INDEX IF NOT EXISTS "idx_payments_delete_flag" ON "payments" ("delete_flag");

-- =====================================================================
-- EXPENSES
-- =====================================================================
CREATE TABLE IF NOT EXISTS "expenses" (
    "expense_uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "assigned_uid" UUID NOT NULL,
    "registered_by_uid" UUID NOT NULL,
    "reason" TEXT NOT NULL,
    "amount" NUMERIC(12, 2) NOT NULL,
    "date" TIMESTAMP NOT NULL,
    "created_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag" SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT "fk_expenses_assigned"
        FOREIGN KEY ("assigned_uid") REFERENCES "users" ("user_uid"),
    CONSTRAINT "fk_expenses_registered_by"
        FOREIGN KEY ("registered_by_uid") REFERENCES "users" ("user_uid")
);

CREATE INDEX IF NOT EXISTS "idx_expenses_date" ON "expenses" ("date");
CREATE INDEX IF NOT EXISTS "idx_expenses_delete_flag" ON "expenses" ("delete_flag");

-- =====================================================================
-- ENTITIES
-- =====================================================================
CREATE TABLE IF NOT EXISTS "entities" (
    "entity_uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "short_name" VARCHAR(36) NOT NULL,
    "name" VARCHAR(52) NOT NULL,
    "currency_code" CHAR(3) NOT NULL DEFAULT 'USD',
    "country_code" CHAR(4) NOT NULL DEFAULT 'USA',
    "colors" VARCHAR(7)[] NOT NULL DEFAULT '{}',
    "logo" TEXT,
    "created_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag" SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT "chk_colors_max_three"
        CHECK (array_length("colors", 1) IS NULL OR array_length("colors", 1) <= 3),
    CONSTRAINT "fk_events_entity"
        FOREIGN KEY ("entity_uid") REFERENCES "entities" ("entity_uid"),
    CONSTRAINT "fk_players_entity"
        FOREIGN KEY ("entity_uid") REFERENCES "entities" ("entity_uid")
);

-- =====================================================================
-- USER ENTITIES ASSIGNATION
-- =====================================================================
CREATE TABLE IF NOT EXISTS "user_entities" (
    "user_entity_uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "user_uid" UUID NOT NULL,
    "entity_uid" UUID NOT NULL,
    "assigned_by_uid" UUID NOT NULL,
    "created_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag" SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT "fk_user_entities_user"
        FOREIGN KEY ("user_uid") REFERENCES "users" ("user_uid"),
    CONSTRAINT "fk_user_entities_entity"
        FOREIGN KEY ("entity_uid") REFERENCES "entities" ("entity_uid"),
    CONSTRAINT "fk_user_entities_assigned_by"
        FOREIGN KEY ("assigned_by_uid") REFERENCES "users" ("user_uid")
);

CREATE UNIQUE INDEX "uq_user_entity_active"
    ON "user_entities" ("user_uid", "entity_uid")
    WHERE "delete_flag" = 0;

-- =====================================================================
-- TEAMS
-- =====================================================================
CREATE TABLE IF NOT EXISTS "entity_teams" (
    "team_uid"         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "description"      VARCHAR(50) NOT NULL,
    "entity_uid"       UUID NOT NULL,
    "created_at_dttm"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag"      SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT "fk_entity_teams_entity"
        FOREIGN KEY ("entity_uid") REFERENCES "entities" ("entity_uid")
);

-- =====================================================================
-- TEAMS CATEGORY
-- =====================================================================
CREATE TABLE IF NOT EXISTS "team_categories" (
    "category_uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "description" VARCHAR(50) NOT NULL,
    "entity_uid" UUID NOT NULL,
    "created_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag" SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT "fk_team_categories_entity" FOREIGN KEY ("entity_uid") REFERENCES "entities" ("entity_uid")
);

-- =====================================================================
-- TEAMS CATEGORY ASSIGNATION
-- =====================================================================
CREATE TABLE IF NOT EXISTS "entity_team_categories" (
    "team_uid" UUID NOT NULL,
    "category_uid" UUID NOT NULL,

    PRIMARY KEY ("team_uid", "category_uid"),

    CONSTRAINT "fk_team_categories_team"
        FOREIGN KEY ("team_uid")
        REFERENCES "entity_teams" ("team_uid")
        ON DELETE CASCADE,

    CONSTRAINT "fk_team_categories_category"
        FOREIGN KEY ("category_uid")
        REFERENCES "team_categories" ("category_uid")
        ON DELETE CASCADE
);

-- =====================================================================
-- PLAYERS ASSIGNATION
-- =====================================================================
CREATE TABLE IF NOT EXISTS "players_assignation" (
    "player_assignation_uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "player_uid" UUID NOT NULL,
    "team_uid" UUID NOT NULL,
    "assigned_by_uid" UUID NOT NULL,
    "created_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag" SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT "fk_players_assignation_player"
        FOREIGN KEY ("player_uid") REFERENCES "players" ("player_uid"),
    CONSTRAINT "fk_players_assignation_team"
        FOREIGN KEY ("team_uid") REFERENCES "entity_teams" ("team_uid"),
    CONSTRAINT "fk_players_assignation_assigned_by"
        FOREIGN KEY ("assigned_by_uid") REFERENCES "users" ("user_uid")
);

-- =====================================================================
-- Catálogos iniciales
-- =====================================================================
INSERT INTO "event_types" ("event_type_uid", "type_name") VALUES
    (uuid_generate_v4(), 'Entrenamiento'),
    (uuid_generate_v4(), 'Partido'),
    (uuid_generate_v4(), 'Reunión');

INSERT INTO "events_state" ("event_state_uid", "event_state") VALUES
    (uuid_generate_v4(), 'Pendiente'),
    (uuid_generate_v4(), 'Confirmado'),
    (uuid_generate_v4(), 'Cancelado'),
    (uuid_generate_v4(), 'Finalizado');

INSERT INTO "asistance_types" ("asistance_type_uid", "name") VALUES
    (uuid_generate_v4(), 'Asistió'),
    (uuid_generate_v4(), 'Faltó'),
    (uuid_generate_v4(), 'Justificado');

INSERT INTO "payment_type" ("payment_type_uid", "payment_name") VALUES
    (uuid_generate_v4(), 'Mensualidad'),
    (uuid_generate_v4(), 'Inscripción'),
    (uuid_generate_v4(), 'Multa');
