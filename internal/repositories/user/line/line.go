package line

import (
	"context"

	"github.com/Rocksus/pogo/internal/repositories/user"
	"github.com/line/line-bot-sdk-go/linebot"
)

type userRepository struct {
	client *linebot.Client
}

func NewUserRepository(client *linebot.Client) user.Repository {
	return &userRepository{client: client}
}

func (u *userRepository) GetByID(ctx context.Context, userID string) (profile user.User, err error) {
	resp, err := u.client.GetProfile(userID).WithContext(ctx).Do()
	if err != nil {
		return
	}

	profile = user.User{
		ID:         resp.UserID,
		Name:       resp.DisplayName,
		PictureURL: resp.PictureURL,
	}
	return
}
