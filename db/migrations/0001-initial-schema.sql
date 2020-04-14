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
CREATE TYPE approval_state AS ENUM ('new', 'active', 'suspended');

-- One record per therapist user. Most profile information is stored
-- in the linked profiles table.
CREATE TABLE therapists (
  id             SERIAL PRIMARY KEY,
  email          TEXT UNIQUE NOT NULL,
  status         approval_state DEFAULT 'new',
  last_login_at  TIMESTAMPTZ DEFAULT now(),
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Profiles for therapists: one or two records per therapist user (a
-- maximum of one public and one pending).
--
-- Newly created therapist profiles are not visible in search results
-- until they are approved, which generally requires that a minimal
-- set of information is filled in (to be really minimal, that would
-- just be email, name and therapist type). Also, subsequent profile
-- edits are not made public until approved by an administrator.
--
-- The mechanism for handling this is to keep parallel public and
-- pending profiles (distinguished by the boolean public column). The
-- following possible states can exist (T = therapists, P = profiles):
--
-- T.status   P.public  Comment
--
-- new        F         Profiles of new users are not visible until
--                      approved.
--
-- active     T         Active user with on pending edits.
-- active     T+F       Active user with pending edits (only public
--                      profile is visible to other non-admin users).
--
-- suspended  F         Users suspended by an administrator do not
--                      appear in public profile listings.
--
-- The administrative interface shows new profiles as "pending
-- approval" once the minimal set of data is included.
--
-- Any user edits to their profile are applied to the non-public
-- (pending) profile version, creating a new pending profile if none
-- already exists for the user.
--
-- TODO: ADD KK CONTRACT STATUS
-- TODO: ADD LAT/LON LOCATION DATA, GEOCODED AT SETUP TIME?
CREATE TABLE profiles (
  id               SERIAL PRIMARY KEY,
  therapist_id     INTEGER NOT NULL REFERENCES therapists(id) ON DELETE CASCADE,
  public           BOOLEAN,
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
  edited_at        TIMESTAMPTZ NOT NULL DEFAULT now(),

  UNIQUE (therapist_id, public)
);

-- Auxiliary table storing therapist profile images.
-- TODO: CHANGE THIS TO USE AN EXTERNAL IMAGE MANAGEMENT SERVICE
CREATE TABLE images (
  id          SERIAL PRIMARY KEY,
  profile_id  INTEGER UNIQUE NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  extension   TEXT NOT NULL,
  data        BYTEA
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

-- Join table for therapist profile/speciality many-to-many relation.
CREATE TABLE therapist_specialities (
  id             SERIAL PRIMARY KEY,
  profile_id     INTEGER REFERENCES profiles(id) ON DELETE CASCADE,
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
DROP TABLE images;
DROP TABLE profiles;
DROP TABLE therapists;
DROP TYPE approval_state;
DROP TYPE therapist_type;
