# REST API backend for `tele-therapie.at`

## Routes

### Fundamental

`GET /`

 * Health check

### Authentication

`POST /auth/request-login-email`

 * Authentication (request login code via email)

`POST /auth/login`

 * Authentication (submit login code)

`POST /auth/logout`

 * Authentication (logout all devices)

### Therapist user profiles

```
GET /me         (authenticated only)
GET /me/public
GET /me/pending
```

 * Get profile of logged in user

`PATCH /me  (authenticated only)`

 * Edit profile of logged in user

`DELETE /me  (authenticated only)`

 * Delete logged in user

### Retrieve therapist information

`GET /therapist/:id`

 * Retrieve full profile of a therapist

```
GET /therapists
    ?type={ergo|physio|logo}   optional therapist type
    &speciality=:labels        optional comma-separated list of speciality labels
    &limit=:num
    &offset=:num
```

 * Retrieve summary profiles of a set of therapists

### Speciality management

`GET /specialities/{ergo|physio|logo}`

 * Retrieve speciality list for a therapist type, in the order that
   the specialities appear in the site

### Therapist avatar images

```
GET /image/:id.{jpg|png}
GET /image/:id.{jpg|png}?thumbnail=true
```

 * Retrieve avatar image for a therapist: these links are generated by
   the backend and passed in responses to the `/therapists` and
   `/therapist/:id` routes


## API usage

### Unauthenticated

`GET /specialities` can be used to retrieve the speciality list for a
therapist type, although it generally shouldn't be needed in the
front-end, since edits to the category list in the admin interface
will trigger changes to JSON files in the site repository used to
generate the front pages for different therapist types.

`GET /therapists` is the main route used for the "matchmaker"
interface.

`GET /therapist/:id` can then be used to retrieve the full user
profile for an individual therapist.

Once a list of therapist summary profiles or a full profile has been
retrieved, the avatar image (thumbnail or full version) can be
retrieved using the `GET /image/:id.{jpg|png}` route, following the
link included in the therapist profile data.

**Resulting requirements**

 * Category list updates in the admin interface should propagate the
   new category data to the main site by updating a
   `site/data/specialities.json` file in GitHub. This means that the
   admin server must be able to perform operations on the GitHub
   repository (possibly as the logged in admin user by assigning the
   right permissions via GitHub OAuth).

 * The summary profiles returned from the `GET /therapists` route
   should include the minimum set of information needed to render a
   list of therapists filterable by speciality, Krankenkasse contract
   status, language and location. This includes a link to access a
   thumbnail image for the therapist. In particular, the summary
   profiles shouldn't include the full therapist profile text (which
   could be quite long).

 * The full profile accessed via the `GET /therapist/:id` should
   include everything, including a link to access the full size
   version of the therapist's avatar image.

### Authenticated

Therapist users authenticate to the system using "magic link" login
via the `/auth/request-login-email` and `/auth/login` routes.

**QUESTION: how to do the token stuff for Vue.js? Use JWT? Or just the
low-rent approach I've used up to now?**

**QUESTION: should it also be possible for therapists to sign up using
social login (Google or Facebook probably)?**

Once authenticated, the `/me` routes can be accessed to view and
modify the therapist's profile. The therapist can see both their
current publicly visible profile and their profile with any pending
edits applied to it. Any edits made to a therapist's profile are saved
as pending, and need to be reviewed by an administrator before
becoming publicly visible. (One exception to this is if a therapist
wants to delete their profile: that happens immediately.)

The therapist profile includes contact data (address, phone, email,
website), an ordered list of specialities the therapist can offer, an
ordered list of languages that the therapist can offer, a short
profile to appear in therapist lists, and a full long-text profile to
appear on the individual therapist's page.

**QUESTION: What to use for editing the therapist profiles? Some sort
of WYSIWYG editor outputting Markdown or something? What about letting
therapists upload images? (Something like https://quilljs.com/? Or
CKEditor, which is free for open source projects.)**

All communications with therapists about the pending, approved or
rejects status of edits to their profile will be sent over email.

**Resulting requirements**

 * There needs to be a way for a therapist to preview the summary and
   full view of their profile before submitting edits. This mechanism
   should be self-contained in the front end, implemented using
   Vue.js. (This should be relatively easy to do, since the profiles
   seen in search results will be rendered from the same data as the
   therapist is editing.)

 * There needs to be a way to send email. What are the GDPR
   implications of using Mailjet? (Actually looks pretty good, at a
   first glance: servers in EU, GDPR compliance statements, etc.)

 * The therapist profile editing page needs the following features:
    - Non-intrusive address completion
    - A good phone number input
    - A reorderable list of specialities (one of those two-column
      things with dragging between the lists?)
    - A reorderable list of languages with a good selector for common
      languages
    - Some sort of rich text editor for the full profile, perhaps
      including image upload

To start with let's have a single profile image, plain ASCII for both
profiles. We can add better editing and representation later on.

## API behaviour

### Fundamental

```
GET /
  => 200 OK with "ok"
```

### Authentication

```
POST /auth/request-login-email { "email": "test@example.com" }
  => 204 No Content
     Generate login token
     Send email with login token to email address
```

```
POST /auth/login { "login_token": "dfkljwrj23fwkefdslkj4t" }
  => 400 Bad Request if token unknown
  => 200 OK + profile if valid token (/me)
```

```
POST /auth/logout
  => 404 Not Found if not authenticated
  => 204 No Content
```

### Therapist user profiles

```
GET /me
  => 404 Not Found if not authenticated
  => 200 OK
  {
    "id": 1234,
    "email": "ian@skybluetrades.net",
    "type": "ergo",
    "name": "Ian Ross",
    "street_address": "Drobollacher Seepromenade 17",
    "city": "Villach",
    "postcode": "9580",
    "country": "AT",
    "phone": "+436804451378",
    "website": "https://www.skybluetrades.net",
    "languages": ["en", "fr", "de"],
    "short_profile": "I'm totally not a therapist...",
    "full_profile": "I mean, really totally not.  Nothing to do with it, man.",
    "status": "approved",
    "photo": "svwe9fwvfj3.jpg",
    "specialities": ["Hand", "Neuro"]
  }
```

```
PATCH /me { "phone": "+436804451379" }
  => 404 Not Found if not authenticated
  => 200 OK returning full profile
```

(Also supports image upload using multipart encoding, which sets a new
image for the therapist)

```
DELETE /me
  => 404 Not Found if not authenticated
  => 204 No Content
```

### Retrieve therapist information

```
GET /therapist/1234
  => 404 Not Found if therapist ID not known
  => 200 OK with full profile
```

```
GET /therapists?type=ergo&speciality=Neuro&limit=5
  => 200 OK with summary profiles
  [
  { "id": 1234,
    "email": "ian@skybluetrades.net",
    "type": "ergo",
    "name": "Ian Ross",
    "city": "Villach",
    "postcode": "9580",
    "country": "AT",
    "phone": "+436804451378",
    "website": "https://www.skybluetrades.net",
    "languages": ["en", "fr", "de"],
    "short_profile": "I'm totally not a therapist...",
    "photo": "svwe9fwvfj3.jpg",
    "specialities": ["Hand", "Neuro"] },
    ...
  ]
```

### Speciality management

```
GET /specialities/ergo
  => 200 OK
  [ { "label": "Hand", "icon": "ergo-hand.png" },
    { "label": "Neuro", "icon": "ergo-neuro.png" },
    { "label": "Paediatric", "icon": "erog-paediatric.png" },
    { "label": "Psychiatric", "icon": "erog-psychiatric.png" },
    { "label": "Something else" },
    { "label": "Something else again" }
  ]
```

### Therapist avatar images

```
GET /image/rfgk43frf.jpg
  => 404 Not Found if not a known image ID
     200 OK with image data at size suitable for main therapist page
```

```
GET /image/rfgk43frf.jpg?thumbnail=true
  => 404 Not Found if not a known image ID
     200 OK with image data at size suitable for therapist list page
```
