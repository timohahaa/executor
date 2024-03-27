package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/timohahaa/executor/internal/service"
)

type commandRoutes struct {
	commandService service.CommandService
}

func newCommandRoutes(g *echo.Group, cs service.CommandService) {
	r := &commandRoutes{
		commandService: cs,
	}

	g.POST("/command", r.CreateCommand)
}

type createCommandInput struct {
	CommandText string `json:"command_text"`
}

type createCommandOutput struct {
	Id uint64 `json:"command_id"`
}

// POST /api/v1/command
func (r *commandRoutes) CreateCommand(c echo.Context) error {
	var input createCommandInput
	if err := c.Bind(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return err
	}

	command, err := r.commandService.CreateCommand(c.Request().Context(), input.CommandText)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}

	output := createCommandOutput{Id: command.Id}

	return c.JSON(http.StatusCreated, output)
}

type getCommandOutput struct {
	Id   uint64 `json:"command_id"`
	Text string `json:"command_text"`
}

// GET /api/v1/command/{commandId}
func (r *commandRoutes) GetCommandById(c echo.Context) error {
	cId := c.Param("commandId")
	commandId, err := strconv.ParseUint(cId, 10, 64)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, ErrInvalidPathParameter.Error())
		return err
	}

	command, err := r.commandService.GetCommandById(c.Request().Context(), commandId)
	if errors.Is(err, service.ErrCommandNotFound) {
		newErrorMessage(c, http.StatusNotFound, ErrCommandNotFound.Error())
		return ErrCommandNotFound
	} else if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}

	output := getCommandOutput{Id: command.Id, Text: command.Text}

	return c.JSON(http.StatusOK, output)
}
