# ReviewsSystem

## Running the project
Just cd in the folder and run `docker-compose up`. A golang webserver and a postgres instance will be started. The frontend is served from the golang server.

### Email & Facebook
If you want to enable email confirmation and Facebook login, you need to specify `FACEBOOK_CLIENT_ID`, `FACEBOOK_CLIENT_SECRET`, `EMAIL_SMTP_USERNAME`, and `EMAIL_SMTP_PASSWORD` in the `docker-compose.yml`.

If you want to skip email confirmation for development purposes, change `SKIP_EMAIL_VERIFICATION` in `docker-compose.yml` to `true`.
   