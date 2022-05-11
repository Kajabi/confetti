# confetti-py

Confetti is a sample Kajabi app to test out the Kajabi Developer Platform.

These instructions assume you've installed the [Kajabi CLI](https://kajabi-platform.stoplight.io/docs/developer-platform/branches/main/ZG9jOjQ3MjM3MTgy-the-kajabi-cli),
have a Kajabi account/site, and have enabled script tags on your site.

(If you don't have a Kajabi account, reach out to us at developer-platform@kajabi.com.)

Full details can be found in the
[Getting Started](https://github.com/Kajabi/confetti#getting-started)
section in this repo's README.

To get started:

1. Clone this Repository: `git clone git@github.com:Kajabi/confetti.git`
2. Change directory to the project directory: `cd py_flask`
3. Copy the .env_sample to .env: `cp ../.env_sample .env`
4. Add a secret key to .env: `echo "FLASK_SECRET_KEY=$(xxd -l 16 -p < /dev/random)" >> .env`
5. Log in to the [Kajabi CLI](https://kajabi-platform.stoplight.io/docs/developer-platform/branches/main/ZG9jOjQ3MjM3MTgy-the-kajabi-cli): `kajabi login`
6. Create a Kajabi app: `kajabi app create`
7. Copy the Client ID and Client Secret into the .env file
8. Run the app (also creates virtual environment and installs dependencies): `make run`
9. Install your new app: `kajabi app install`
