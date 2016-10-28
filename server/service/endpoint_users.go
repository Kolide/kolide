package service

import (
	"net/http"

	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/kolide"
)

////////////////////////////////////////////////////////////////////////////////
// Create User
////////////////////////////////////////////////////////////////////////////////

type createUserRequest struct {
	payload kolide.UserPayload
}

type createUserResponse struct {
	User *kolide.User `json:"user,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r createUserResponse) error() error { return r.Err }

func makeCreateUserEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createUserRequest)
		user, err := svc.NewUser(ctx, req.payload)
		if err != nil {
			return createUserResponse{Err: err}, nil
		}
		return createUserResponse{User: user}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get User
////////////////////////////////////////////////////////////////////////////////

type getUserRequest struct {
	ID uint `json:"id"`
}

type getUserResponse struct {
	User *kolide.User `json:"user,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r getUserResponse) error() error { return r.Err }

func makeGetUserEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getUserRequest)
		user, err := svc.User(ctx, req.ID)
		if err != nil {
			return getUserResponse{Err: err}, nil
		}
		return getUserResponse{User: user}, nil
	}
}

func makeGetSessionUserEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user, err := svc.AuthenticatedUser(ctx)
		if err != nil {
			return getUserResponse{Err: err}, nil
		}
		return getUserResponse{User: user}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// List Users
////////////////////////////////////////////////////////////////////////////////

type listUsersRequest struct {
	ListOptions kolide.ListOptions
}

type listUsersResponse struct {
	Users []kolide.User `json:"users"`
	Err   error         `json:"error,omitempty"`
}

func (r listUsersResponse) error() error { return r.Err }

func makeListUsersEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listUsersRequest)
		users, err := svc.ListUsers(ctx, req.ListOptions)
		if err != nil {
			return listUsersResponse{Err: err}, nil
		}

		resp := listUsersResponse{Users: []kolide.User{}}
		for _, user := range users {
			resp.Users = append(resp.Users, *user)
		}
		return resp, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Reset Password
////////////////////////////////////////////////////////////////////////////////

type resetPasswordRequest struct {
	PasswordResetToken string `json:"password_reset_token"`
	NewPassword        string `json:"new_password"`
}

type resetPasswordResponse struct {
	Err error `json:"error,omitempty"`
}

func (r resetPasswordResponse) error() error { return r.Err }

func makeResetPasswordEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(resetPasswordRequest)
		err := svc.ResetPassword(ctx, req.PasswordResetToken, req.NewPassword)
		return resetPasswordResponse{Err: err}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Modify User
////////////////////////////////////////////////////////////////////////////////

type modifyUserRequest struct {
	ID      uint
	payload kolide.UserPayload
}

type modifyUserResponse struct {
	User *kolide.User `json:"user,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r modifyUserResponse) error() error { return r.Err }

func makeModifyUserEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyUserRequest)
		user, err := svc.ModifyUser(ctx, req.ID, req.payload)
		if err != nil {
			return modifyUserResponse{Err: err}, nil
		}

		return modifyUserResponse{User: user}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Forgot Password
////////////////////////////////////////////////////////////////////////////////

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

type forgotPasswordResponse struct {
	Err error `json:"error,omitempty"`
}

func (r forgotPasswordResponse) error() error { return r.Err }
func (r forgotPasswordResponse) status() int  { return http.StatusAccepted }

func makeForgotPasswordEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(forgotPasswordRequest)
		err := svc.RequestPasswordReset(ctx, req.Email)
		if err != nil {
			return forgotPasswordResponse{Err: err}, nil
		}
		return forgotPasswordResponse{}, nil
	}
}
