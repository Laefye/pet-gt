package controllers

import (
	"encoding/json"
	"errors"
	"gt/internal/middleware"
	"gt/internal/services"
	"gt/internal/templates"
	"net/http"
	"net/url"
)

type GameController struct {
	gameService *services.GameService
	userService *services.UserService
}

func NewGameController(gameService *services.GameService, userService *services.UserService) *GameController {
	return &GameController{gameService: gameService, userService: userService}
}

func (c *GameController) jsonResponse(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (c *GameController) renderTemplate(w http.ResponseWriter, data *templates.GameData) {
	err := templates.GameLoginTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *GameController) renderGameOKTemplate(w http.ResponseWriter, data *templates.GameData) {
	err := templates.GameOKTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type gameLoginRequestResponse struct {
	ID    string `json:"id"`
	URL   string `json:"url,omitempty"`
	Token string `json:"token,omitempty"`
}

type gameErrorResponse struct {
	Message string `json:"message"`
}

type gameLoginState struct {
	ID     string  `json:"id"`
	UserID *string `json:"user_id"`
}

type gameUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type gameLogin struct {
	ID    string   `json:"id"`
	Token string   `json:"token"`
	User  gameUser `json:"user"`
}

func (c *GameController) CreateGameLoginRequest(w http.ResponseWriter, r *http.Request) {
	req, err := c.gameService.CreateGameLoginRequest(r.Context())
	if err != nil {
		c.jsonResponse(w, gameErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	u := url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   "/game",
	}
	query := u.Query()
	query.Set("id", req.GameLoginRequest.ID)
	u.RawQuery = query.Encode()
	c.jsonResponse(w, gameLoginRequestResponse{ID: req.GameLoginRequest.ID, URL: u.String(), Token: req.Token}, http.StatusCreated)
}

func (c *GameController) GetGameLoginPage(w http.ResponseWriter, r *http.Request) {
	requestID := r.URL.Query().Get("id")
	if requestID == "" {
		c.renderTemplate(w, &templates.GameData{Error: "Missing request ID"})
		return
	}
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		query := url.Values{}
		query.Set("redirect", LoginRedirectData{
			Action:             LoginActionGameLogin,
			GameLoginRequestID: requestID,
		}.ToQuery())
		http.Redirect(w, r, "/login?"+query.Encode(), http.StatusSeeOther)
		return
	}
	gameLoginRequest, err := c.gameService.GetGameLoginRequest(r.Context(), requestID)
	if err != nil {
		c.renderTemplate(w, &templates.GameData{Error: err.Error()})
		return
	}
	c.renderTemplate(w, &templates.GameData{GameLoginRequest: gameLoginRequest, User: user})
}

func (c *GameController) PostGameLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	requestID := r.FormValue("request_id")
	if requestID == "" {
		http.Error(w, "Missing game login request ID", http.StatusBadRequest)
		return
	}
	user := middleware.UserFromContext(r.Context())
	err := c.gameService.Login(r.Context(), requestID, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.renderGameOKTemplate(w, &templates.GameData{User: user})
}

func (c *GameController) GetGameLoginState(w http.ResponseWriter, r *http.Request) {
	requestID := r.URL.Query().Get("id")
	token := r.URL.Query().Get("token")
	if requestID == "" || token == "" {
		c.jsonResponse(w, gameErrorResponse{Message: "Missing request ID or token"}, http.StatusBadRequest)
		return
	}
	req, err := c.gameService.GetGameLoginRequestState(r.Context(), requestID, token)
	if err != nil {
		if errors.Is(err, services.ErrGameLoginRequestNotFound) {
			c.jsonResponse(w, gameErrorResponse{Message: "Game login request not found"}, http.StatusNotFound)
		} else {
			c.jsonResponse(w, gameErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		}
		return
	}
	response := gameLoginState{ID: req.ID, UserID: req.UserID}
	c.jsonResponse(w, response, http.StatusOK)
}

func (c *GameController) ExchangeGameLoginCode(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	token := r.URL.Query().Get("token")
	if id == "" || token == "" {
		c.jsonResponse(w, gameErrorResponse{Message: "Missing ID or token"}, http.StatusBadRequest)
		return
	}
	code, err := c.gameService.Exchange(r.Context(), id, token)
	if err != nil {
		if errors.Is(err, services.ErrGameLoginRequestNotFound) {
			c.jsonResponse(w, gameErrorResponse{Message: "Game login request not found"}, http.StatusNotFound)
		} else {
			c.jsonResponse(w, gameErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		}
		return
	}
	user, err := c.userService.GetUserByID(r.Context(), code.GameLogin.UserID)
	if err != nil {
		c.jsonResponse(w, gameErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	c.jsonResponse(w, gameLogin{
		ID:    code.GameLogin.ID,
		Token: code.Token,
		User: gameUser{
			ID:       user.ID,
			Username: user.Username,
		},
	}, http.StatusOK)
}
