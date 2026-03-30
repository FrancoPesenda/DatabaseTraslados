-- ============================================================
--  Base de datos: Festival / Eventos
--  Motor: MySQL 8+
--  Generado a partir del diagrama de clases UML
-- ============================================================

CREATE DATABASE IF NOT EXISTS eventra
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE eventra;

-- ------------------------------------------------------------
-- LOCATION
-- ------------------------------------------------------------
CREATE TABLE location (
    id            INT             NOT NULL AUTO_INCREMENT,
    city          VARCHAR(100)    NOT NULL,
    province      VARCHAR(100)    NOT NULL,
    gps_coords    VARCHAR(100)    NOT NULL COMMENT 'Format: lat,lng  eg: -31.4135,-64.1811',
    PRIMARY KEY (id)
);

-- ------------------------------------------------------------
-- USER (base table — single-table inheritance with type)
-- ------------------------------------------------------------
CREATE TABLE user (
    id            INT             NOT NULL AUTO_INCREMENT,
    name          VARCHAR(100)    NOT NULL,
    last_name     VARCHAR(100)    NOT NULL,
    email         VARCHAR(150)    NOT NULL UNIQUE,
    password      VARCHAR(255)    NOT NULL,
    type          ENUM('admin','company','client') NOT NULL,
    -- Company-specific field (NULL for admin and client)
    company_name  VARCHAR(200)    NULL,
    PRIMARY KEY (id)
);

-- ------------------------------------------------------------
-- EVENT  (published by an Admin)
-- ------------------------------------------------------------
CREATE TABLE event (
    id            INT             NOT NULL AUTO_INCREMENT,
    name          VARCHAR(200)    NOT NULL,
    date          DATE            NOT NULL,
    location_id   INT             NOT NULL,
    admin_id      INT             NOT NULL COMMENT 'Admin user who published the event',
    PRIMARY KEY (id),
    CONSTRAINT fk_event_location FOREIGN KEY (location_id) REFERENCES location (id),
    CONSTRAINT fk_event_admin    FOREIGN KEY (admin_id)    REFERENCES user (id)
);

-- ------------------------------------------------------------
-- SERVICE  (published by a Company, appears in an Event)
-- ------------------------------------------------------------
CREATE TABLE service (
    id            INT             NOT NULL AUTO_INCREMENT,
    name          VARCHAR(200)    NOT NULL,
    description   TEXT            NULL,
    price         DECIMAL(10,2)   NOT NULL,
    company_id    INT             NOT NULL COMMENT 'Company user who published the service',
    PRIMARY KEY (id),
    CONSTRAINT fk_service_company FOREIGN KEY (company_id) REFERENCES user (id)
);

-- ------------------------------------------------------------
-- EVENT_SERVICE  (an event can have multiple services)
-- ------------------------------------------------------------
CREATE TABLE event_service (
    event_id      INT             NOT NULL,
    service_id    INT             NOT NULL,
    PRIMARY KEY (event_id, service_id),
    CONSTRAINT fk_es_event    FOREIGN KEY (event_id)   REFERENCES event    (id),
    CONSTRAINT fk_es_service  FOREIGN KEY (service_id) REFERENCES service  (id)
);

-- ------------------------------------------------------------
-- STAGE  (belongs to an Event)
-- ------------------------------------------------------------
CREATE TABLE stage (
    id            INT             NOT NULL AUTO_INCREMENT,
    name          VARCHAR(200)    NOT NULL,
    description   TEXT            NULL,
    schedule      VARCHAR(100)    NULL COMMENT 'eg: 18:00 - 23:00',
    event_id      INT             NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_stage_event FOREIGN KEY (event_id) REFERENCES event (id)
);

-- ------------------------------------------------------------
-- BAND
-- ------------------------------------------------------------
CREATE TABLE band (
    id            INT             NOT NULL AUTO_INCREMENT,
    name          VARCHAR(200)    NOT NULL,
    PRIMARY KEY (id)
);

-- ------------------------------------------------------------
-- STAGE_BAND  (a stage has multiple bands)
-- ------------------------------------------------------------
CREATE TABLE stage_band (
    stage_id      INT             NOT NULL,
    band_id       INT             NOT NULL,
    PRIMARY KEY (stage_id, band_id),
    CONSTRAINT fk_sb_stage FOREIGN KEY (stage_id) REFERENCES stage (id),
    CONSTRAINT fk_sb_band  FOREIGN KEY (band_id)  REFERENCES band  (id)
);

-- ------------------------------------------------------------
-- ARTIST
-- ------------------------------------------------------------
CREATE TABLE artist (
    id            INT             NOT NULL AUTO_INCREMENT,
    name          VARCHAR(100)    NOT NULL,
    last_name     VARCHAR(100)    NOT NULL,
    genre         VARCHAR(100)    NULL,
    PRIMARY KEY (id)
);

-- ------------------------------------------------------------
-- BAND_ARTIST  (a band has multiple artists, M:N)
-- ------------------------------------------------------------
CREATE TABLE band_artist (
    band_id       INT             NOT NULL,
    artist_id     INT             NOT NULL,
    PRIMARY KEY (band_id, artist_id),
    CONSTRAINT fk_ba_band    FOREIGN KEY (band_id)   REFERENCES band   (id),
    CONSTRAINT fk_ba_artist  FOREIGN KEY (artist_id) REFERENCES artist (id)
);