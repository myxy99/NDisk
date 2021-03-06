/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/12 14:48
 **/
package middleware

import (
	"errors"
	"github.com/coder2z/g-server/xtrace"
	R "github.com/coder2z/ndisk/pkg/response"

	"github.com/coder2z/g-saber/xcast"
	xclient "github.com/coder2z/ndisk/internal/getway/client"
	_map "github.com/coder2z/ndisk/internal/getway/map"
	AuthorityPb "github.com/coder2z/ndisk/pkg/pb/authority"
	NUserPb "github.com/coder2z/ndisk/pkg/pb/nuser"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span, context := xtrace.StartSpanFromContext(ctx.Request.Context(), "Auth Middleware")
		defer span.Finish()
		token := ctx.GetHeader("Authorization")
		if token == "" {
			R.HandleForbidden(ctx)
			ctx.Abort()
			return
		}
		ctx.Set("token", token)
		userInfo, err := xclient.NUserServer.VerifyUsers(context, &NUserPb.Token{
			AccountToken: token,
		})
		if err != nil {
			R.HandleForbidden(ctx)
			ctx.Abort()
			return
		}
		rolesData, err := xclient.AuthorityServer.GetUsersRoles(context, &AuthorityPb.Ids{
			To: []uint32{xcast.ToUint32(userInfo.Uid)},
		})
		if !errors.Is(err, nil) {
			R.HandleForbidden(ctx)
			ctx.Abort()
			return
		}
		var info = _map.UserInfo{
			Uid:         userInfo.Uid,
			Name:        userInfo.Name,
			Alias:       userInfo.Alias,
			Tel:         userInfo.Tel,
			Email:       userInfo.Email,
			Authority:   rolesData.Data[xcast.ToUint32(userInfo.Uid)],
			Status:      userInfo.Status,
			EmailStatus: userInfo.EmailStatus,
			CreatedAt:   userInfo.CreatedAt,
			UpdatedAt:   userInfo.UpdatedAt,
		}
		ctx.Set("user", info)
		ctx.Request = ctx.Request.WithContext(context)
		ctx.Next()
		return
	}
}

func Authority() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span, context := xtrace.StartSpanFromContext(ctx.Request.Context(), "Authority Middleware")
		defer span.Finish()
		if i, ok := ctx.Get("user"); ok {
			info := i.(_map.UserInfo)
			ctx.Set("Uid", info.Uid)
			rep, _ := xclient.AuthorityServer.Enforce(context, &AuthorityPb.Resources{
				Role:   xcast.ToString(info.Uid),
				Obj:    ctx.FullPath(),
				Action: ctx.Request.Method,
			})
			if rep != nil && rep.Ok {
				ctx.Next()
			} else {
				R.HandleForbidden(ctx)
				ctx.Abort()
			}
		} else {
			R.HandleForbidden(ctx)
			ctx.Abort()
		}
		return
	}
}
