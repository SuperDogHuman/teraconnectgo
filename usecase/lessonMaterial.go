package usecase

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/super-dog-human/teraconnectgo/domain"
)

// CreateLessonMaterialParams
type CreateLessonMaterialParams struct {
	DurationSec       float32                     `json:"durationSec"`
	AvatarID          int64                       `json:"avatarID"`
	AvatarLightColor  string                      `json:"avatarLightColor"`
	BackgroundImageID int64                       `json:"backgroundImageID"`
	BackgroundMusicID int64                       `json:"backgroundMusicID"`
	AvatarMovings     []domain.LessonAvatarMoving `json:"avatarMovings"`
	Graphics          []domain.LessonGraphic      `json:"graphics"`
	Drawings          []domain.LessonDrawing      `json:"drawings"`
}

type LessonMaterialErrorCode uint

const (
	LessonMaterialNotAvailable LessonMaterialErrorCode = 1
	LessonMaterialNotFound     LessonMaterialErrorCode = 2
)

func (e LessonMaterialErrorCode) Error() string {
	switch e {
	case LessonMaterialNotAvailable:
		return "lesson material not available"
	case LessonMaterialNotFound:
		return "lesson material not found"
	default:
		return "unknown lesson error"
	}
}

func GetLessonMaterial(request *http.Request, lessonID int64) (domain.LessonMaterial, error) {
	ctx := request.Context()

	var lessonMaterial domain.LessonMaterial
	if _, err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return lessonMaterial, LessonMaterialNotAvailable
	}

	if err := domain.GetLessonMaterial(ctx, lessonID, &lessonMaterial); err != nil {
		return lessonMaterial, err
	}

	if lessonMaterial.ID == 0 {
		return lessonMaterial, LessonMaterialNotFound
	}

	return lessonMaterial, nil
}

func CreateLessonMaterial(request *http.Request, lessonID int64, params CreateLessonMaterialParams) error {
	ctx := request.Context()

	userID, err := currentUserAccessToLesson(ctx, request, lessonID)
	if err != nil {
		return LessonMaterialNotAvailable
	}

	var lessonMaterial domain.LessonMaterial
	copier.Copy(&lessonMaterial, &params)
	lessonMaterial.UserID = userID

	if err := domain.CreateLessonMaterial(ctx, lessonID, &lessonMaterial); err != nil {
		return err
	}

	return nil
}
