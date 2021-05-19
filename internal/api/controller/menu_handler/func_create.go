package menu_handler

import (
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/api/code"
	"github.com/xinliangnote/go-gin-api/internal/api/service/menu_service"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/pkg/errno"

	"github.com/spf13/cast"
)

type createRequest struct {
	Id    string `form:"id"`    // ID
	Pid   int32  `form:"pid"`   // Parent class ID
	Name  string `form:"name"`  // Menu name
	Link  string `form:"link"`  // link address
	Icon  string `form:"icon"`  // icon
	Level int32  `form:"level"` // Menu type 1: Level 1 menu 2: Level 2 menu
}

type createResponse struct {
	Id int32 `json:"id"` // Primary Key ID
}

// Create create/edit menu
// @Summary create/edit menu
// @Description create/edit menu
// @Tags API.menu
// @Accept multipart/form-data
// @Produce json
// @Param Request body createRequest true "Request information"
// @Success 200 {object} createResponse
// @Failure 400 {object} code.Failure
// @Router /api/menu [post]
func (h *handler) Create() core.HandlerFunc {
	return func(c core.Context) {
		req := new(createRequest)
		res := new(createResponse)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(errno.NewError(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithErr(err),
			)
			return
		}

		if req.Id != "" { // Edit function
			ids, err := h.hashids.HashidsDecode(req.Id)
			if err != nil {
				c.AbortWithError(errno.NewError(
					http.StatusBadRequest,
					code.HashIdsDecodeError,
					code.Text(code.HashIdsDecodeError)).WithErr(err),
				)
				return
			}

			id := int32(ids[0])

			updateData := new(menu_service.UpdateMenuData)
			updateData.Name = req.Name
			updateData.Icon = req.Icon
			updateData.Link = req.Link

			err = h.menuService.Modify(c, id, updateData)
			if err != nil {
				c.AbortWithError(errno.NewError(
					http.StatusBadRequest,
					code.MenuUpdateError,
					code.Text(code.MenuUpdateError)).WithErr(err),
				)
				return
			}

			res.Id = id
			c.Payload(res)

		} else { // New features

			pid := req.Level
			level := 2

			if req.Level == -1 {
				pid = 0
				level = 1
			}

			createData := new(menu_service.CreateMenuData)
			createData.Pid = pid
			createData.Name = req.Name
			createData.Icon = req.Icon
			createData.Link = req.Link
			createData.Level = cast.ToInt32(level)

			id, err := h.menuService.Create(c, createData)
			if err != nil {
				c.AbortWithError(errno.NewError(
					http.StatusBadRequest,
					code.MenuCreateError,
					code.Text(code.MenuCreateError)).WithErr(err),
				)
				return
			}

			res.Id = id
			c.Payload(res)
		}
	}
}
