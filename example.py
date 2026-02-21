import time
import requests
from pydantic import BaseModel
from rich.console import Console
from rich.table import Table
from rich.panel import Panel
from rich.live import Live
from rich.spinner import Spinner
from rich.align import Align

console = Console()

BASE_URL = "http://localhost:8080"


# ---------- MODELS ----------

class GameLoginResponse(BaseModel):
    id: str
    url: str
    token: str


class GameUser(BaseModel):
    id: str
    username: str
    email: str


class GameLogin(BaseModel):
    id: str
    token: str


class GameLoginStateResponse(BaseModel):
    id: str
    user_id: str | None


# ---------- API ----------

def create_game_login_request() -> GameLoginResponse:
    r = requests.post(f"{BASE_URL}/api/game/login")
    r.raise_for_status()
    return GameLoginResponse(**r.json())


def wait_for_user_login(id: str, token: str) -> GameLoginStateResponse:
    spinner = Spinner("dots", text="Ожидание входа пользователя…")

    with Live(
        Align.center(spinner, vertical="middle"),
        refresh_per_second=12,
        console=console,
    ):
        while True:
            r = requests.get(
                f"{BASE_URL}/api/game/login",
                params={"id": id, "token": token},
            )
            r.raise_for_status()
            state = GameLoginStateResponse(**r.json())

            if state.user_id is not None:
                return state

            time.sleep(2)


def exchange_code_for_token(id: str, token: str) -> GameLogin:
    r = requests.get(
        f"{BASE_URL}/api/game/exchange",
        params={"id": id, "token": token},
    )
    r.raise_for_status()
    return GameLogin(**r.json())


class GameAPI:
    def __init__(self, id: str, token: str):
        self.id = id
        self.token = token

    def get_user(self) -> GameUser:
        r = requests.get(
            f"{BASE_URL}/api/game/user",
            headers={
                "X-Game-Login-ID": self.id,
                "X-Game-Login-Token": self.token,
            },
        )
        r.raise_for_status()
        return GameUser(**r.json())


# ---------- UI HELPERS ----------

def show_login_info(login: GameLoginResponse):
    table = Table(title="🎮 Game Login", show_header=False)
    table.add_row("Login ID", f"[bold]{login.id}[/bold]")
    table.add_row("URL", f"[cyan]{login.url}[/cyan]")

    console.print(Panel(table, border_style="green"))


def show_user(user: GameUser):
    table = Table(show_header=False)
    table.add_row("ID", user.id)
    table.add_row("Username", f"[bold cyan]{user.username}[/bold cyan]")
    table.add_row("Email", f"[green]{user.email}[/green]")

    console.print(
        Panel(
            table,
            title="👤 User Profile",
            border_style="blue",
        )
    )


# ---------- MAIN ----------

if __name__ == "__main__":
    console.rule("[bold green]Game Login Flow")

    login = create_game_login_request()
    show_login_info(login)

    console.print(
        f"\n🔗 Откройте в браузере:\n[underline cyan]{login.url}[/underline cyan]\n"
    )

    state = wait_for_user_login(login.id, login.token)

    console.print("✅ Пользователь вошёл", style="bold green")

    game_login = exchange_code_for_token(state.id, login.token)

    api = GameAPI(game_login.id, game_login.token)
    user = api.get_user()

    show_user(user)