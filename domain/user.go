package domain

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type UserErrorCode uint

const (
	AlreadyProviderIDExists UserErrorCode = 1
)

func (e UserErrorCode) Error() string {
	switch e {
	case AlreadyProviderIDExists:
		return "provider id is already existed"
	default:
		return "unknown error"
	}
}

// GetCurrentUser returns user from valid token.
func GetCurrentUser(request *http.Request) (User, error) {
	user := new(User) // for return blank user when error

	providerID, err := ProviderID(request)
	if err != nil {
		return *user, err
	}

	var users []User
	ctx := request.Context()
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *user, FailedDatastoreInitialize
	}

	query := datastore.NewQuery("User").Filter("ProviderID =", providerID).Limit(1)
	keys, err := client.GetAll(ctx, query, &users)
	if err != nil {
		return *user, FailedGettingUser
	}

	if len(users) == 0 {
		return *user, UserNotFound
	}

	user = &users[0]
	user.ID = strconv.FormatInt(keys[0].ID, 10)
	return users[0], nil
}

// GetUserByID is return user has ID.
func GetUserByID(ctx context.Context, id int64) (User, error) {
	user := new(User)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *user, err
	}

	key := datastore.IDKey("User", id, nil)
	if err := client.Get(ctx, key, user); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return *user, UserNotFound
		}
		return *user, err
	}

	return *user, nil
}

// ReserveUserProviderIDInTransaction creates user's ProviderID for exclusion control.
func ReserveUserProviderIDInTransaction(tx *datastore.Transaction, providerID string) error {
	key := datastore.NameKey("UserProviderID", providerID, nil)
	userProviderID := new(UserProviderID)

	err := tx.Get(key, userProviderID)
	if err == nil {
		return AlreadyProviderIDExists
	}
	if err != datastore.ErrNoSuchEntity {
		return err
	}

	// Put only when ErrNoSuchEntity
	_, err = tx.Put(key, userProviderID)
	return err
}

// CreateUserInTransaction creates new user.
func CreateUserInTransaction(tx *datastore.Transaction, user *User) (*datastore.PendingKey, error) {
	key := datastore.IncompleteKey("User", nil)

	user.Created = time.Now()

	pendingKey, err := tx.Put(key, user)
	if err != nil {
		return nil, err
	}

	return pendingKey, nil
}

func UpdateUser(ctx context.Context, user *User) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	userID, err := strconv.ParseInt(user.ID, 10, 64)
	if err != nil {
		return err
	}
	key := datastore.IDKey("User", userID, nil)

	user.Updated = time.Now()

	if _, err := client.Put(ctx, key, user); err != nil {
		return err
	}

	return nil
}

func DeleteUser(ctx context.Context, id string) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	key := datastore.NameKey("User", id, nil)
	if err := client.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}
