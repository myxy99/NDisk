package admin_server

import (
	"context"
	"errors"
	"github.com/myxy99/component/pkg/xcast"
	xclient "github.com/myxy99/ndisk/internal/getway/client"
	xerror "github.com/myxy99/ndisk/internal/getway/error"
	_map "github.com/myxy99/ndisk/internal/getway/map"
	AuthorityPb "github.com/myxy99/ndisk/pkg/pb/authority"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
	xrpc "github.com/myxy99/ndisk/pkg/rpc"
)

func CreateUser(ctx context.Context, user _map.CreateUser) (data _map.Batch, errs *xerror.Err) {
	var list = make([]*NUserPb.UserInfo, len(user.Data))
	for i, datum := range user.Data {
		list[i] = &NUserPb.UserInfo{
			Name:     datum.Name,
			Alias:    datum.Alias,
			Tel:      datum.Tel,
			Email:    datum.Email,
			Password: datum.Password,
		}
	}
	rep, err := xclient.NUserServer.CreateUsers(ctx, &NUserPb.UserList{
		List: list,
	})
	if !errors.Is(err, nil) {
		e := xerror.NewErrRPC(err)
		if e.ErrorCode == xrpc.DataExistErrCode {
			e = e.SetMessage("data exist")
		}
		return data, e
	}
	return _map.Batch{Count: rep.Count}, nil
}

func UpdateUser(ctx context.Context, user _map.UpdateUser) *xerror.Err {
	_, err := xclient.NUserServer.UpdateUser(ctx, &NUserPb.UserInfo{
		Uid:      user.Uid,
		Name:     user.Name,
		Alias:    user.Alias,
		Tel:      user.Tel,
		Email:    user.Email,
		Password: user.Password,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

func DeleteUser(ctx context.Context, user _map.UidList) (data _map.Batch, errs *xerror.Err) {
	rep, err := xclient.NUserServer.DelUsers(ctx, &NUserPb.UidList{
		Uid: user.List,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	return _map.Batch{Count: rep.Count}, nil
}

func UserList(ctx context.Context, req _map.PageList) (data _map.UserList, errs *xerror.Err) {
	rep, err := xclient.NUserServer.GetUserList(ctx, &NUserPb.PageRequest{
		Keyword:  req.Keyword,
		Page:     xcast.ToUint32(req.Page),
		Limit:    xcast.ToUint32(req.PageSize),
		IsDelete: req.IsDelete,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	uidList := make([]uint32, len(rep.List))
	for i, info := range rep.List {
		uidList[i] = xcast.ToUint32(info.Uid)
	}
	rolesData, err := xclient.AuthorityServer.GetUsersRoles(ctx, &AuthorityPb.Ids{
		To: uidList,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	list := make([]_map.UserInfo, len(rep.List))
	for i, info := range rep.List {
		list[i] = _map.UserInfo{
			Uid:         info.Uid,
			Name:        info.Name,
			Alias:       info.Alias,
			Tel:         info.Tel,
			Authority:   rolesData.Data[xcast.ToUint32(info.Uid)],
			Email:       info.Email,
			Status:      info.Status,
			EmailStatus: info.EmailStatus,
			CreatedAt:   info.CreatedAt,
			UpdatedAt:   info.UpdatedAt,
			DeletedAt:   info.DeletedAt,
		}
	}
	return _map.UserList{
		Count: rep.Count,
		Data:  list,
	}, nil
}

func UpdateStatusUser(ctx context.Context, req _map.UpdateStatus) *xerror.Err {
	_, err := xclient.NUserServer.UpdateUserStatus(ctx, &NUserPb.UserInfo{
		Uid:    req.Uid,
		Status: req.Status,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

func RestoreDeleteUser(ctx context.Context, req _map.UidList) (data _map.Batch, errs *xerror.Err) {
	rep, err := xclient.NUserServer.RecoverDelUsers(ctx, &NUserPb.UidList{
		Uid: req.List,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	return _map.Batch{Count: rep.Count}, nil
}

func UserById(ctx context.Context, req _map.Uid) (data _map.UserInfo, errs *xerror.Err) {
	rep, err := xclient.NUserServer.GetUserById(ctx, &NUserPb.UserInfo{
		Uid: req.Uid,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	rolesData, err := xclient.AuthorityServer.GetUsersRoles(ctx, &AuthorityPb.Ids{
		To: []uint32{xcast.ToUint32(req.Uid)},
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	return _map.UserInfo{
		Uid:         rep.Uid,
		Name:        rep.Name,
		Alias:       rep.Alias,
		Tel:         rep.Tel,
		Email:       rep.Email,
		Authority:   rolesData.Data[xcast.ToUint32(req.Uid)],
		Status:      rep.Status,
		EmailStatus: rep.EmailStatus,
		CreatedAt:   rep.CreatedAt,
		UpdatedAt:   rep.UpdatedAt,
		DeletedAt:   rep.DeletedAt,
	}, nil
}

func GetAllRoles(ctx context.Context, req _map.PageList) (data _map.RolesListRes, errs *xerror.Err) {
	rep, err := xclient.AuthorityServer.GetAllRoles(ctx, &AuthorityPb.PageRequest{
		Keyword:  req.Keyword,
		Page:     xcast.ToUint32(req.Page),
		Limit:    xcast.ToUint32(req.PageSize),
		IsDelete: req.IsDelete,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}

	var list = make([]_map.RolesInfoRes, len(rep.List))
	for i, info := range rep.List {
		list[i] = _map.RolesInfoRes{
			Id:          info.Id,
			Name:        info.Name,
			Description: info.Description,
			CreatedAt:   info.CreatedAt,
			UpdatedAt:   info.UpdatedAt,
			DeletedAt:   info.DeletedAt,
		}
	}

	return _map.RolesListRes{
		Count: rep.Count,
		Data:  list,
	}, nil
}

func UserByRole(ctx context.Context, role string) (data _map.UserList, errs *xerror.Err) {
	rep, err := xclient.AuthorityServer.GetUsersForRole(ctx, &AuthorityPb.Single{
		To: role,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	var uidList = make([]uint32, len(rep.Data))
	for i, datum := range rep.Data {
		uidList[i] = xcast.ToUint32(datum)
	}
	reqUser, err := xclient.NUserServer.GetUserListByUid(ctx, &NUserPb.UidList{
		Uid: uidList,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	list := make([]_map.UserInfo, len(reqUser.List))
	for i, info := range reqUser.List {
		list[i] = _map.UserInfo{
			Uid:       info.Uid,
			Name:      info.Name,
			Alias:     info.Alias,
			Tel:       info.Tel,
			Email:     info.Email,
			CreatedAt: info.CreatedAt,
			UpdatedAt: info.UpdatedAt,
			DeletedAt: info.DeletedAt,
		}
	}
	return _map.UserList{
		Count: reqUser.Count,
		Data:  list,
	}, nil
}

func UserAddRoles(ctx context.Context, req _map.UserRolesReq) (errs *xerror.Err) {
	_, err := xclient.AuthorityServer.AddRolesForUser(ctx, &AuthorityPb.Batch{
		To:      xcast.ToString(req.Uid),
		Operate: req.Role,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

func UserDelRoles(ctx context.Context, req _map.UserRoleReq) (errs *xerror.Err) {
	_, err := xclient.AuthorityServer.DeleteRoleForUser(ctx, &AuthorityPb.Single{
		To:      xcast.ToString(req.Uid),
		Operate: req.Role,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

func MenuList(ctx context.Context, req _map.PageList) (data _map.MenuListRes, errs *xerror.Err) {
	rep, err := xclient.AuthorityServer.GetAllMenu(ctx, &AuthorityPb.PageRequest{
		Keyword:  req.Keyword,
		Page:     xcast.ToUint32(req.Page),
		Limit:    xcast.ToUint32(req.PageSize),
		IsDelete: req.IsDelete,
	})
	if !errors.Is(err, nil) {
		e := xerror.NewErrRPC(err)
		return data, e
	}
	list := make([]_map.MenuInfoRes, len(rep.List))
	for i, info := range rep.List {
		list[i] = _map.MenuInfoRes{
			Id:          info.Id,
			ParentId:    info.ParentId,
			Path:        info.Path,
			Name:        info.Name,
			Description: info.Description,
			IconClass:   info.IconClass,
			CreatedAt:   info.CreatedAt,
			UpdatedAt:   info.UpdatedAt,
			DeletedAt:   info.DeletedAt,
		}
	}
	return _map.MenuListRes{
		Count: rep.Count,
		Data:  list,
	}, nil
}

func DelMenu(ctx context.Context, req _map.UidList) (data _map.Batch, errs *xerror.Err) {
	rep, err := xclient.AuthorityServer.DeleteMenu(ctx, &AuthorityPb.Ids{
		To: req.List,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	return _map.Batch{Count: rep.Count}, nil
}

func AddMenu(ctx context.Context, req _map.MenuReq) (errs *xerror.Err) {
	_, err := xclient.AuthorityServer.AddMenu(ctx, &AuthorityPb.MenuInfo{
		ParentId:    req.ParentId,
		Path:        req.Path,
		Name:        req.Name,
		Description: req.Description,
		IconClass:   req.IconClass,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

func UpdateMenu(ctx context.Context, id uint32, req _map.MenuReq) (errs *xerror.Err) {
	_, err := xclient.AuthorityServer.UpdateMenu(ctx, &AuthorityPb.MenuInfo{
		Id:          id,
		ParentId:    req.ParentId,
		Path:        req.Path,
		Name:        req.Name,
		Description: req.Description,
		IconClass:   req.IconClass,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

func ResourcesList(ctx context.Context, req _map.PageList) (data _map.ResourcesListRes, errs *xerror.Err) {
	rep, err := xclient.AuthorityServer.GetAllResources(ctx, &AuthorityPb.PageRequest{
		Keyword:  req.Keyword,
		Page:     xcast.ToUint32(req.Page),
		Limit:    xcast.ToUint32(req.PageSize),
		IsDelete: req.IsDelete,
	})
	if !errors.Is(err, nil) {
		e := xerror.NewErrRPC(err)
		return data, e
	}
	list := make([]_map.ResourcesInfoRes, len(rep.List))
	for i, info := range rep.List {
		list[i] = _map.ResourcesInfoRes{
			Id:          info.Id,
			Name:        info.Name,
			Path:        info.Path,
			Action:      info.Action,
			Description: info.Description,
			CreatedAt:   info.CreatedAt,
			UpdatedAt:   info.UpdatedAt,
			DeletedAt:   info.DeletedAt,
		}
	}
	return _map.ResourcesListRes{
		Count: rep.Count,
		Data:  list,
	}, nil
}

func DelResources(ctx context.Context, req _map.UidList) (data _map.Batch, errs *xerror.Err) {
	rep, err := xclient.AuthorityServer.DeleteResources(ctx, &AuthorityPb.Ids{
		To: req.List,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	return _map.Batch{Count: rep.Count}, nil
}

func AddResources(ctx context.Context, req _map.ResourcesReq) (errs *xerror.Err) {
	_, err := xclient.AuthorityServer.AddResources(ctx, &AuthorityPb.ResourcesInfo{
		Name:        req.Name,
		Path:        req.Path,
		Action:      req.Action,
		Description: req.Description,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

func UpdateResources(ctx context.Context, id uint32, req _map.ResourcesReq) (errs *xerror.Err) {
	_, err := xclient.AuthorityServer.UpdateResources(ctx, &AuthorityPb.ResourcesInfo{
		Id:          id,
		Name:        req.Name,
		Path:        req.Path,
		Action:      req.Action,
		Description: req.Description,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

// 更新角色下的所有菜单权限
func UpdateRolesMenuAndResources(ctx context.Context, req _map.UpdateRolesMenuAndResourcesReq) (errs *xerror.Err) {
	_, err := xclient.AuthorityServer.UpdateRolesMenuAndResources(ctx, &AuthorityPb.UpdateRolesMenuAndResourcesReq{
		Id:        req.Id,
		Menus:     req.Menus,
		Resources: req.Resources,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

// 删除角色
func DelRoles(ctx context.Context, req _map.UidList) (data _map.Batch, errs *xerror.Err) {
	rep, err := xclient.AuthorityServer.DeleteRoles(ctx, &AuthorityPb.Ids{
		To: req.List,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	return _map.Batch{Count: rep.Count}, nil
}

// 添加角色
func AddRoles(ctx context.Context, req _map.RoleInfoReq) (errs *xerror.Err) {
	_, err := xclient.AuthorityServer.AddRoles(ctx, &AuthorityPb.RolesInfo{
		Name:        req.Name,
		Description: req.Description,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

// 更新角色
func UpdateRoles(ctx context.Context, id uint32, req _map.RoleInfoReq) (errs *xerror.Err) {
	_, err := xclient.AuthorityServer.UpdateRoles(ctx, &AuthorityPb.RolesInfo{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}
