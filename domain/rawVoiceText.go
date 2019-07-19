package domain

import (
	"context"

	"google.golang.org/appengine/datastore"
)

func GetRawVoiceTexts(ctx context.Context, lessonID string) ([]RawVoiceText, error) {
	var voiceTexts []RawVoiceText
	query := datastore.NewQuery("RawVoiceText").Filter("LessonID =", lessonID).Order("FileID")
	if _, err := query.GetAll(ctx, &voiceTexts); err != nil {
		return voiceTexts, err
	}

	return voiceTexts, nil
}

func DeleteRawVoiceTextsByLessonID(ctx context.Context, lessonID string) error {
	var voiceTexts []RawVoiceText
	var keys []*datastore.Key
	query := datastore.NewQuery("RawVoiceText").Filter("LessonID =", lessonID)

	if _, err := query.GetAll(ctx, &voiceTexts); err != nil {
		return err
	}

	if err := datastore.DeleteMulti(ctx, keys); err != nil {
		return err
	}

	return nil
}