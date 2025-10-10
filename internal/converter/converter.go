package converter

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/sborsh1kmusora/auth/internal/model"
	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserDescFromService(user *model.User) *desc.User {
	var updatedAt *timestamp.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        user.ID,
		UserInfo:  ToUserInfoFromService(&user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserInfoFromService(info *model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Name:    info.Name,
		Email:   info.Email,
		IsAdmin: info.IsAdmin,
	}
}

func ToCreateInputFromDesc(info *desc.CreateRequest) *model.CreateInput {
	return &model.CreateInput{
		UserInfo:        ToUserInfoFromDesc(info.UserInfo),
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
	}
}

func ToUserInfoFromDesc(info *desc.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:    info.Name,
		Email:   info.Email,
		IsAdmin: info.IsAdmin,
	}
}

func ToUpdateInputFromDesc(info *desc.UpdateRequest) *model.UpdateInput {
	var name, email *string

	if info.Name != nil {
		name = &info.Name.Value
	}

	if info.Email != nil {
		email = &info.Email.Value
	}

	return &model.UpdateInput{
		ID:    info.Id,
		Name:  name,
		Email: email,
	}
}
