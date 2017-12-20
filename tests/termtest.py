import requests

api = {
    "login": "/api/1/login",
    "user":  "/api/1/user"
}

u = {
    "email": "meh@meow.com",
    "password": "badpassword"
}

def post(url, endpoint, payload):
    global api
    return requests.post(f"{url}{api[endpoint]}", json=payload)
