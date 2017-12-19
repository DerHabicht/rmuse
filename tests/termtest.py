import requests

api = {
    "login": "/api/v1/login",
    "user":  "/api/v1/user"
}

u = {
    "email": "meh@meow.com",
    "password": "badpassword"
}

def post(url, endpoint, payload):
    global api
    return requests.post(f"{url}{api[endpoint]}", json=payload)
