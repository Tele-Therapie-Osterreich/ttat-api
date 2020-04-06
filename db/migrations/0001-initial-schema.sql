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

-- Top-level therapist user types.
CREATE TYPE therapist_type AS ENUM ('ot', 'physio', 'speech');

-- Approval status for therapist accounts.
CREATE TYPE approval_state AS ENUM ('new', 'approved', 'edits_pending', 'suspended');

-- One record per therapist user. Newly created user profiles are not
-- visible in search results until they are approved, which generally
-- requires that a minimal set of information is filled in (to be
-- really minimal, that would just be email, name and therapist type).
--
-- The administrative interface shows new profiles as "pending
-- approval" once the minimal set of data is included.
--
-- TODO: ADD KK CONTRACT STATUS
-- TODO: ADD LAT/LON LOCATION DATA, GEOCODED AT SETUP TIME?
CREATE TABLE users (
  id               SERIAL PRIMARY KEY,
  email            TEXT UNIQUE NOT NULL,
  type             therapist_type NOT NULL DEFAULT 'ot',
  name             TEXT,
  street_address   TEXT,
  city             TEXT,
  postcode         TEXT,
  country          TEXT,
  phone            TEXT,
  short_profile    TEXT,
  full_profile     TEXT,
  status           approval_state DEFAULT 'new',
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- When a user edits their profile, those edits are not immediately
-- reflected in the public view of their profile. Instead the edits
-- are held in a row of this table as a JSON-based field-level patch
-- between the public view and the newly edited version of the
-- profile. When the edits are approved by an administrator, the patch
-- is folded into the profile in the main users table. Any extra edits
-- before approval are folded into the patch in this table so that
-- edits can be approved all at once.
CREATE TABLE user_pending_edits (
  id         SERIAL PRIMARY KEY,
  user_id    INTEGER UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  patch      JSONB,
  edited_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Auxiliary table storing user profile images.
-- TODO: CHANGE THIS TO USE AN EXTERNAL IMAGE MANAGEMENT SERVICE
CREATE TABLE images (
  id         SERIAL PRIMARY KEY,
  user_id    INTEGER UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  extension  TEXT NOT NULL,
  data       BYTEA
);

-- -- SELECT id, email, name,
-- --        street_address, city, postcode, country, phone,
-- --        photo, short_profile FROM users WHERE id IN
-- --  (SELECT user_id FROM user_sub_specialities WHERE sub_speciality_id = ?);

-- Sub-specialities for each therapist type.
CREATE TABLE sub_specialities (
  id          SERIAL PRIMARY KEY,
  speciality  therapist_type,
  name        TEXT,
  label       TEXT,
  icon        INTEGER REFERENCES images(id) ON DELETE SET NULL,
  deleted     BOOLEAN DEFAULT false,

  UNIQUE(speciality, name)
);

-- -- Join table for therapist user/sub-speciality many-to-many relation.
CREATE TABLE user_sub_specialities (
  id                 SERIAL PRIMARY KEY,
  user_id            INTEGER REFERENCES users(id) ON DELETE CASCADE,
  sub_speciality_id  INTEGER REFERENCES sub_specialities(id) ON DELETE RESTRICT
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
  user_id  INTEGER REFERENCES users(id) ON DELETE CASCADE
);


-- Daily statistics about user creation and deletion.
CREATE TABLE user_statistics (
  id             SERIAL PRIMARY KEY,
  date           DATE UNIQUE NOT NULL DEFAULT current_date,
  new_users      INTEGER DEFAULT 0,
  deleted_users  INTEGER DEFAULT 0
);

-- Daily search statistics.
CREATE TABLE search_statistics (
  id              SERIAL PRIMARY KEY,
  date            DATE NOT NULL DEFAULT current_date,
  sub_speciality  INTEGER NOT NULL REFERENCES sub_specialities(id) ON DELETE RESTRICT,
  searches        INTEGER DEFAULT 0,

  UNIQUE (date, sub_speciality)
);


-- +migrate Down

DROP TABLE search_statistics;
DROP TABLE user_statistics;
DROP TABLE sessions;
DROP TABLE login_tokens;
DROP TABLE user_sub_specialities;
DROP TABLE sub_specialities;
DROP TABLE user_pending_edits;
DROP TABLE users;
DROP TABLE images;
DROP TYPE approval_state;
DROP TYPE therapist_type;
