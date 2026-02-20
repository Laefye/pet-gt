package controllers

import (
	"encoding/json"
	"gt/internal/middleware"
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

func (c *GameController) jsonResponse(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (c *GameController) renderTemplate(w http.ResponseWriter, data *templates.GameData) {
	err := templates.GameTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type gameLoginRequestResponse struct {
	ID  string `json:"id"`
	URL string `json:"url,omitempty"`
}

type gameErrorResponse struct {
	Message string `json:"message"`
}

type gameLoginStateResponse struct {
	ID   string              `json:"id"`
	User *gameLoginStateUser `json:"user,omitempty"`
}

type gameLoginStateUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (c *GameController) CreateGameLoginRequest(w http.ResponseWriter, r *http.Request) {
	req, err := c.gameAPIService.CreateGameLoginRequest(r.Context())
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
	query.Set("id", req.ID)
	u.RawQuery = query.Encode()
	c.jsonResponse(w, gameLoginRequestResponse{ID: req.ID, URL: u.String()}, http.StatusCreated)
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
		query.Set("game_login_request_id", requestID)
		http.Redirect(w, r, "/login?"+query.Encode(), http.StatusSeeOther)
		return
	}
	gameLoginRequest, err := c.gameAPIService.GetGameLoginRequestByID(r.Context(), requestID)
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
	requestID := r.FormValue("game_login_request_id")
	if requestID == "" {
		http.Error(w, "Missing game login request ID", http.StatusBadRequest)
		return
	}
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		query := url.Values{}
		query.Set("game_login_request_id", requestID)
		http.Redirect(w, r, "/login?"+query.Encode(), http.StatusSeeOther)
		return
	}
	err := c.gameAPIService.Login(r.Context(), requestID, user)
	if err != nil {
		c.jsonResponse(w, gameErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *GameController) GetGameLoginState(w http.ResponseWriter, r *http.Request) {
	requestID := r.URL.Query().Get("id")
	if requestID == "" {
		c.jsonResponse(w, gameErrorResponse{Message: "Missing request ID"}, http.StatusBadRequest)
		return
	}
	req, err := c.gameAPIService.GetGameLoginRequestState(r.Context(), requestID)
	if err != nil {
		c.jsonResponse(w, gameErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	if req == nil {
		c.jsonResponse(w, gameErrorResponse{Message: "Game login request not found"}, http.StatusNotFound)
		return
	}
	response := gameLoginStateResponse{ID: req.ID}
	if req.AuthorizedUserID != nil {
		response.User = &gameLoginStateUser{
			ID:       req.AuthorizedUser.ID,
			Username: req.AuthorizedUser.Username,
		}
	}
	c.jsonResponse(w, response, http.StatusOK)
}
