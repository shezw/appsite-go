// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package account

import (
	"errors"

	"gorm.io/gorm"

	"appsite-go/internal/core/model"
	"appsite-go/internal/services/user/dto"
	"appsite-go/internal/services/user/entity"
)

// Add creates a new user with strict field control (UserAccount.php addFields)
func (s *AuthService) Add(input dto.UserCreateReq) (*entity.User, error) {
	// 1. Check uniqueness (Username, Email, Mobile)
	// Build flexible query
	query := s.db.Model(&entity.User{})
	conds := []string{}
	args := []interface{}{}

	if input.Username != "" {
		conds = append(conds, "username = ?")
		args = append(args, input.Username)
	}
	if input.Email != "" {
		conds = append(conds, "email = ?")
		args = append(args, input.Email)
	}
	if input.Mobile != "" {
		conds = append(conds, "mobile = ?")
		args = append(args, input.Mobile)
	} else if input.UID != "" { 
		// Check UID if provided manually (though rare)
		conds = append(conds, "id = ?")
		args = append(args, input.UID)
	}

	if len(conds) > 0 {
		var count int64
		// "username = ? OR email = ? OR ..."
		queryStr := ""
		for i, c := range conds {
			if i > 0 {
				queryStr += " OR "
			}
			queryStr += c
		}
		if err := query.Where(queryStr, args...).Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, ErrUserExists
		}
	}

	// 2. Hash Password (beforeAdd hook)
	// If password is provided, hash it. If not, maybe allow empty? (e.g. 3rd party login)
	var hashedPwd string
	if input.Password != "" {
		var err error
		hashedPwd, err = s.pwd.Hash(input.Password)
		if err != nil {
			return nil, err
		}
	}

	// 3. Map Fields
	// Handle Pointers for optional fields to avoid empty string unique violation
	var email *string
	if input.Email != "" {
		email = &input.Email
	}
	var mobile *string
	if input.Mobile != "" {
		mobile = &input.Mobile
	}

	user := &entity.User{
		// Base
		Tenant: model.Tenant{SaasID: input.SaasID},
		
		// Auth
		Username: input.Username,
		Password: hashedPwd,
		Email:    email,
		Mobile:   mobile,

		// Profile
		Nickname:    input.Nickname,
		Avatar:      input.Avatar,
		Status:      input.Status,
		
		// Extended (UserAccount.php)
		Cover:       input.Cover,
		Description: input.Description,
		Introduce:   input.Introduce,
		Birthday:    input.Birthday,
		Gender:      input.Gender,
		AreaID:      input.AreaID,
		GroupID:     input.GroupID,
	}

	// Manual UID?
	if input.UID != "" {
		user.ID = input.UID
	}
	
	// Defaults if missing (Go zero values are often valid, but check constraints)
	if user.Status == "" {
		user.Status = "enabled"
	}
	if user.GroupID == "" {
		user.GroupID = "100"
	}
	if user.AreaID == "" {
		user.AreaID = "1"
	}
	if user.Gender == "" {
		user.Gender = "private"
	}

	// 4. Create
	if res := s.repo.Add(user); !res.Success {
		if res.Error != nil {
			return nil, res.Error
		}
		return nil, errors.New(res.Message)
	}

	return user, nil
}

// Update updates a user with strict field control (UserAccount.php updateFields)
func (s *AuthService) Update(uid string, input dto.UserUpdateReq) error {
	var user entity.User
	if err := s.db.First(&user, "id = ?", uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	updates := make(map[string]interface{})

	// helper to set if not nil
	if input.Password != nil {
		hash, err := s.pwd.Hash(*input.Password)
		if err != nil {
			return err
		}
		updates["password"] = hash
	}
	if input.Email != nil {
		if *input.Email == "" {
			updates["email"] = nil
		} else {
			updates["email"] = *input.Email
		}
	}
	if input.Mobile != nil {
		if *input.Mobile == "" {
			updates["mobile"] = nil
		} else {
			updates["mobile"] = *input.Mobile
		}
	}
	if input.Nickname != nil {
		updates["nickname"] = *input.Nickname
	}
	if input.Avatar != nil {
		updates["avatar"] = *input.Avatar
	}
	if input.Cover != nil {
		updates["cover"] = *input.Cover
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Introduce != nil {
		updates["introduce"] = *input.Introduce
	}
	if input.Birthday != nil {
		updates["birthday"] = *input.Birthday
	}
	if input.Gender != nil {
		updates["gender"] = *input.Gender
	}
	if input.GroupID != nil {
		updates["group_id"] = *input.GroupID
	}
	if input.AreaID != nil {
		updates["area_id"] = *input.AreaID
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}

	if len(updates) == 0 {
		return nil
	}

	return s.db.Model(&user).Updates(updates).Error
}

// GetDetail retrieves full user details
func (s *AuthService) GetDetail(uid string) (*dto.UserDetailResp, error) {
	var user entity.User
	if err := s.db.First(&user, "id = ?", uid).Error; err != nil {
		return nil, err
	}
	
	return &dto.UserDetailResp{
		UID:         user.ID,
		SaasID:      user.SaasID,
		Username:    user.Username,
		Email:       valOrEmpty(user.Email),
		Mobile:      valOrEmpty(user.Mobile),
		Password:    user.Password, // Caution: Hashed
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Cover:       user.Cover,
		Description: user.Description,
		Introduce:   user.Introduce,
		Birthday:    user.Birthday,
		Gender:      user.Gender,
		GroupID:     user.GroupID,
		AreaID:      user.AreaID,
		Status:      user.Status,
		CreateTime:  user.CreatedAt,
		LastTime:    user.UpdatedAt,
	}, nil
}

// valOrEmpty helper
func valOrEmpty(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// ListUsers retrieves users based on filter
func (s *AuthService) ListUsers(filter dto.UserFilterReq, page, pageSize int) ([]dto.UserListResp, int64, error) {
	var users []entity.User
	var count int64

	query := s.db.Model(&entity.User{})

	if filter.UID != "" {
		query = query.Where("id = ?", filter.UID)
	}
	if filter.Username != "" {
		query = query.Where("username LIKE ?", "%"+filter.Username+"%")
	}
	if filter.Email != "" {
		query = query.Where("email = ?", filter.Email)
	}
	if filter.Mobile != "" {
		query = query.Where("mobile = ?", filter.Mobile)
	}
	if filter.Nickname != "" {
		query = query.Where("nickname LIKE ?", "%"+filter.Nickname+"%")
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.GroupID != "" {
		query = query.Where("group_id = ?", filter.GroupID)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	res := make([]dto.UserListResp, len(users))
	for i, u := range users {
		res[i] = dto.UserListResp{
			UID:         u.ID,
			Username:    u.Username,
			Email:       valOrEmpty(u.Email),
			Mobile:      valOrEmpty(u.Mobile),
			SaasID:      u.SaasID,
			Nickname:    u.Nickname,
			Avatar:      u.Avatar,
			Cover:       u.Cover,
			Description: u.Description,
			GroupID:     u.GroupID,
			Gender:      u.Gender,
			AreaID:      u.AreaID,
			Status:      u.Status,
			CreateTime:  u.CreatedAt,
			LastTime:    u.UpdatedAt,
		}
	}

	return res, count, nil
}
