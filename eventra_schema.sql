-- eventra_schema.sql
-- Eventos Traslados — Esquema de base de datos relacional
-- Compatible con MySQL 8.0.13+
-- ============================================================

CREATE DATABASE IF NOT EXISTS eventra
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE eventra;

-- ============================================================
-- 1. admins
--    Usuarios administradores. Soporta tipo 'admin' (completo)
--    y 'event_creator' (solo creación de eventos).
-- ============================================================
CREATE TABLE admins (
    id                   VARCHAR(100)  NOT NULL,
    username             VARCHAR(160)  NOT NULL,
    email                VARCHAR(160)  NOT NULL DEFAULT '',
    password_hash        TEXT          NOT NULL,
    current_password     VARCHAR(255)  NOT NULL DEFAULT '',
    name                 VARCHAR(160)  NOT NULL DEFAULT '',
    last_name            VARCHAR(160)  NOT NULL DEFAULT '',
    account_type         VARCHAR(30)   NOT NULL DEFAULT 'admin',
    can_manage_admins    BOOLEAN       NOT NULL DEFAULT FALSE,
    recovery_email       VARCHAR(160)  NOT NULL DEFAULT '',
    responsable_admin_id VARCHAR(100)  NULL DEFAULT NULL,
    created_at           BIGINT        NOT NULL,
    updated_at           BIGINT        NULL DEFAULT NULL,
    PRIMARY KEY (id),
    CONSTRAINT chk_admins_account_type
        CHECK (account_type IN ('admin', 'event_creator')),
    CONSTRAINT fk_admins_responsable
        FOREIGN KEY (responsable_admin_id) REFERENCES admins(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 2. company_profiles
--    Perfil de empresa 1:1 con cada admin.
-- ============================================================
CREATE TABLE company_profiles (
    id             INT           NOT NULL AUTO_INCREMENT,
    admin_id       VARCHAR(100)  NOT NULL,
    company_name   VARCHAR(120)  NOT NULL DEFAULT '',
    instagram      VARCHAR(160)  NOT NULL DEFAULT '',
    transfer_alias VARCHAR(120)  NOT NULL DEFAULT '',
    cbu            VARCHAR(40)   NOT NULL DEFAULT '',
    owner_name     VARCHAR(120)  NOT NULL DEFAULT '',
    bank           VARCHAR(120)  NOT NULL DEFAULT '',
    transfer_phone VARCHAR(40)   NOT NULL DEFAULT '',
    updated_at     BIGINT        NULL DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_company_profiles_admin (admin_id),
    CONSTRAINT fk_company_profiles_admin
        FOREIGN KEY (admin_id) REFERENCES admins(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 3. company_payment_accounts
--    Cuentas de cobro asociadas al perfil de empresa de un admin.
-- ============================================================
CREATE TABLE company_payment_accounts (
    id             VARCHAR(100)  NOT NULL,
    admin_id       VARCHAR(100)  NOT NULL,
    alias_name     VARCHAR(80)   NOT NULL DEFAULT '',
    owner_name     VARCHAR(120)  NOT NULL DEFAULT '',
    bank           VARCHAR(120)  NOT NULL DEFAULT '',
    cbu            VARCHAR(40)   NOT NULL DEFAULT '',
    transfer_alias VARCHAR(120)  NOT NULL DEFAULT '',
    PRIMARY KEY (id),
    CONSTRAINT fk_payment_accounts_admin
        FOREIGN KEY (admin_id) REFERENCES admins(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 4. company_contact_numbers
--    Números de contacto de la empresa.
-- ============================================================
CREATE TABLE company_contact_numbers (
    id           VARCHAR(100)  NOT NULL,
    admin_id     VARCHAR(100)  NOT NULL,
    alias_name   VARCHAR(80)   NOT NULL DEFAULT '',
    phone_number VARCHAR(40)   NOT NULL DEFAULT '',
    PRIMARY KEY (id),
    CONSTRAINT fk_contact_numbers_admin
        FOREIGN KEY (admin_id) REFERENCES admins(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 5. company_auxiliary_links
--    Links auxiliares de la empresa (redes, páginas, etc.).
-- ============================================================
CREATE TABLE company_auxiliary_links (
    id         VARCHAR(100)  NOT NULL,
    admin_id   VARCHAR(100)  NOT NULL,
    alias_name VARCHAR(80)   NOT NULL DEFAULT '',
    url        VARCHAR(300)  NOT NULL DEFAULT '',
    PRIMARY KEY (id),
    CONSTRAINT fk_auxiliary_links_admin
        FOREIGN KEY (admin_id) REFERENCES admins(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 6. drivers
--    Cuentas de choferes. Tienen TTL (expires_at).
-- ============================================================
CREATE TABLE drivers (
    id                     VARCHAR(100)  NOT NULL,
    username               VARCHAR(160)  NOT NULL,
    password_hash          TEXT          NOT NULL,
    current_password       VARCHAR(255)  NOT NULL DEFAULT '',
    name                   VARCHAR(160)  NOT NULL DEFAULT '',
    created_by_admin_id    VARCHAR(100)  NULL DEFAULT NULL,
    created_by_admin_name  VARCHAR(160)  NOT NULL DEFAULT '',
    created_by_admin_email VARCHAR(160)  NOT NULL DEFAULT '',
    created_at             BIGINT        NOT NULL,
    expires_at             BIGINT        NULL DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_drivers_username (username),
    CONSTRAINT fk_drivers_admin
        FOREIGN KEY (created_by_admin_id) REFERENCES admins(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 7. events
--    FK a event_requests se agrega con ALTER TABLE al final
--    para evitar la referencia circular.
-- ============================================================
CREATE TABLE events (
    id                        VARCHAR(100)  NOT NULL,
    name                      VARCHAR(255)  NOT NULL DEFAULT '',
    code                      VARCHAR(80)   NOT NULL DEFAULT '',
    date                      DATE          NULL DEFAULT NULL,
    image                     VARCHAR(300)  NOT NULL DEFAULT '',
    thumbnail                 VARCHAR(300)  NOT NULL DEFAULT '',
    province                  VARCHAR(80)   NOT NULL DEFAULT '',
    reservation_link          VARCHAR(300)  NOT NULL DEFAULT '',
    whatsapp_number           VARCHAR(30)   NOT NULL DEFAULT '',
    departure_place           VARCHAR(255)  NOT NULL DEFAULT '',
    departure_time            VARCHAR(60)   NOT NULL DEFAULT '',
    return_time               VARCHAR(60)   NOT NULL DEFAULT '',
    payment_methods           TEXT          NOT NULL,
    departure_info            TEXT          NOT NULL,
    company_name              VARCHAR(120)  NOT NULL DEFAULT '',
    transfer_alias            VARCHAR(120)  NOT NULL DEFAULT '',
    transfer_cbu              VARCHAR(40)   NOT NULL DEFAULT '',
    transfer_banco            VARCHAR(120)  NOT NULL DEFAULT '',
    image_focus_x             INT           NOT NULL DEFAULT 50,
    image_focus_y             INT           NOT NULL DEFAULT 50,
    image_scale               INT           NOT NULL DEFAULT 100,
    transfer_account_info     TEXT          NOT NULL,
    payment_proof_destination VARCHAR(300)  NOT NULL DEFAULT '',
    post_payment_instructions TEXT          NOT NULL,
    refund_policy_notice      TEXT          NOT NULL,
    status                    VARCHAR(30)   NOT NULL DEFAULT 'en_curso',
    source_type               VARCHAR(50)   NOT NULL DEFAULT 'event',
    request_id                VARCHAR(100)  NULL DEFAULT NULL,
    is_visible_in_home        BOOLEAN       NOT NULL DEFAULT TRUE,
    is_paused_in_home         BOOLEAN       NOT NULL DEFAULT FALSE,
    is_deleted_from_home      BOOLEAN       NOT NULL DEFAULT FALSE,
    assigned_driver           VARCHAR(100)  NULL DEFAULT NULL,
    tracking_enabled          BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at                BIGINT        NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_events_driver
        FOREIGN KEY (assigned_driver) REFERENCES drivers(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 8. event_dates
--    Fechas adicionales de un evento (array dates[]).
-- ============================================================
CREATE TABLE event_dates (
    id       INT           NOT NULL AUTO_INCREMENT,
    event_id VARCHAR(100)  NOT NULL,
    date     DATE          NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_event_dates (event_id, date),
    CONSTRAINT fk_event_dates_event
        FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 9. event_departure_places
--    Lugares de salida del evento (array departurePlaces[]).
-- ============================================================
CREATE TABLE event_departure_places (
    id       INT           NOT NULL AUTO_INCREMENT,
    event_id VARCHAR(100)  NOT NULL,
    place    VARCHAR(255)  NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_departure_places_event
        FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 10. event_departure_times
--     Horarios de salida del evento (array departureTimes[]).
-- ============================================================
CREATE TABLE event_departure_times (
    id       INT           NOT NULL AUTO_INCREMENT,
    event_id VARCHAR(100)  NOT NULL,
    time     VARCHAR(60)   NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_departure_times_event
        FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 11. event_requests
--     Solicitudes de eventos creadas por admins de tipo
--     'event_creator' y revisadas por el admin principal.
-- ============================================================
CREATE TABLE event_requests (
    id                   VARCHAR(100)  NOT NULL,
    event_name           VARCHAR(120)  NOT NULL DEFAULT '',
    event_code           VARCHAR(60)   NOT NULL DEFAULT '',
    whatsapp_number      VARCHAR(30)   NOT NULL DEFAULT '',
    event_date           DATE          NULL DEFAULT NULL,
    image                VARCHAR(300)  NOT NULL DEFAULT '',
    thumbnail            VARCHAR(300)  NOT NULL DEFAULT '',
    province             VARCHAR(80)   NOT NULL DEFAULT '',
    departure_city       VARCHAR(80)   NOT NULL DEFAULT '',
    departure_point      VARCHAR(180)  NOT NULL DEFAULT '',
    activity_info        VARCHAR(220)  NOT NULL DEFAULT '',
    notes                VARCHAR(800)  NOT NULL DEFAULT '',
    company_name         VARCHAR(120)  NOT NULL DEFAULT '',
    created_by_admin_id  VARCHAR(100)  NULL DEFAULT NULL,
    created_by_name      VARCHAR(160)  NOT NULL DEFAULT '',
    created_by_email     VARCHAR(160)  NOT NULL DEFAULT '',
    status               VARCHAR(20)   NOT NULL DEFAULT 'pendiente',
    created_at           BIGINT        NOT NULL,
    reviewed_at          BIGINT        NULL DEFAULT NULL,
    reviewed_by_admin_id VARCHAR(100)  NULL DEFAULT NULL,
    reviewed_by_name     VARCHAR(160)  NOT NULL DEFAULT '',
    published_event_id   VARCHAR(100)  NULL DEFAULT NULL,
    PRIMARY KEY (id),
    CONSTRAINT chk_event_requests_status
        CHECK (status IN ('pendiente', 'aceptado', 'rechazado')),
    CONSTRAINT fk_event_requests_created_by
        FOREIGN KEY (created_by_admin_id) REFERENCES admins(id) ON DELETE SET NULL,
    CONSTRAINT fk_event_requests_reviewed_by
        FOREIGN KEY (reviewed_by_admin_id) REFERENCES admins(id) ON DELETE SET NULL,
    CONSTRAINT fk_event_requests_published_event
        FOREIGN KEY (published_event_id) REFERENCES events(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- FK circular: events.request_id → event_requests
ALTER TABLE events
    ADD CONSTRAINT fk_events_request_id
        FOREIGN KEY (request_id) REFERENCES event_requests(id) ON DELETE SET NULL;

-- ============================================================
-- 12. event_request_dates
--     Fechas propuestas en una solicitud (array eventDates[]).
-- ============================================================
CREATE TABLE event_request_dates (
    id               INT           NOT NULL AUTO_INCREMENT,
    event_request_id VARCHAR(100)  NOT NULL,
    date             DATE          NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_event_request_dates (event_request_id, date),
    CONSTRAINT fk_event_request_dates_request
        FOREIGN KEY (event_request_id) REFERENCES event_requests(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 13. event_request_departure_options
--     Opciones de punto de salida definidas en la solicitud.
-- ============================================================
CREATE TABLE event_request_departure_options (
    id               VARCHAR(100)  NOT NULL,
    event_request_id VARCHAR(100)  NOT NULL,
    city             VARCHAR(80)   NOT NULL DEFAULT '',
    point            VARCHAR(180)  NOT NULL DEFAULT '',
    single_day_only  BOOLEAN       NOT NULL DEFAULT FALSE,
    PRIMARY KEY (id),
    CONSTRAINT fk_departure_options_request
        FOREIGN KEY (event_request_id) REFERENCES event_requests(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 14. offered_date_groups
--     Grupos de fechas ofrecidas con precio y cuenta de pago.
--     Pertenecen a un evento O a una solicitud (nunca los dos).
-- ============================================================
CREATE TABLE offered_date_groups (
    id                             VARCHAR(100)  NOT NULL,
    event_id                       VARCHAR(100)  NULL DEFAULT NULL,
    event_request_id               VARCHAR(100)  NULL DEFAULT NULL,
    alias_name                     VARCHAR(80)   NOT NULL DEFAULT '',
    departure_option_id            VARCHAR(80)   NOT NULL DEFAULT '',
    province                       VARCHAR(80)   NOT NULL DEFAULT '',
    departure_city                 VARCHAR(80)   NOT NULL DEFAULT '',
    departure_point                VARCHAR(180)  NOT NULL DEFAULT '',
    service_amount                 VARCHAR(40)   NOT NULL DEFAULT '',
    single_day_only                BOOLEAN       NOT NULL DEFAULT FALSE,
    payment_mode                   VARCHAR(20)   NOT NULL DEFAULT '',
    selected_payment_account_id    VARCHAR(100)  NULL DEFAULT NULL,
    payment_account_alias_name     VARCHAR(80)   NOT NULL DEFAULT '',
    payment_account_owner_name     VARCHAR(120)  NOT NULL DEFAULT '',
    payment_account_bank           VARCHAR(120)  NOT NULL DEFAULT '',
    payment_account_cbu            VARCHAR(40)   NOT NULL DEFAULT '',
    payment_account_transfer_alias VARCHAR(120)  NOT NULL DEFAULT '',
    external_payment_link          VARCHAR(300)  NOT NULL DEFAULT '',
    activity_info                  VARCHAR(220)  NOT NULL DEFAULT '',
    PRIMARY KEY (id),
    CONSTRAINT chk_offered_parent CHECK (
        (event_id IS NOT NULL AND event_request_id IS NULL)
        OR (event_id IS NULL AND event_request_id IS NOT NULL)
    ),
    CONSTRAINT chk_offered_payment_mode
        CHECK (payment_mode IN ('transferencias', 'otros', 'ambos', '')),
    CONSTRAINT fk_offered_groups_event
        FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    CONSTRAINT fk_offered_groups_request
        FOREIGN KEY (event_request_id) REFERENCES event_requests(id) ON DELETE CASCADE,
    CONSTRAINT fk_offered_groups_payment_account
        FOREIGN KEY (selected_payment_account_id) REFERENCES company_payment_accounts(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 15. offered_date_group_dates
--     Fechas individuales dentro de un offered_date_group.
-- ============================================================
CREATE TABLE offered_date_group_dates (
    id                    INT           NOT NULL AUTO_INCREMENT,
    offered_date_group_id VARCHAR(100)  NOT NULL,
    date                  DATE          NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_offered_group_dates (offered_date_group_id, date),
    CONSTRAINT fk_offered_group_dates_group
        FOREIGN KEY (offered_date_group_id) REFERENCES offered_date_groups(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 16. locations
--     Puntos GPS enviados por el chofer durante un traslado.
-- ============================================================
CREATE TABLE locations (
    id         VARCHAR(100)  NOT NULL,
    event_id   VARCHAR(100)  NOT NULL,
    lat        DOUBLE        NOT NULL,
    lng        DOUBLE        NOT NULL,
    timestamp  BIGINT        NOT NULL,
    point_type VARCHAR(30)   NULL DEFAULT NULL,          -- NULL = punto normal, 'arrival'
    PRIMARY KEY (id),
    CONSTRAINT fk_locations_event
        FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 17. meeting_points
--     Último punto de encuentro conocido por evento (1:1).
-- ============================================================
CREATE TABLE meeting_points (
    id         INT           NOT NULL AUTO_INCREMENT,
    event_id   VARCHAR(100)  NOT NULL,
    lat        DOUBLE        NOT NULL,
    lng        DOUBLE        NOT NULL,
    timestamp  BIGINT        NOT NULL,
    point_type VARCHAR(30)   NOT NULL DEFAULT 'meeting_point',
    PRIMARY KEY (id),
    UNIQUE KEY uk_meeting_points_event (event_id),
    CONSTRAINT fk_meeting_points_event
        FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 18. comments
--     Comentarios de clientes sobre eventos ya finalizados.
-- ============================================================
CREATE TABLE comments (
    id            VARCHAR(100)  NOT NULL,
    event_id      VARCHAR(100)  NOT NULL,
    event_name    VARCHAR(255)  NOT NULL DEFAULT '',
    event_date    DATE          NULL DEFAULT NULL,
    customer_name VARCHAR(80)   NOT NULL DEFAULT '',
    comment       TEXT          NOT NULL,
    rating        SMALLINT      NULL DEFAULT NULL,
    created_at    BIGINT        NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT chk_comments_rating
        CHECK (rating >= 1 AND rating <= 5),
    CONSTRAINT fk_comments_event
        FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 19. driver_arrival_notifications
--     Notificaciones generadas cuando un chofer marca llegada.
-- ============================================================
CREATE TABLE driver_arrival_notifications (
    id                   VARCHAR(100)  NOT NULL,
    recipient_admin_id   VARCHAR(100)  NULL DEFAULT NULL,
    recipient_admin_name VARCHAR(160)  NOT NULL DEFAULT '',
    driver_id            VARCHAR(100)  NULL DEFAULT NULL,
    driver_name          VARCHAR(160)  NOT NULL DEFAULT '',
    driver_username      VARCHAR(160)  NOT NULL DEFAULT '',
    event_id             VARCHAR(100)  NULL DEFAULT NULL,
    event_name           VARCHAR(255)  NOT NULL DEFAULT '',
    lat                  DOUBLE        NOT NULL DEFAULT 0,
    lng                  DOUBLE        NOT NULL DEFAULT 0,
    created_at           BIGINT        NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_arrival_notif_admin
        FOREIGN KEY (recipient_admin_id) REFERENCES admins(id) ON DELETE SET NULL,
    CONSTRAINT fk_arrival_notif_driver
        FOREIGN KEY (driver_id) REFERENCES drivers(id) ON DELETE SET NULL,
    CONSTRAINT fk_arrival_notif_event
        FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 20. settings
--     Configuración global de la aplicación (fila única).
-- ============================================================
CREATE TABLE settings (
    id                                INT           NOT NULL AUTO_INCREMENT,
    logo                              VARCHAR(300)  NOT NULL DEFAULT '',
    driver_login_logo                 VARCHAR(300)  NOT NULL DEFAULT '',
    instagram_link                    VARCHAR(300)  NOT NULL DEFAULT '',
    default_whatsapp_number           VARCHAR(30)   NOT NULL DEFAULT '',
    default_reservation_link          VARCHAR(300)  NOT NULL DEFAULT '',
    default_payment_methods           VARCHAR(255)  NOT NULL DEFAULT '',
    default_transfer_alias            VARCHAR(120)  NOT NULL DEFAULT '',
    default_transfer_cbu              VARCHAR(40)   NOT NULL DEFAULT '',
    default_transfer_banco            VARCHAR(120)  NOT NULL DEFAULT '',
    default_transfer_amount           DECIMAL(12,2) NOT NULL DEFAULT 0,
    default_payment_proof_destination VARCHAR(300)  NOT NULL DEFAULT '',
    default_post_payment_instructions TEXT          NOT NULL,
    default_refund_policy_notice      TEXT          NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================================
-- 21. settings_social_links
--     Redes sociales de la organización (array socialLinks[]).
-- ============================================================
CREATE TABLE settings_social_links (
    id         INT           NOT NULL AUTO_INCREMENT,
    name       VARCHAR(60)   NOT NULL,
    icon       VARCHAR(60)   NOT NULL,
    link       VARCHAR(300)  NOT NULL,
    sort_order INT           NOT NULL DEFAULT 0,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 22. settings_reservation_whatsapp_numbers
--     Números de WhatsApp para reservas frecuentes.
-- ============================================================
CREATE TABLE settings_reservation_whatsapp_numbers (
    id           INT           NOT NULL AUTO_INCREMENT,
    phone_number VARCHAR(30)   NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_reservation_whatsapp (phone_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 23. settings_payment_methods
--     Formas de pago frecuentes reutilizables.
-- ============================================================
CREATE TABLE settings_payment_methods (
    id     INT           NOT NULL AUTO_INCREMENT,
    method VARCHAR(80)   NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_settings_payment_method (method)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 24. settings_transfer_account_presets
--     Cuentas de transferencia frecuentes.
-- ============================================================
CREATE TABLE settings_transfer_account_presets (
    id    INT           NOT NULL AUTO_INCREMENT,
    alias VARCHAR(120)  NOT NULL DEFAULT '',
    cbu   VARCHAR(40)   NOT NULL DEFAULT '',
    banco VARCHAR(120)  NOT NULL DEFAULT '',
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 25. settings_transfer_whatsapp_numbers
--     Números de WhatsApp para reporte de transferencias.
-- ============================================================
CREATE TABLE settings_transfer_whatsapp_numbers (
    id           INT           NOT NULL AUTO_INCREMENT,
    phone_number VARCHAR(30)   NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_transfer_whatsapp (phone_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- Índices para consultas frecuentes
-- ============================================================
CREATE INDEX idx_events_status              ON events(status);
CREATE INDEX idx_events_date               ON events(date);
CREATE INDEX idx_events_home_visibility    ON events(is_visible_in_home, is_paused_in_home, is_deleted_from_home);
CREATE INDEX idx_events_source_type        ON events(source_type);
CREATE INDEX idx_events_company_name       ON events(company_name);
CREATE INDEX idx_event_requests_status     ON event_requests(status);
CREATE INDEX idx_event_requests_created_by ON event_requests(created_by_admin_id);
CREATE INDEX idx_offered_groups_event      ON offered_date_groups(event_id);
CREATE INDEX idx_offered_groups_request    ON offered_date_groups(event_request_id);
CREATE INDEX idx_locations_event_id        ON locations(event_id);
CREATE INDEX idx_locations_timestamp       ON locations(event_id, timestamp);
CREATE INDEX idx_comments_event_id         ON comments(event_id);
CREATE INDEX idx_comments_created_at       ON comments(created_at DESC);
CREATE INDEX idx_driver_notif_admin        ON driver_arrival_notifications(recipient_admin_id);
CREATE INDEX idx_drivers_expires_at        ON drivers(expires_at);
