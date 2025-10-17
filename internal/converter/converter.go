package converter

import (
	"github.com/sborsh1kmusora/auth/internal/model"
	desc "github.com/sborsh1kmusora/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserInfoFromDesc(desc *desc.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Username: desc.GetUsername(),
		Password: desc.GetPassword(),
		Role:     desc.GetRole().String(),
	}
}

func ToDescFromUser(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        user.ID,
		UserInfo:  ToDescFromUserInfo(&user.UserInfo),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToDescFromUserInfo(userInfo *model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Username: userInfo.Username,
		Password: userInfo.Password,
		Role:     toDescRole(userInfo.Role),
	}
}

func ToUpdateUserFromDesc(desc *desc.UpdateRequest) *model.UpdateUser {
	return &model.UpdateUser{
		ID:       desc.GetId(),
		Username: desc.Username,
		Password: desc.Password,
	}
}

func toDescRole(role string) desc.Role {
	if role == "admin" {
		return desc.Role_ADMIN
	}

	return desc.Role_USER
}
