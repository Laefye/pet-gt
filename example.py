import requests
from rich import print
from pydantic import BaseModel
import time

class GameLoginResponse(BaseModel):
    id: str
    url: str
    token: str

class GameLoginUser(BaseModel):
    id: str
    username: str

class GameCode(BaseModel):
    id: str
    user: GameLoginUser

class GameLogin(BaseModel):
    id: str
    token: str
    user: GameLoginUser

class GameLoginStateResponse(BaseModel):
    id: str
    code: GameCode | None

BASE_URL = "http://localhost:8080"

def create_game_login_request() -> GameLoginResponse:
    response = requests.post(f"{BASE_URL}/api/game/login")
    response.raise_for_status()
    return GameLoginResponse(**response.json())

def wait_for_user_login(id: str, token: str) -> GameLoginStateResponse:
    while True:
        params = {
            "id": id,
            "token": token
        }
        response = requests.get(f"{BASE_URL}/api/game/login", params=params)
        response.raise_for_status()
        state = GameLoginStateResponse(**response.json())
        if state.code is not None:
            return state
        time.sleep(5)

def exchange_code_for_token(code_id: str, user_id: str) -> GameLogin:
    params = {
        "code_id": code_id,
        "user_id": user_id
    }
    response = requests.get(f"{BASE_URL}/api/game/exchange", params=params)
    response.raise_for_status()
    return GameLogin(**response.json())

if __name__ == "__main__":
    login_response = create_game_login_request()
    print(login_response)
    print(f"Visit the following URL to log in: {login_response.url}")
    state = wait_for_user_login(login_response.id, login_response.token)
    print(state)
    game_login = exchange_code_for_token(state.code.id, state.code.user.id)
    print(game_login)