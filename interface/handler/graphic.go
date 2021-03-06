package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func getGraphic(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	graphic, err := usecase.GetGraphicByID(c.Request(), id)
	if err != nil {
		if ok := errors.Is(err, domain.GraphicNotFound); ok {
			warnLog(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, graphic)
}

func getGraphics(c echo.Context) error {
	lessonID, err := strconv.ParseInt(c.QueryParam("lesson_id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	graphics, err := usecase.GetGraphicsByLessonID(c.Request(), lessonID)

	if err != nil {
		fatalLog(err)
		graphicErr, ok := err.(domain.GraphicErrorCode)
		if ok && graphicErr == domain.GraphicNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(graphics) == 0 {
		errMessage := "graphics not found"
		warnLog(errMessage)
		return c.JSON(http.StatusNotFound, errMessage)
	}

	return c.JSON(http.StatusOK, graphics)
}

func postGraphics(c echo.Context) error {
	objectRequest := new(infrastructure.StorageObjectRequest)
	if err := c.Bind(objectRequest); err != nil {
		fatalLog(err)
		authErr, ok := err.(domain.AuthErrorCode)
		if ok && authErr == domain.UserNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	signedURLs, err := usecase.CreateGraphicsAndBlankFiles(c.Request(), *objectRequest)
	if err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, signedURLs)
}

func deleteGraphic(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	err = usecase.DeleteGraphic(c.Request(), id)
	if err != nil {
		fatalLog(err)

		if ok := errors.Is(err, domain.GraphicNotFound); ok {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "the graphic has deleted.")
}
