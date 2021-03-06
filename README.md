# Confetti

Confetti is a simple app to demonstrate the Kajabi Developer Platform.

The Kajabi Developer Platform allows third-parties to extend the functionality
of the core Kajabi product. Apps are granted access to Kajabi via OAuth, and they interact with Kajabi using a [RESTful API](https://kajabi-platform.stoplight.io/docs/developer-platform).

A [Kajabi Hero](https://kajabi.com/hero) can install apps to their
website using the CLI or the GUI interface in their site settings.

## Getting Started

You can explore how to set up apps using the samples in this repo.
Installation instructions for each language/platform are in their respective
directories.

You'll need a Kajabi site; if you don't have one, reach out to us at
[developer-platform@kajabi.com](mailto:developer-platform@kajabi.com).

You'll need to enable script tags on your Kajabi site. On the Kajabi admin page
for your site, go to Settings > Site Details and scroll down to the Page Scripts
section.

![Kajabi Admin Panel](/assets/kajabi-admin-page-scripts.png)

Add the following snippet to the Page Scripts field:

```js
<script src="https://scripts.kajabi.dev/{YOUR_SITE_ID}/scripts.js"></script>
```

Your site id is the number after `/admin/sites/` in the URL of the admin page.

You'll also need to have the [Kajabi CLI](https://kajabi-platform.stoplight.io/docs/developer-platform/ZG9jOjQ3MjM3MTgy-the-kajabi-cli)
installed on your development machine to set up an app.

1. Clone this repository: `git clone git@github.com:Kajabi/confetti.git`
1. Copy `.sample_env` to your preferred platform's directory, e.g. `cp .env_sample ./go/.env`
1. Follow the installation instructions in the directory for your preferred language/framework.

## Adding Another Piece of Confetti to the Pile

If you'd like to add a new version of Confetti to the repo, create a new directory
with the following `README.md` template. The newly-created directory should use
the naming convention `language_framework`.

Care should be taken to implement the app using idioms common to the language
and framework you're using. When serving your application locally, use port 3000.

```
# {{LANGUAGE OR FRAMEWORK}} Confetti

Confetti is a sample Kajabi app to test out the Kajabi Developer Platform.

These instructions assume you've installed the [Kajabi CLI](https://kajabi-platform.stoplight.io/docs/developer-platform/branches/main/ZG9jOjQ3MjM3MTgy-the-kajabi-cli), have a Kajabi account/site, and have enabled script tags on your site.

(If you don't have a Kajabi account, reach out to us at developer-platform@kajabi.com.)

Full details can be found in the
[Getting Started](https://github.com/Kajabi/confetti#getting-started)
section in this repo's README.

1. Clone this Repository: `git clone git@github.com:Kajabi/confetti.git`
2. Change to this directory: `cd ./{{DIR}}`
3. Copy .env_sample to .env: `cp ../.env_sample .env`
4. Log in to the Kajabi CLI: `kajabi login`
5. Create a Kajabi app: `kajabi app create`
6. Copy the Client ID and Client Secret into the .env file.  **This is the only time the Client Secret will be displayed???be sure to save it!**
7. Run the app: `{{COMMAND TO RUN APP LOCALLY}}`
8. Install your new app: `kajabi app install`
```

After installing the app, the Confetti admin panel will open in your default browser.

![Confetti Admin Panel](/assets/confetti-admin.png)

Click the Enable toggle, and then click the Save button, to enable confetti on your site. Party time!

### Confetti Specification

The Confetti app is granted access to the user's site with OAuth, and must implement
the endpoints below.

- You should provide helper functions to set up/run your application (assume a
Unix-like environment)
- When running the app locally, use port 3000
- Copy `assets/templates/edit.html` to your project and modify it for your templating language
- Do NOT copy `assets/static` to your project???create a symlink to this folder instead

#### `GET /authorize`

When the app is installed via Kajabi's site settings or the CLI, this is the
endpoint that is called. The site's id will be apppended to the URL as a query
parameter.

This endpoint must create a valid OAuth request/URL and redirect to the auth provider.

The site's id must be included in the request.

When the user grants the app access, the auth provider will redirect to the following URL.

#### `GET /callback`

This endpoint validates the response, persists the token in a cookie/session, and redirects to the following URL.

#### `GET /sites/{site_id}`

If a request to this endpoint does not have a token present, redirect to `/authorize`.

The Confetti app is loaded onto the user's site with a script tag. The app's state
is stored as query parameters on the script's URL.

Check to see if there's already a script URL for this site, and use the values
from the query parameters as the initial values for the form.

If no URL exists yet for this site, use the following default values:

- `enabled` False
- `max` 80
- `size` 1
- `speed` 25

Render and return the form in `edit.html` with the values.

#### `POST /sites/{site_id}`

If a request to this endpoint does not have a token present, redirect to `/authorize`.

The following form data must be posted to this endpoint:

- `enabled` boolean
- `max` integer
- `size` integer
- `speed` integer

Create or update the script URL for this site with the provided values as
query parameters.

Render and return the form in `edit.html` with the new values.
