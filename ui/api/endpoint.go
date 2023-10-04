// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/mainflux/mainflux/pkg/errors"
	"github.com/ultravioletrs/mainflux-ui/ui"

	"github.com/go-kit/kit/endpoint"
	sdk "github.com/mainflux/mainflux/pkg/sdk/go"
)

func indexEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(indexReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.Index(req.token)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func loginEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		res, err := svc.Login()
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func showUpdatePasswordEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		res, err := svc.PasswordUpdate()
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func updatePasswordEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateUserPasswordReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		res, err := svc.UpdatePassword(req.token, req.OldPass, req.NewPass)
		if err != nil {
			return nil, err
		}

		cookies := []*http.Cookie{
			{
				Name:     "token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			},
			{
				Name:     "refresh_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			},
		}

		return uiRes{
			code:    http.StatusOK,
			html:    res,
			cookies: cookies,
		}, nil
	}
}

func tokenEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(tokenReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		credentials := sdk.Credentials{
			Identity: req.Identity,
			Secret:   req.Secret,
		}
		user := sdk.User{
			Credentials: credentials,
		}

		token, err := svc.Token(user)
		if err != nil {
			return nil, errors.Wrap(errUnauthorized, err)
		}

		cookies := []*http.Cookie{
			{
				Name:     "token",
				Value:    token.AccessToken,
				Path:     "/",
				HttpOnly: true,
			},
			{
				Name:     "refresh_token",
				Value:    token.RefreshToken,
				Path:     "/",
				HttpOnly: true,
			},
		}

		tkr := uiRes{
			code:    http.StatusFound,
			cookies: cookies,
			headers: map[string]string{"Location": "/"},
		}

		return tkr, nil
	}
}

func refreshTokenEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(refreshTokenReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		token, err := svc.RefreshToken(req.refreshToken)
		if err != nil {
			return nil, errors.Wrap(errUnauthorized, err)
		}

		cookies := []*http.Cookie{
			{
				Name:     "token",
				Value:    token.AccessToken,
				Path:     "/",
				HttpOnly: true,
			},
			{
				Name:     "refresh_token",
				Value:    token.RefreshToken,
				Path:     "/",
				HttpOnly: true,
			},
		}

		tkr := uiRes{
			code:    http.StatusSeeOther,
			headers: map[string]string{"Location": req.ref},
			cookies: cookies,
		}

		return tkr, nil
	}
}

func logoutEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		res, err := svc.Logout()
		if err != nil {
			return nil, err
		}

		cookies := []*http.Cookie{
			{
				Name:     "token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			},
			{
				Name:     "refresh_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			},
		}
		return uiRes{
			code:    http.StatusFound,
			html:    res,
			cookies: cookies,
			headers: map[string]string{"Location": "/login"},
		}, nil
	}
}

func passwordResetRequestEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(passwordResetRequestReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.PasswordResetRequest(req.Email)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code:    http.StatusFound,
			html:    res,
			headers: map[string]string{"Location": "/login"},
		}, nil
	}
}

func passwordResetEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(passwordResetReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		res, err := svc.PasswordReset(req.token, req.Password, req.ConfirmPassword)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code:    http.StatusFound,
			html:    res,
			headers: map[string]string{"Location": "/login"},
		}, nil
	}
}

func showPasswordResetEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		res, err := svc.ShowPasswordReset()
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func createUserEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createUserReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.CreateUsers(req.token, req.User)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}

func createUsersEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createUsersReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		users := []sdk.User{}
		for i := range req.Emails {
			credential := sdk.Credentials{
				Identity: req.Emails[i],
				Secret:   req.Passwords[i],
			}
			user := sdk.User{
				Name:        req.Names[i],
				Credentials: credential,
			}
			users = append(users, user)
		}
		res, err := svc.CreateUsers(req.token, users...)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}

func listUsersEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listUsersReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.ListUsers(req.token)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func viewUserEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		res, err := svc.ViewUser(req.token, req.id)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func updateUserEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateUserReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		user := sdk.User{
			ID:       req.id,
			Name:     req.Name,
			Metadata: req.Metadata,
		}

		res, err := svc.UpdateUser(req.token, req.id, user)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func updateUserTagsEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateUserTagsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		user := sdk.User{
			ID:   req.id,
			Tags: req.Tags,
		}
		res, err := svc.UpdateUserTags(req.token, req.id, user)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func updateUserIdentityEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateUserIdentityReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		credential := sdk.Credentials{
			Identity: req.Identity,
		}
		user := sdk.User{
			ID:          req.id,
			Credentials: credential,
		}
		res, err := svc.UpdateUserIdentity(req.token, req.id, user)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func enableUserEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateUserStatusReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.EnableUser(req.token, req.UserID)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code:    http.StatusFound,
			html:    res,
			headers: map[string]string{"Location": "/users"},
		}, nil
	}
}

func disableUserEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateUserStatusReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.DisableUser(req.token, req.UserID)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code:    http.StatusFound,
			html:    res,
			headers: map[string]string{"Location": "/users"},
		}, nil
	}
}

func createThingEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createThingReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.CreateThings(req.token, req.Thing)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}

func createThingsEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createThingsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		things := []sdk.Thing{}
		for i := range req.Names {
			th := sdk.Thing{
				Name: req.Names[i],
			}
			things = append(things, th)
		}
		res, err := svc.CreateThings(req.token, things...)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}

func listThingsEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listThingsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.ListThings(req.token)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func viewThingEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		res, err := svc.ViewThing(req.token, req.id)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func updateThingEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateThingReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		uth := sdk.Thing{
			ID:       req.id,
			Name:     req.Name,
			Metadata: req.Metadata,
		}

		res, err := svc.UpdateThing(req.token, req.id, uth)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func updateThingTagsEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateThingTagsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		thing := sdk.Thing{
			ID:   req.id,
			Tags: req.Tags,
		}
		res, err := svc.UpdateThingTags(req.token, req.id, thing)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func updateThingSecretEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateThingSecretReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		res, err := svc.UpdateThingSecret(req.token, req.id, req.Secret)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func enableThingEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateThingStatusReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.EnableThing(req.token, req.ThingID)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code:    http.StatusFound,
			html:    res,
			headers: map[string]string{"Location": "/things"},
		}, nil
	}
}

func disableThingEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateThingStatusReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.DisableThing(req.token, req.ThingID)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code:    http.StatusFound,
			html:    res,
			headers: map[string]string{"Location": "/things"},
		}, nil
	}
}

func updateThingOwnerEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateThingOwnerReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		thing := sdk.Thing{
			ID:    req.id,
			Owner: req.Owner,
		}
		res, err := svc.UpdateThingOwner(req.token, req.id, thing)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func createChannelEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createChannelReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		res, err := svc.CreateChannels(req.token, req.Channel)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}

func createChannelsEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createChannelsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		channels := []sdk.Channel{}
		for i := range req.Names {
			ch := sdk.Channel{
				Name: req.Names[i],
			}
			channels = append(channels, ch)
		}
		res, err := svc.CreateChannels(req.token, channels...)
		if err != nil {
			return nil, err
		}
		return uiRes{
			html: res,
		}, nil
	}
}

func viewChannelEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		res, err := svc.ViewChannel(req.token, req.id)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func updateChannelEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateChannelReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		uch := sdk.Channel{
			ID:          req.id,
			Name:        req.Name,
			Metadata:    req.Metadata,
			Description: req.Description,
		}

		res, err := svc.UpdateChannel(req.token, req.id, uch)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func listChannelsEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listChannelsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		res, err := svc.ListChannels(req.token)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func connectEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		cr := request.(connectReq)

		if err := cr.validate(); err != nil {
			return nil, err
		}

		res, err := svc.Connect(cr.token, cr.ConnIDs)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}

func disconnectEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		dcr := request.(disconnectReq)

		if err := dcr.validate(); err != nil {
			return nil, err
		}

		res, err := svc.Disconnect(dcr.token, dcr.ConnIDs)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func connectThingEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		cr := request.(connectThingReq)

		if err := cr.validate(); err != nil {
			return nil, err
		}

		res, err := svc.ConnectThing(cr.token, cr.ConnIDs)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}

func shareThingEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		sr := request.(shareThingReq)

		if err := sr.validate(); err != nil {
			return nil, err
		}

		res, err := svc.ShareThing(sr.token, sr.ChanID, sr.UserID, sr.Actions)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, err
	}
}

func connectChannelEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		cr := request.(connectChannelReq)

		if err := cr.validate(); err != nil {
			return nil, err
		}

		res, err := svc.ConnectChannel(cr.token, cr.ConnIDs)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}

func disconnectThingEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		dcr := request.(disconnectThingReq)

		if err := dcr.validate(); err != nil {
			return nil, err
		}

		res, err := svc.DisconnectThing(dcr.ThingID, dcr.ChanID, dcr.token)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func disconnectChannelEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		dcr := request.(disconnectChannelReq)

		if err := dcr.validate(); err != nil {
			return nil, err
		}

		res, err := svc.DisconnectChannel(dcr.ThingID, dcr.ChanID, dcr.token)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func listThingsByChannelEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		res, err := svc.ListThingsByChannel(req.token, req.id)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func listChannelsByThingEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		res, err := svc.ListChannelsByThing(req.token, req.id)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func enableChannelEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateChannelStatusReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.EnableChannel(req.token, req.ChannelID)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code:    http.StatusFound,
			html:    res,
			headers: map[string]string{"Location": "/channels"},
		}, err
	}
}

func disableChannelEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateChannelStatusReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.DisableChannel(req.token, req.ChannelID)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code:    http.StatusFound,
			html:    res,
			headers: map[string]string{"Location": "/channels"},
		}, err
	}
}

func listThingsPoliciesEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listPoliciesReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.ListThingsPolicies(req.token)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func addThingsPolicyEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(addThingsPolicyReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.AddThingsPolicy(req.token, req.Policy)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}

func updateThingsPolicyEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updatePolicyReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.UpdateThingsPolicy(req.token, req.Policy)
		if err != nil {
			return nil, err
		}
		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func deleteThingsPolicyEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteThingsPolicyReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.DeleteThingsPolicy(req.token, req.Policy)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func createGroupEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		cgr := request.(createGroupReq)

		if err := cgr.validate(); err != nil {
			return nil, err
		}

		res, err := svc.CreateGroups(cgr.token, cgr.Group)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}

func createGroupsEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createGroupsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		groups := []sdk.Group{}
		for i := range req.Names {
			gr := sdk.Group{
				Name: req.Names[i],
			}
			groups = append(groups, gr)
		}
		res, err := svc.CreateGroups(req.token, groups...)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}

func listGroupsEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listGroupsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		res, err := svc.ListGroups(req.token)
		if err != nil {
			return nil, err
		}
		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func listGroupMembersEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		res, err := svc.ListGroupMembers(req.token, req.id)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func viewGroupEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		res, err := svc.ViewGroup(req.token, req.id)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func updateGroupEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateGroupReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		ugr := sdk.Group{
			ID:          req.id,
			Name:        req.Name,
			Metadata:    req.Metadata,
			Description: req.Description,
		}

		res, err := svc.UpdateGroup(req.token, req.id, ugr)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func assignEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(assignReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.Assign(req.token, req.groupID, req.MemberID, req.Type)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}

func unassignEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(unassignReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		res, err := svc.Unassign(req.token, req.groupID, req.MemberID)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func enableGroupEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateGroupStatusReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.EnableGroup(req.token, req.GroupID)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code:    http.StatusFound,
			html:    res,
			headers: map[string]string{"Location": "/groups"},
		}, nil
	}
}

func disableGroupEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateGroupStatusReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.DisableGroup(req.token, req.GroupID)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code:    http.StatusFound,
			html:    res,
			headers: map[string]string{"Location": "/groups"},
		}, nil
	}
}

func listPoliciesEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listPoliciesReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.ListPolicies(req.token)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func addPolicyEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(addPolicyReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.AddPolicy(req.token, req.Policy)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}

func updatePolicyEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updatePolicyReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.UpdatePolicy(req.token, req.Policy)
		if err != nil {
			return nil, err
		}
		return uiRes{
			code:    http.StatusFound,
			html:    res,
			headers: map[string]string{"Location": "/users/policies"},
		}, nil
	}
}

func deletePolicyEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(deletePolicyReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.DeletePolicy(req.token, req.Policy)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code:    http.StatusFound,
			html:    res,
			headers: map[string]string{"Location": "/users/policies"},
		}, nil
	}
}

func publishMessageEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(publishReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.Publish(req.token, req.thingKey, req.Msg)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func readMessageEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(readMessageReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.ReadMessage(req.token)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func wsConnectionEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(wsConnectionReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.WsConnection(req.token, req.ChanID, req.ThingKey)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}

func listDeletedClientsEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listDeletedClientsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.ListDeletedClients(req.token)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func getTerminalEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.GetRemoteTerminal(req.id)
		if err != nil {
			return nil, err
		}
		return uiRes{
			html: res,
		}, nil
	}
}

func handleTerminalInputEndpoint(svc ui.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(bootstrapCommandReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		// Create a channel to receive the command result
		ch := make(chan string)

		g, ctx := errgroup.WithContext(ctx)

		// Start a goroutine to process the command asynchronously
		g.Go(func() error {
			return svc.ProcessTerminalCommand(ctx, req.id, req.token, req.command, ch)
		})

		if err := g.Wait(); err != nil {
			return nil, err
		}

		// Receive the command result from the channel
		result := <-ch

		return terminalResponse{
			Command: req.command,
			Result:  result,
		}, nil
	}
}

func viewBootstrap(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.ViewBootstrap(req.token, req.id)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func listBootstrap(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listBootstrapReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.ListBootstrap(req.token)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func updateBootstrap(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateBootstrapReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		cfg := sdk.BootstrapConfig{
			ThingID: req.id,
			Name:    req.Name,
			Content: req.Content,
		}
		res, err := svc.UpdateBootstrap(req.token, cfg)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func updateBootstrapCerts(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateBootstrapCertReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		cfg := sdk.BootstrapConfig{
			ThingID:    req.thingID,
			ClientCert: req.ClientCert,
			ClientKey:  req.ClientKey,
			CACert:     req.CACert,
		}
		res, err := svc.UpdateBootstrapCerts(req.token, cfg)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func updateBootstrapConnections(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateBootstrapConnReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		cfg := sdk.BootstrapConfig{
			ThingID:  req.id,
			Channels: req.Channels,
		}
		res, err := svc.UpdateBootstrapConnections(req.token, cfg)
		if err != nil {
			return nil, err
		}

		return uiRes{
			code: http.StatusOK,
			html: res,
		}, nil
	}
}

func createBootstrap(svc ui.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createBootstrapReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		cfg := sdk.BootstrapConfig{
			ThingID:     req.ThingID,
			Channels:    req.Channels,
			ExternalID:  req.ExternalID,
			ExternalKey: req.ExternalKey,
			Name:        req.Name,
			ClientCert:  req.ClientCert,
			ClientKey:   req.ClientKey,
			CACert:      req.CACert,
			Content:     req.Content,
		}
		res, err := svc.CreateBootstrap(req.token, cfg)
		if err != nil {
			return nil, err
		}

		return uiRes{
			html: res,
		}, nil
	}
}
