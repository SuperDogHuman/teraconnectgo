package handler

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/SuperDogHuman/teraconnectgo/usecase"
	"github.com/labstack/echo/v4"
)

func getAuthoringLesson(c echo.Context) error {
	id := c.Param("id")

	ids := []string{id}
	if !IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lesson, err := usecase.GetAuthoringLesson(c.Request(), id)
	if err != nil {
		lessonErr, ok := err.(usecase.AuthoringLessonErrorCode)
		if ok && lessonErr == usecase.AuthoringLessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lesson)
}

func createAuthoringLesson(c echo.Context) error {
	postedLesson := new(domain.Lesson)

	if err := c.Bind(postedLesson); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	lesson, err := usecase.CreateAuthoringLesson(c.Request(), *postedLesson)
	if err != nil {
		fatalLog(err)
		authoringLessonErr, ok := err.(usecase.AuthoringLessonErrorCode)
		if ok && authoringLessonErr == usecase.InvalidAuthoringLessonParams {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, lesson)
}

func updateAuthoringLesson(c echo.Context) error {
	id := c.Param("id")

	ids := []string{id}
	// TODO add checking of avatarID, graphicIDs
	if !IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		fatalLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lesson, err := usecase.UpdateAuthoringLesson(id, c.Request())
	if err != nil {
		fatalLog(err)

		authoringLessonErr, ok := err.(usecase.AuthoringLessonErrorCode)
		if ok && authoringLessonErr == usecase.AuthoringLessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok && authoringLessonErr == usecase.InvalidAuthoringLessonParams {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lesson)
}