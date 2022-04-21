from os import environ as env
from urllib.parse import urlencode, urlparse, parse_qs

from authlib.integrations.flask_client import OAuth
from authlib.integrations.requests_client import OAuth2Session
from flask import Flask, redirect, render_template, request, session, url_for
from flask_wtf import FlaskForm
from wtforms import BooleanField, IntegerField


app = Flask(__name__)
app.secret_key = env.get("FLASK_SECRET_KEY")

scopes = env.get("SCOPES").replace(",", " ")
SCRIPT_URL = "https://s3.amazonaws.com/sandbox-integrations-scripts-development.kajabi.dev/scripts/confetti.js"


oauth = OAuth(app)

client = OAuth2Session(
    client_id=env.get("CLIENT_ID"),
    client_secret=env.get("CLIENT_SECRET"),
    scope=env.get("SCOPES"),
)

oauth.register(
    "kajabi",
    client_id=env.get("CLIENT_ID"),
    client_secret=env.get("CLIENT_SECRET"),
    access_token_url=f"{env.get('AUTH_DOMAIN')}/oauth/token",
    access_token_params=None,
    authorize_url=f"{env.get('AUTH_DOMAIN')}/authorize",
    client_kwargs={"scope": scopes},
)


class ConfettiForm(FlaskForm):
    enabled = BooleanField(default=False)
    max = IntegerField(default=80)
    size = IntegerField(default=1)
    speed = IntegerField(default=25)


@app.route("/authorize")
def authorize():
    site_id = request.args.get("site_id")

    if not site_id:
        return "Site id must be provided as a query parameter", 400

    url = url_for(
        "callback",
        _external=True,
        site_id=site_id,
    )

    # NOTE Ensure your authorize URL is passing audience and scope, otherwise you will get an opaque token back :-(
    return oauth.kajabi.authorize_redirect(url, audience=env.get("AUTH_AUDIENCE"))


@app.route("/callback")
def callback():
    token = oauth.kajabi.authorize_access_token()
    session["user"] = token

    return redirect(url_for("site", site_id=request.args.get("site_id")))


@app.route("/sites/<int:site_id>", methods=["GET", "POST"])
def site(site_id):
    # Do we have a token and is it still valid?
    if "user" not in session:
        redirect(url_for("authorize", site_id=site_id))

    token = session["user"]
    form = ConfettiForm()

    # POST
    if form.validate_on_submit():
        query_params = form.data.copy()
        query_params.pop("csrf_token")
        url = SCRIPT_URL + f"?{urlencode(query_params)}"
        print(url)

        # create_or_update_script_tag
        script_tags = get_script_tags(site_id, token)
        print(script_tags)

        if not script_tags:
            print(create_script_tag(site_id, token, url))
        else:
            print(update_script_tag(site_id, script_tags[0]["id"], token, url))

        return render_template("edit.html", form=form)

    # GET
    script_tags = get_script_tags(site_id, token)
    print(script_tags)

    if script_tags:
        script_tag = script_tags[0]
        query_string = urlparse(script_tag["source_url"])[4]
        parsed = parse_qs(query_string)  # Parsed values are lists
        data = {
            "enabled": True
            if parsed.get("enabled", "")[0].lower() == "true"
            else False,
            "max": int(parsed.get("max", 80)[0]),
            "size": int(parsed.get("size", 1)[0]),
            "speed": int(parsed.get("speed", 25)[0]),
        }
        print(data)
        form = ConfettiForm(**data)

    return render_template("edit.html", form=form)


def get_script_tags(site_id, token):
    rsp = oauth.kajabi.get(
        f"{env.get('API_DOMAIN')}/sites/{site_id}/script_tags", token=token
    )

    return rsp.json()["script_tags"]


def create_script_tag(site_id, token, url):
    rsp = oauth.kajabi.post(
        f"{env.get('API_DOMAIN')}/sites/{site_id}/script_tags",
        token=token,
        json={"source_url": url},
    )

    return rsp.json()


def update_script_tag(site_id, script_tag_id, token, url):
    rsp = oauth.kajabi.put(
        f"{env.get('API_DOMAIN')}/sites/{site_id}/script_tags/{script_tag_id}",
        token=token,
        json={"source_url": url},
    )

    return rsp.json()
