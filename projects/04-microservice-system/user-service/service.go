package main

import (
	"context"

	"learning-go-with-claude/projects/04-microservice-system/proto"
)

// userService implements proto.UserService.
//
// TODO: Implement all methods:
//   - Register: validate input, hash password, store user, return User
//   - Login: look up user by email, verify password, generate JWT, return TokenPair
//   - GetProfile: look up user by ID, return User
//   - UpdateProfile: look up user, apply changes, return updated User

type userService struct {
	// TODO: Add dependencies (store, JWT secret, etc.)
}

// Verify at compile time that userService implements the interface.
var _ proto.UserService = (*userService)(nil)

func (s *userService) Register(ctx context.Context, req proto.RegisterRequest) (proto.User, error) {
	return proto.User{}, nil
}

func (s *userService) Login(ctx context.Context, req proto.LoginRequest) (proto.TokenPair, error) {
	return proto.TokenPair{}, nil
}

func (s *userService) GetProfile(ctx context.Context, userID string) (proto.User, error) {
	return proto.User{}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, req proto.UpdateProfileRequest) (proto.User, error) {
	return proto.User{}, nil
}
