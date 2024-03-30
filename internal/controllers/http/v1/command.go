package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/timohahaa/executor/internal/service"
)

const (
	DEFAULT_LIMIT uint64 = 50
)

type commandRoutes struct {
	commandService service.CommandService
}

func newCommandRoutes(g *echo.Group, cs service.CommandService) {
	r := &commandRoutes{
		commandService: cs,
	}

	g.POST("/command", r.CreateCommand)
	g.GET("/command/:commandId", r.GetCommandById)
	g.POST("/command/:commandId/run", r.RunCommand)
	g.GET("/commands", r.ListCommands)
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
	Id         uint64 `json:"command_id"`
	Text       string `json:"command_text"`
	LastOutput string `json:"last_output"`
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

	output := getCommandOutput{Id: command.Id, Text: command.Text, LastOutput: command.LastOutput}

	return c.JSON(http.StatusOK, output)
}

type listCommandsQueryParams struct {
	Limit  uint64 `query:"limit"`
	Offset uint64 `query:"offset"`
}

type listCommandsOutput struct {
	Data []getCommandOutput `json:"data"`
	Next struct {
		Limit  uint64 `json:"limit"`
		Offset uint64 `json:"offset"`
	} `json:"next"`
}

// GET /api/v1/commands
func (r *commandRoutes) ListCommands(c echo.Context) error {
	var queryParams listCommandsQueryParams
	if err := c.Bind(&queryParams); err != nil {
		queryParams.Limit = DEFAULT_LIMIT
		queryParams.Offset = 0
	}

	commands, err := r.commandService.ListCommands(c.Request().Context(), queryParams.Limit, queryParams.Offset)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}

	output := listCommandsOutput{}
	output.Next.Limit = queryParams.Limit
	output.Next.Offset = queryParams.Offset + queryParams.Limit + 1
	for _, command := range commands {
		output.Data = append(output.Data, getCommandOutput{Id: command.Id, Text: command.Text, LastOutput: command.LastOutput})
	}

	return c.JSON(http.StatusOK, output)
}

// POST /api/v1/command/{commandId}/run
func (r *commandRoutes) RunCommand(c echo.Context) error {
	return nil
}
