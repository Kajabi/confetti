# Go Confetti

These instructions assume you've installed the [Kajabi CLI](https://kajabi-platform.stoplight.io/docs/developer-platform/branches/main/ZG9jOjQ3MjM3MTgy-the-kajabi-cli) and have a Kajabi account/site. (If you don't have a Kajabi account, reach out to us
at developer-platform@kajabi.com.)

1. Clone this Repository: `git clone git@github.com:Kajabi/confetti.git`
2. Change to this directory: `cd ./go`
3. Copy .env_sample to .env: `cp ../.env_sample .env`
4. Log in to the Kajabi CLI: `kajabi login`
5. Create a Kajabi app: `kajabi app create`
6. Copy the Client ID and Client Secret into the .env file.  **This is the only time the Client Secret will be displayed—be sure to save it!**
7. Run the app: `make run`
8. Install your new app: `kajabi app install`
