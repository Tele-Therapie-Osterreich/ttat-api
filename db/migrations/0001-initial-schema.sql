-- +migrate Up

-- QUESTION: WHAT'S THE RIGHT WAY TO HANDLE DATABASE UPDATES AND
-- HISTORY HERE? SHOULD HISTORY BE ELIDED, I.E. THE DATA THAT'S STORED
-- IN THE DATABASE IS ALLOWED ONLY TO MIRROR EXACTLY WHAT THE USER
-- SAYS, NOT TO STORE HISTORICAL (AND NO LONGER CORRECT) DATA.
--
-- MY THINKING ABOUT THIS IS THAT IF YOU WANT TO BE STRICT ABOUT DATA
-- PROTECTION, THEN YOU DO NOT STORE *ANY* HISTORICAL DATA BEYOND
-- WHAT'S EXPLICITLY REQUIRED TO PROVIDE THE FUNCTIONALITY OF THE SITE
-- TO ITS USERS. IN THIS CONTEXT, THAT MEANS STORING NO HISTORICAL
-- INFORMATION AT ALL.

-- Top-level therapist types.
CREATE TYPE therapist_type AS ENUM ('unknown', 'ergo', 'physio', 'logo');

-- Approval status for therapist accounts.
CREATE TYPE approval_state AS ENUM ('new', 'approved', 'edits_pending', 'suspended');

-- One record per therapist user. Newly created therapist profiles are
-- not visible in search results until they are approved, which
-- generally requires that a minimal set of information is filled in
-- (to be really minimal, that would just be email, name and therapist
-- type).
--
-- The administrative interface shows new profiles as "pending
-- approval" once the minimal set of data is included.
--
-- TODO: ADD KK CONTRACT STATUS
-- TODO: ADD LAT/LON LOCATION DATA, GEOCODED AT SETUP TIME?
CREATE TABLE therapists (
  id               SERIAL PRIMARY KEY,
  email            TEXT UNIQUE NOT NULL,
  type             therapist_type NOT NULL DEFAULT 'unknown',
  name             TEXT,
  street_address   TEXT,
  city             TEXT,
  postcode         TEXT,
  country          TEXT,
  phone            TEXT,
  website          TEXT,
  languages        TEXT[],
  short_profile    TEXT,
  full_profile     TEXT,
  status           approval_state DEFAULT 'new',
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- When a therapist edits their profile, those edits are not
-- immediately reflected in the public view of their profile. Instead
-- the edits are held in a row of this table as a JSON-based
-- field-level patch between the public view and the newly edited
-- version of the profile. When the edits are approved by an
-- administrator, the patch is folded into the profile in the main
-- therapists table. Any extra edits before approval are folded into
-- the patch in this table so that edits can be approved all at once.
CREATE TABLE therapist_pending_edits (
  id            SERIAL PRIMARY KEY,
  therapist_id  INTEGER UNIQUE NOT NULL REFERENCES therapists(id) ON DELETE CASCADE,
  patch         JSONB,
  edited_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Auxiliary table storing therapist profile images.
-- TODO: CHANGE THIS TO USE AN EXTERNAL IMAGE MANAGEMENT SERVICE
CREATE TABLE images (
  id            SERIAL PRIMARY KEY,
  therapist_id  INTEGER UNIQUE NOT NULL REFERENCES therapists(id) ON DELETE CASCADE,
  extension     TEXT NOT NULL,
  data          BYTEA
);

-- -- SELECT id, email, name,
-- --        street_address, city, postcode, country, phone,
-- --        photo, short_profile FROM therapists WHERE id IN
-- --  (SELECT therapist_id FROM therapist_specialities WHERE speciality_id = ?);

-- Specialities for each therapist type.
--
-- The optional icon field gives the filename of an icon in site
-- static storage, accessed as
-- https://www.tele-therapie.at/images/<icon>.
CREATE TABLE specialities (
  id       SERIAL PRIMARY KEY,
  type     therapist_type,
  label    TEXT,
  icon     TEXT,

  UNIQUE(type, label)
);

-- Join table for therapist therapist/speciality many-to-many relation.
CREATE TABLE therapist_specialities (
  id             SERIAL PRIMARY KEY,
  therapist_id   INTEGER REFERENCES therapists(id) ON DELETE CASCADE,
  speciality_id  INTEGER REFERENCES specialities(id) ON DELETE CASCADE,
  index          INTEGER NOT NULL
);


-- "Magic link" login tokens.
CREATE TABLE login_tokens (
  token       CHAR(6)      PRIMARY KEY,
  email       TEXT         UNIQUE NOT NULL,
  language    TEXT,
  expires_at  TIMESTAMPTZ  NOT NULL
);
CREATE INDEX login_tokens_expired_index ON login_tokens(expires_at);

-- Session IDs for logins.
CREATE TABLE sessions (
  token    TEXT PRIMARY KEY,
  therapist_id  INTEGER REFERENCES therapists(id) ON DELETE CASCADE
);


-- Daily statistics about therapist creation and deletion.
CREATE TABLE therapist_statistics (
  id                  SERIAL PRIMARY KEY,
  date                DATE UNIQUE NOT NULL DEFAULT current_date,
  new_therapists      INTEGER DEFAULT 0,
  deleted_therapists  INTEGER DEFAULT 0
);

-- Daily search statistics.
CREATE TABLE search_statistics (
  id          SERIAL PRIMARY KEY,
  date        DATE NOT NULL DEFAULT current_date,
  speciality  INTEGER NOT NULL REFERENCES specialities(id) ON DELETE RESTRICT,
  searches    INTEGER DEFAULT 0,

  UNIQUE (date, speciality)
);


-- +migrate Down

DROP TABLE search_statistics;
DROP TABLE therapist_statistics;
DROP TABLE sessions;
DROP TABLE login_tokens;
DROP TABLE therapist_specialities;
DROP TABLE specialities;
DROP TABLE therapist_pending_edits;
DROP TABLE images;
DROP TABLE therapists;
DROP TYPE approval_state;
DROP TYPE therapist_type;
