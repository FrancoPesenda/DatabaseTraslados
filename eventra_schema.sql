-- =============================================================
--  DATABASE SCHEMA
--  Generated from ERD v5
--  Tables: 23
-- =============================================================

CREATE DATABASE IF NOT EXISTS eventra
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE eventra;

-- =============================================================
--  USERS & PERMISSIONS
-- =============================================================

CREATE TABLE role (
  id    INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  name  VARCHAR(50)   NOT NULL UNIQUE,
  PRIMARY KEY (id)
);

CREATE TABLE permission (
  id           INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  name         VARCHAR(100)  NOT NULL UNIQUE,
  description  VARCHAR(255)      NULL,
  PRIMARY KEY (id)
);

CREATE TABLE role_permission (
  role_id       INT UNSIGNED  NOT NULL,
  permission_id INT UNSIGNED  NOT NULL,
  PRIMARY KEY (role_id, permission_id),
  CONSTRAINT fk_rp_role
    FOREIGN KEY (role_id)       REFERENCES role(id)       ON DELETE CASCADE,
  CONSTRAINT fk_rp_permission
    FOREIGN KEY (permission_id) REFERENCES permission(id) ON DELETE CASCADE
);

CREATE TABLE user (
  id              INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  username        VARCHAR(100)      NULL UNIQUE,
  password        VARCHAR(255)  NOT NULL,
  email           VARCHAR(255)      NULL,
  name            VARCHAR(100)      NULL,
  last_name       VARCHAR(100)      NULL,
  role_id         INT UNSIGNED      NULL,
  account_expiry  DATE              NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_user_role
    FOREIGN KEY (role_id) REFERENCES role(id)
);

INSERT INTO role (name) VALUES ('admin'), ('company');

-- =============================================================
--  GEOGRAPHY
-- =============================================================

CREATE TABLE country (
  id    INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  name  VARCHAR(100)  NOT NULL UNIQUE,
  PRIMARY KEY (id)
);

CREATE TABLE province (
  id          INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  name        VARCHAR(100)  NOT NULL,
  country_id  INT UNSIGNED  NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_province_country
    FOREIGN KEY (country_id) REFERENCES country(id)
);

CREATE TABLE city (
  id           INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  name         VARCHAR(100)  NOT NULL,
  province_id  INT UNSIGNED  NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_city_province
    FOREIGN KEY (province_id) REFERENCES province(id)
);

CREATE TABLE location (
  id       INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  stadium  VARCHAR(150)  NOT NULL,
  city_id  INT UNSIGNED  NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_location_city
    FOREIGN KEY (city_id) REFERENCES city(id)
);

-- =============================================================
--  COMPANIES & SERVICES
-- =============================================================

CREATE TABLE company (
  id    INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  name  VARCHAR(150)  NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE service (
  id          INT UNSIGNED    NOT NULL AUTO_INCREMENT,
  date        DATE            NOT NULL,
  price       DECIMAL(10, 2)  NOT NULL,
  company_id  INT UNSIGNED    NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_service_company
    FOREIGN KEY (company_id) REFERENCES company(id)
);

-- =============================================================
--  REQUESTS
-- =============================================================

CREATE TABLE request (
  id            INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  type          VARCHAR(50)   NOT NULL,
  description   TEXT              NULL,
  requester_id  INT UNSIGNED  NOT NULL COMMENT 'User who creates the request',
  approver_id   INT UNSIGNED      NULL COMMENT 'Admin who approves the request',
  PRIMARY KEY (id),
  CONSTRAINT fk_request_requester
    FOREIGN KEY (requester_id) REFERENCES user(id),
  CONSTRAINT fk_request_approver
    FOREIGN KEY (approver_id)  REFERENCES user(id)
);

-- =============================================================
--  MUSIC
-- =============================================================

CREATE TABLE band (
  id    INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  name  VARCHAR(150)  NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE artist (
  id    INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  name  VARCHAR(150)  NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE band_artist (
  band_id    INT UNSIGNED  NOT NULL,
  artist_id  INT UNSIGNED  NOT NULL,
  PRIMARY KEY (band_id, artist_id),
  CONSTRAINT fk_ba_band
    FOREIGN KEY (band_id)   REFERENCES band(id)   ON DELETE CASCADE,
  CONSTRAINT fk_ba_artist
    FOREIGN KEY (artist_id) REFERENCES artist(id) ON DELETE CASCADE
);

-- =============================================================
--  EVENTS
-- =============================================================

CREATE TABLE event (
  id           INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  name         VARCHAR(150)  NOT NULL,
  location_id  INT UNSIGNED  NOT NULL,
  start_date   DATE          NOT NULL,
  end_date     DATE          NOT NULL,
  image        VARCHAR(500)      NULL,
  request_id   INT UNSIGNED      NULL,
  service_id   INT UNSIGNED      NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_event_location
    FOREIGN KEY (location_id) REFERENCES location(id),
  CONSTRAINT fk_event_request
    FOREIGN KEY (request_id)  REFERENCES request(id),
  CONSTRAINT fk_event_service
    FOREIGN KEY (service_id)  REFERENCES service(id)
);

CREATE TABLE event_day (
  id        INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  event_id  INT UNSIGNED  NOT NULL,
  date      DATE          NOT NULL,
  name      VARCHAR(100)      NULL COMMENT 'e.g. Day 1 - Opening Night',
  PRIMARY KEY (id),
  CONSTRAINT fk_event_day_event
    FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE
);

CREATE TABLE transport (
  id          INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  position    VARCHAR(150)  NOT NULL,
  event_id    INT UNSIGNED  NOT NULL,
  service_id  INT UNSIGNED  NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_transport_event
    FOREIGN KEY (event_id)   REFERENCES event(id),
  CONSTRAINT fk_transport_service
    FOREIGN KEY (service_id) REFERENCES service(id)
);

-- =============================================================
--  ROSTER
-- =============================================================

CREATE TABLE roster (
  id        INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  event_id  INT UNSIGNED  NOT NULL UNIQUE COMMENT '1:1 with event',
  PRIMARY KEY (id),
  CONSTRAINT fk_roster_event
    FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE
);

CREATE TABLE stage (
  id         INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  name       VARCHAR(100)  NOT NULL,
  capacity   INT UNSIGNED      NULL,
  roster_id  INT UNSIGNED  NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_stage_roster
    FOREIGN KEY (roster_id) REFERENCES roster(id) ON DELETE CASCADE
);

CREATE TABLE roster_item (
  id            INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  roster_id     INT UNSIGNED  NOT NULL,
  stage_id      INT UNSIGNED  NOT NULL,
  band_id       INT UNSIGNED  NOT NULL,
  event_day_id  INT UNSIGNED  NOT NULL,
  start_time    TIME          NOT NULL,
  end_time      TIME          NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_ri_roster
    FOREIGN KEY (roster_id)    REFERENCES roster(id),
  CONSTRAINT fk_ri_stage
    FOREIGN KEY (stage_id)     REFERENCES stage(id),
  CONSTRAINT fk_ri_band
    FOREIGN KEY (band_id)      REFERENCES band(id),
  CONSTRAINT fk_ri_event_day
    FOREIGN KEY (event_day_id) REFERENCES event_day(id)
);

-- =============================================================
--  PIVOT TABLES
-- =============================================================

CREATE TABLE home (
  id  INT UNSIGNED  NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (id)
);

CREATE TABLE home_event (
  home_id    INT UNSIGNED         NOT NULL,
  event_id   INT UNSIGNED         NOT NULL,
  order_num  TINYINT UNSIGNED     NOT NULL DEFAULT 0,
  PRIMARY KEY (home_id, event_id),
  CONSTRAINT fk_he_home
    FOREIGN KEY (home_id)  REFERENCES home(id)  ON DELETE CASCADE,
  CONSTRAINT fk_he_event
    FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE
);

CREATE TABLE event_company (
  event_id        INT UNSIGNED  NOT NULL,
  company_id      INT UNSIGNED  NOT NULL,
  coordinator_id  INT UNSIGNED  NOT NULL COMMENT 'User with coordinator role',
  PRIMARY KEY (event_id, company_id),
  CONSTRAINT fk_ec_event
    FOREIGN KEY (event_id)       REFERENCES event(id)   ON DELETE CASCADE,
  CONSTRAINT fk_ec_company
    FOREIGN KEY (company_id)     REFERENCES company(id),
  CONSTRAINT fk_ec_coordinator
    FOREIGN KEY (coordinator_id) REFERENCES user(id)
);