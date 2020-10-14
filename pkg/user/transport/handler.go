package transport

import (
	"net/http"

	userent "github.com/dev681999/tpo-platform/pkg/ent/user"
	"github.com/dev681999/tpo-platform/pkg/user"
	"github.com/dev681999/tpo-platform/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Handler is user transport handler
type Handler struct {
	userSvc user.Service
	logger  zerolog.Logger
}

//NewHandler returns a new Hadler
func NewHandler(userSvc user.Service, logger zerolog.Logger, g *echo.Group) Handler {
	h := Handler{
		userSvc: userSvc,
		logger:  logger,
	}

	h.initRoutes(g)

	return h
}

func (h *Handler) initRoutes(g *echo.Group) {
	g.POST("/register", h.Register)
	g.POST("/register/", h.Register)

	g.GET("", h.GetAll)
	g.GET("/", h.GetAll)

	g.GET("/:id", h.Get)
	g.GET("/:id/", h.Get)
}

// Get gets a single user
func (h *Handler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	uuid, err := utils.ParseUUID(id)
	if err != nil {
		h.logger.Err(err).Msg("")
		return utils.NewEchoErrorResponse(http.StatusBadRequest, err)
	}

	u, err := h.userSvc.Get(ctx, uuid)
	if err != nil {
		h.logger.Err(err).Msg("")
		return utils.NewEchoErrorResponse(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user": u,
	})
}

// GetAll get a all users
func (h *Handler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := h.userSvc.GetAll(ctx)
	if err != nil {
		h.logger.Err(err).Msg("")
		return utils.NewEchoErrorResponse(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"users": users,
	})
}

type registrationDetails struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Register registers a user
func (h *Handler) Register(c echo.Context) error {
	ctx := c.Request().Context()
	rd := registrationDetails{}
	err := c.Bind(&rd)
	if err != nil {
		h.logger.Err(err).Msg("")
		return utils.NewEchoErrorResponse(http.StatusBadRequest, err)
	}

	u, err := h.userSvc.Create(ctx, rd.Email, rd.Password, rd.Name, userent.RoleStudent)
	if err != nil {
		h.logger.Err(err).Msg("")
		return utils.NewEchoErrorResponse(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id": u.ID,
	})
}
