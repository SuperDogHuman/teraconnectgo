package domain

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// TODO move to infrastructure for at development settings.
const avatarThumbnailURL = "https://storage.googleapis.com/teraconn_thumbnail/avatar/{id}.png"

// GetAvailableAvatars for fetch avatar object from Cloud Datastore
func GetAvailableAvatars(request *http.Request) ([]Avatar, error) {
	ctx := appengine.NewContext(request)

	var avatars []Avatar

	currentUser, err := GetCurrentUser(request)
	if err != nil {
		return nil, err
	}

	usersAvatars, err := getCurrentUsersAvatars(ctx, currentUser.ID)
	if err != nil {
		return nil, err
	}
	avatars = append(avatars, usersAvatars...)

	publicAvatars, err := getPublicAvatars(ctx)
	if err != nil {
		return nil, err
	}
	avatars = append(avatars, publicAvatars...)

	return avatars, nil
}

func getCurrentUsersAvatars(ctx context.Context, userID string) ([]Avatar, error){
	var avatars []Avatar

	query := datastore.NewQuery("Avatar").Filter("UserId =", userID)
	keys, err := query.GetAll(ctx, &avatars)
	if err != nil {
		return nil, err
	}

	storeAvatarThumbnailUrl(&avatars, keys)

	return avatars, nil
}

func getPublicAvatars(ctx context.Context) ([]Avatar, error){
	var avatars []Avatar

	query := datastore.NewQuery("Avatar").Filter("IsPublic =", true)
	keys, err := query.GetAll(ctx, &avatars)
	if err != nil {
		return nil, err
	}

	storeAvatarThumbnailUrl(&avatars, keys)

	return avatars, nil
}

func storeAvatarThumbnailUrl(avatars *[]Avatar, keys []*datastore.Key) {
	for i, key := range keys {
		id := key.StringID()
		(*avatars)[i].ID = id
		(*avatars)[i].ThumbnailURL = strings.Replace(avatarThumbnailURL, "{id}", id, 1)
	}
}