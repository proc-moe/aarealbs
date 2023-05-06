package auth

import (
	"errors"

	"github.com/proc-moe/aarealbs/server/model"
	"github.com/proc-moe/aarealbs/server/utils/klog"
)

func GetUserStatus(token string) (int, int, error) {
	x := model.Token{}
	r := model.DB.Where("token = ?", token).First(&x)
	if r.RowsAffected == 0 {
		klog.E("no such token")
		return 0, 0, errors.New("no such token")
	}
	id := x.UserInfoID

	userInfo := model.UserInfo{}
	r = model.DB.Where("user_id = ?", id).First(&userInfo)
	if r.RowsAffected == 0 {
		klog.E("no such user_id")
		return 0, 0, errors.New("no such user_id")
	}
	klog.I("[GetUserStatus] userID=%v,userStatus=%v", userInfo.UserId, userInfo.Status)
	return int(userInfo.UserId), int(userInfo.Status), nil
}

func UserIsAdmin(token string) error {
	_, status, err := GetUserStatus(token)
	if err != nil {
		return err
	}
	if status != 2 {
		return errors.New("unauthorized, user is not admin")
	}
	return nil
}

func UserIsSelfAndUnbanned(token string, uid int) error {
	id, status, err := GetUserStatus(token)
	if err != nil {
		return err
	}
	if status == 3 {
		return errors.New("unauthorized, user banned")
	}
	if id != uid {
		return errors.New("unauthorzied, user not self")
	}
	return nil
}
