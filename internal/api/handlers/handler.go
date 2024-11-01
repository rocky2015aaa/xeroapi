package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/shmoulana/xeroapi/internal/api/models"
	"golang.org/x/oauth2"
)

const (
	EnvXeroConnectionUrl             = "XERO_CONNECTION_URL"
	EnvXeroItemUrl                   = "XERO_ITEM_URL"
	EnvXeroContractUrl               = "XERO_CONTRACT_URL"
	EnvXeroTrackingCategoryUrl       = "XERO_TRACKING_CATEGORY_URL"
	EnvXeroTrackingCategoryOptionUrl = "XERO_TRACKING_CATEGORY_OPTION_URL"
	EnvXeroAccountUrl                = "XERO_ACCOUNT_URL"
	PUTMethod                        = "PUT"

	responseSuccess     = "success"
	responseData        = "data"
	responseError       = "error"
	responseDescription = "description"
)

type Handler struct {
	httpClient   *http.Client
	oAuth2Config *oauth2.Config
	oAuthToken   *oauth2.Token
	connections  []*models.Connection
}

func NewHandler(oAuth2Config *oauth2.Config) *Handler {
	return &Handler{
		oAuth2Config: oAuth2Config,
	}
}

func getResponse(success bool, data interface{}, err, description string) gin.H {
	return gin.H{
		responseSuccess:     success,
		responseData:        data,
		responseError:       err,
		responseDescription: description,
	}
}

func (h *Handler) IndexHandler(ctx *gin.Context) {
	var loginInformationJSON []byte

	if h.oAuthToken.Valid() {
		connectionsResponse, err := h.httpClient.Get(os.Getenv(EnvXeroConnectionUrl))

		if err != nil && connectionsResponse.StatusCode != 200 {
			log.Fatalln(err)
		}
		respBody, err := io.ReadAll(connectionsResponse.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var connections []*models.Connection
		err = json.Unmarshal(respBody, &connections)
		if err != nil {
			log.Fatalln(err)
		}
		h.connections = connections
		loginInformation := models.LoginInformation{
			Connections: connections,
			TokenDetails: &models.TokenDetails{
				AccessToken:  h.oAuthToken.AccessToken,
				RefreshToken: h.oAuthToken.RefreshToken,
				Expiry:       h.oAuthToken.Expiry.String(),
				TokenType:    h.oAuthToken.TokenType,
				Scope:        h.oAuthToken.Extra("scope").(string),
			},
		}
		loginInformationJSON, _ = json.MarshalIndent(loginInformation, "", "    ")
	}
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"IsTokenValid": h.oAuthToken.Valid(),
		"OAuthToken":   string(loginInformationJSON),
	})
}

func (h *Handler) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, getResponse(true, nil, "", "pong"))
}

func (h *Handler) LoginHandler(c *gin.Context) {
	// Generate the OAuth2 URL for redirection
	redirectURL := h.oAuth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)

	// Use Gin's context to set the Location header and issue a redirect
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *Handler) HandleOAuthCallback(ctx *gin.Context) {
	baseToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("Basic %s:%s", h.oAuth2Config.ClientID, h.oAuth2Config.ClientSecret)))
	// Exchange the authorization code for a token
	tok, err := h.oAuth2Config.Exchange(
		ctx,
		ctx.Query("code"),
		oauth2.SetAuthURLParam("authorization", baseToken),
	)
	if err != nil {
		log.Println(err)
	}

	// Update the server's OAuth token
	h.oAuthToken = tok
	h.httpClient = h.oAuth2Config.Client(ctx, tok)
	// Redirect the user back to the home page
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *Handler) HandleTokenRefreshRequest(ctx *gin.Context) {
	if !h.oAuthToken.Valid() {
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// Write global styles if needed
	src := h.oAuth2Config.TokenSource(ctx, &oauth2.Token{RefreshToken: h.oAuthToken.RefreshToken})
	newToken, err := src.Token()
	if err != nil {
		log.Println("An error occurred while get refresh token.")
		log.Fatalln(err)
	}
	// Also update the Server struct properties
	h.oAuthToken = newToken

	if !h.oAuthToken.Valid() {
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *Handler) CreateItems(ctx *gin.Context) {
	req := models.Items{}
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest,
			getResponse(false, nil, err.Error(), "Binding data has failed"))
		return
	}
	payload, err := json.Marshal(req)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "Failed to marshal payload"))
		return
	}

	responseData, err := callXeroApi(PUTMethod, os.Getenv(EnvXeroItemUrl),
		h.oAuthToken.AccessToken, h.connections[0].TenantID, payload)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "Failed to call xero api"))
		return
	}

	// Respond with success
	ctx.JSON(http.StatusOK, getResponse(true, responseData, "", "Creating the item has succeeded"))
}

func (h *Handler) CreateContacts(ctx *gin.Context) {
	req := models.Contacts{}
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest,
			getResponse(false, nil, err.Error(), "Binding data has failed"))
		return
	}
	payload, err := json.Marshal(req)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "Failed to marshal payload"))
		return
	}

	responseData, err := callXeroApi(PUTMethod, os.Getenv(EnvXeroContractUrl),
		h.oAuthToken.AccessToken, h.connections[0].TenantID, payload)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "Failed to call xero api"))
		return
	}

	// Respond with success
	ctx.JSON(http.StatusOK, getResponse(true, responseData, "", "Creating the contacts has succeeded"))
}

func (h *Handler) CreateTrackingCategory(ctx *gin.Context) {
	req := struct {
		Name string `json:"name"`
	}{}
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest,
			getResponse(false, nil, err.Error(), "Binding data has failed"))
		return
	}
	payload, err := json.Marshal(req)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "Failed to marshal payload"))
		return
	}

	responseData, err := callXeroApi(PUTMethod, os.Getenv(EnvXeroTrackingCategoryUrl),
		h.oAuthToken.AccessToken, h.connections[0].TenantID, payload)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "Failed to call xero api"))
		return
	}

	// Respond with success
	ctx.JSON(http.StatusOK, getResponse(true, responseData, "", "Creating the tracking category has succeeded"))
}

func (h *Handler) CreateTrackingOption(ctx *gin.Context) {
	trackingCategoryID := ctx.Param("tracking_category_ID")
	req := struct {
		Name string `json:"name"`
	}{}
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest,
			getResponse(false, nil, err.Error(), "Binding data has failed"))
		return
	}
	payload, err := json.Marshal(req)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "Failed to marshal payload"))
		return
	}

	responseData, err := callXeroApi(PUTMethod, fmt.Sprintf(os.Getenv(EnvXeroTrackingCategoryOptionUrl), trackingCategoryID),
		h.oAuthToken.AccessToken, h.connections[0].TenantID, payload)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "Failed to call xero api"))
		return
	}

	// Respond with success
	ctx.JSON(http.StatusOK, getResponse(true, responseData, "", "Creating the option has succeeded"))
}

func (h *Handler) CreateAccount(ctx *gin.Context) {
	req := struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
	}{}
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest,
			getResponse(false, nil, err.Error(), "Binding data has failed"))
		return
	}
	payload, err := json.Marshal(req)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "Failed to marshal payload"))
		return
	}

	responseData, err := callXeroApi(PUTMethod, os.Getenv(EnvXeroAccountUrl),
		h.oAuthToken.AccessToken, h.connections[0].TenantID, payload)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "Failed to call xero api"))
		return
	}

	// Respond with success
	ctx.JSON(http.StatusOK, getResponse(true, responseData, "", "Creating the option has succeeded"))
}
