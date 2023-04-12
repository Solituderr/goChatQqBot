package model

//操作消息等

func SaveMsg(msg UserInfo) error {
	err := DB.Create(msg).Error
	return err
}

func FindMsg(userid string) ([]UserInfo, error) {
	var users []UserInfo
	err := DB.Where("user_id = ?", userid).Find(&users).Error
	return users, err
}

func DeleteMsg(userid string, checkSys bool) error {
	if checkSys {
		err := DB.Where("user_id = ?", userid).Delete(&UserInfo{}).Error
		return err
	} else {
		err := DB.Where("user_id = ? and (role='user' or role='assistant')", userid).Delete(&UserInfo{}).Error
		return err
	}

}

// 更改当前设定
func ChangeSys(userid string, role string, define string) error {
	err := DB.Model(&UserInfo{}).Where("user_id = ? and role = ?", userid, role).Update("content", define).Error
	//err := DB.Raw("update user_infos set content=? where user_id=?", define, userid).Error
	return err
}

// 查看用户是否与机器人沟通过
func CheckUser(userid string) error {
	var user UserInfo
	err := DB.Where("user_id = ?", userid).First(&user).Error
	return err
}

// 操作权限等
func CheckAccess(userid string) error {
	var acc UserAccess
	err := DB.Model(&UserAccess{}).Where("user_id = ?", userid).First(&acc).Error
	return err
}

func AddAccess(userid string) error {
	var acc UserAccess
	acc.Acc = true
	acc.UserId = userid
	var err error
	if CheckAccess(userid) != nil {
		err = DB.Model(&UserAccess{}).Create(acc).Error
	} else {
		err = nil
	}
	return err
}

func DeleteAccess(userid string) error {
	err := DB.Model(&UserAccess{}).Where("user_id = ?", userid).Delete(&UserAccess{}).Error
	return err
}
