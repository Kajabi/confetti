# confetti-py

Confetti is a sample Kajabi app to test out the Kajabi Developer Platform.

To get started:

1. Clone this Repository: `git clone git@github.com:Kajabi/confetti.git`
1. Change directory to the project directory: `cd py_flask`
1. Copy the .env_sample to .env: `cp .env_sample .env`
1. Log in to the [Kajabi CLI](https://kajabi-platform.stoplight.io/docs/developer-platform/branches/main/ZG9jOjQ3MjM3MTgy-the-kajabi-cli): `kajabi login`
1. Create a Kajabi app: `kajabi app create`
1. Copy the Client ID and Client Secret into the .env file
1. Run the app (also creates virtual environment and installs dependencies): `make run`
1. Install your new app: `kajabi app install`
