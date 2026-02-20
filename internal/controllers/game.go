package controllers

import (
	"encoding/json"
	"gt/internal/repository"
	"gt/internal/services"
	"gt/internal/templates"
	"net/http"
	"net/url"
)

type GameController struct {
	gameAPIService *services.GameAPIService
}

func NewGameController(gameAPIService *services.GameAPIService) *GameController {
	return &GameController{gameAPIService: gameAPIService}
}

func (c *GameController) jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *GameController) renderTemplate(w http.ResponseWriter, data *templates.GameData) {
	err := templates.GameTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type GameLoginRequestResponse struct {
	ID  string `json:"id"`
	Url string `json:"url,omitempty"`
}

type GameErrorResponse struct {
	Message string `json:"message"`
}

func (c *GameController) CreateGameLoginRequest(w http.ResponseWriter, r *http.Request) {
	gameLoginRequest, err := c.gameAPIService.CreateGameLoginRequest(r.Context())
	if err != nil {
		c.jsonResponse(w, GameErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	url := url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   "/game",
	}
	query := url.Query()
	query.Set("id", gameLoginRequest.ID)
	url.RawQuery = query.Encode()
	c.jsonResponse(w, GameLoginRequestResponse{ID: gameLoginRequest.ID, Url: url.String()}, http.StatusCreated)
}

func (c *GameController) GetGameLoginPage(user *repository.User, w http.ResponseWriter, r *http.Request) {
	requestID := r.URL.Query().Get("id")
	if requestID == "" {
		c.renderTemplate(w, &templates.GameData{
			Error: "Missing request ID",
		})
		return
	}
	if user == nil {
		query := url.Values{}
		query.Set("game_login_request_id", requestID)
		http.Redirect(w, r, "/login?"+query.Encode(), http.StatusSeeOther)
		return
	}
	gameLoginRequest, err := c.gameAPIService.GetGameLoginRequestByID(r.Context(), requestID)
	if err != nil {
		c.renderTemplate(w, &templates.GameData{
			Error: err.Error(),
		})
		return
	}
	c.renderTemplate(w, &templates.GameData{GameLoginRequest: gameLoginRequest, User: user})
}

func (c *GameController) PostGameLogin(user *repository.User, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	requestID := r.FormValue("game_login_request_id")
	if requestID == "" {
		http.Error(w, "Missing game login request ID", http.StatusBadRequest)
		return
	}
	if user == nil {
		query := url.Values{}
		query.Set("game_login_request_id", requestID)
		http.Redirect(w, r, "/login?"+query.Encode(), http.StatusSeeOther)
		return
	}
	err := c.gameAPIService.Login(r.Context(), requestID, user)
	if err != nil {
		c.jsonResponse(w, GameErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type GameLoginStateResponse struct {
	ID   string `json:"id"`
	User *struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"user,omitempty"`
}

func (c *GameController) GetGameLoginState(w http.ResponseWriter, r *http.Request) {
	requestID := r.URL.Query().Get("id")
	if requestID == "" {
		c.jsonResponse(w, GameErrorResponse{Message: "Missing request ID"}, http.StatusBadRequest)
		return
	}
	gameLoginRequest, err := c.gameAPIService.GetGameLoginRequestState(r.Context(), requestID)
	if err != nil {
		c.jsonResponse(w, GameErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	if gameLoginRequest == nil {
		c.jsonResponse(w, GameErrorResponse{Message: "Game login request not found"}, http.StatusNotFound)
		return
	}
	response := GameLoginStateResponse{
		ID: gameLoginRequest.ID,
	}
	if gameLoginRequest.LoginedUserID != nil {
		response.User = &struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		}{
			ID:       gameLoginRequest.LoginedUser.ID,
			Username: gameLoginRequest.LoginedUser.Username,
		}
	}
	c.jsonResponse(w, response, http.StatusOK)
}
