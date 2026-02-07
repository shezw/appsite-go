// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

// UserCreateReq defines fields allowed when creating a user (UserAccount).
// Corresponds to PHP 'addFields'.
type UserCreateReq struct {
	UID         string `json:"uid"` // Optional, usually generated
	SaasID      string `json:"saasid"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
	Introduce   string `json:"introduce"`
	Birthday    int64  `json:"birthday"`
	Gender      string `json:"gender"`
	GroupID     string `json:"groupid"`
	AreaID      string `json:"areaid"`
	Status      string `json:"status"`
}

// UserUpdateReq defines fields allowed when updating a user.
// Corresponds to PHP 'updateFields'.
type UserUpdateReq struct {
	Password    *string `json:"password,omitempty"`
	Email       *string `json:"email,omitempty"`
	Mobile      *string `json:"mobile,omitempty"`
	Nickname    *string `json:"nickname,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
	Cover       *string `json:"cover,omitempty"`
	Description *string `json:"description,omitempty"`
	Introduce   *string `json:"introduce,omitempty"`
	Birthday    *int64  `json:"birthday,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	GroupID     *string `json:"groupid,omitempty"`
	AreaID      *string `json:"areaid,omitempty"`
	Status      *string `json:"status,omitempty"`
}

// UserDetailResp defines fields for internal full detail view.
// Corresponds to PHP 'detailFields'.
type UserDetailResp struct {
	UID         string `json:"uid"`
	SaasID      string `json:"saasid"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	Password    string `json:"password"` // Encrypted?
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
	Introduce   string `json:"introduce"`
	Birthday    int64  `json:"birthday"`
	Gender      string `json:"gender"`
	GroupID     string `json:"groupid"`
	AreaID      string `json:"areaid"`
	Status      string `json:"status"`
	CreateTime  int64  `json:"createtime"`
	LastTime    int64  `json:"lasttime"`
}

// UserPublicDetailResp defines fields for public detail view.
// Corresponds to PHP 'publicDetailFields'.
type UserPublicDetailResp struct {
	UID         string `json:"uid"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	SaasID      string `json:"saasid"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
	Introduce   string `json:"introduce"`
	Birthday    int64  `json:"birthday"`
	Gender      string `json:"gender"`
	GroupID     string `json:"groupid"`
	AreaID      string `json:"areaid"`
	Status      string `json:"status"`
	CreateTime  int64  `json:"createtime"`
	LastTime    int64  `json:"lasttime"`
}

// UserOverviewResp defines fields for brief overview.
// Corresponds to PHP 'overviewFields'.
type UserOverviewResp struct {
	UID         string `json:"uid"`
	Username    string `json:"username"`
	SaasID      string `json:"saasid"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	Introduce   string `json:"introduce"`
	GroupID     string `json:"groupid"`
	AreaID      string `json:"areaid"`
	CreateTime  int64  `json:"createtime"`
	LastTime    int64  `json:"lasttime"`
}

// UserListResp defines fields for admin/internal list.
// Corresponds to PHP 'listFields'.
type UserListResp struct {
	UID         string `json:"uid"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	SaasID      string `json:"saasid"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
	GroupID     string `json:"groupid"`
	Gender      string `json:"gender"`
	AreaID      string `json:"areaid"`
	Status      string `json:"status"`
	CreateTime  int64  `json:"createtime"`
	LastTime    int64  `json:"lasttime"`
}

// UserPublicListResp defines fields for public list.
// Corresponds to PHP 'publicListFields'.
type UserPublicListResp struct {
	UID         string `json:"uid"`
	Username    string `json:"username"`
	SaasID      string `json:"saasid"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	GroupID     string `json:"groupid"`
	Gender      string `json:"gender"`
	AreaID      string `json:"areaid"`
	Status      string `json:"status"`
	CreateTime  int64  `json:"createtime"`
	LastTime    int64  `json:"lasttime"`
}

// UserFilterReq defines allowed filter fields.
// Corresponds to PHP 'filterFields'.
type UserFilterReq struct {
	UID        string `json:"uid,omitempty" form:"uid"`
	Username   string `json:"username,omitempty" form:"username"`
	Email      string `json:"email,omitempty" form:"email"`
	Mobile     string `json:"mobile,omitempty" form:"mobile"`
	SaasID     string `json:"saasid,omitempty" form:"saasid"`
	Nickname   string `json:"nickname,omitempty" form:"nickname"`
	Gender     string `json:"gender,omitempty" form:"gender"`
	GroupID    string `json:"groupid,omitempty" form:"groupid"`
	AreaID     string `json:"areaid,omitempty" form:"areaid"`
	Status     string `json:"status,omitempty" form:"status"`
	CreateTime int64  `json:"createtime,omitempty" form:"createtime"`
	LastTime   int64  `json:"lasttime,omitempty" form:"lasttime"`
}
