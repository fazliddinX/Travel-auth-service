package server

import (
	pb "Auth-service/genproto/auth_service"
	"Auth-service/token"
	"context"
	"fmt"
)

func (s *Server) Register(ctx context.Context, in *pb.RegisterUserReq) (*pb.RegisterUserRes, error) {
	user, err := s.User.Create(in)
	if err != nil {
		s.Logger.Error("Error in register", "error", err)
		return nil, err
	}
	return user, err
}
func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.Tokens, error) {
	user, err := s.User.Login(in)
	if err != nil {
		s.Logger.Error("Error in login", "error", err)
		return nil, err
	}
	AccessToken, err := token.NewAccessToken(user)
	if err != nil {
		s.Logger.Error("Error in get AccessToken", "error", err)
		return nil, err
	}
	RefreshToken, err := token.NewRefreshToken(user)
	if err != nil {
		s.Logger.Error("Error in get RefreshToken", "error", err)
		return nil, err
	}
	token := &pb.Tokens{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
		ExpireIn:     "72",
	}
	return token, nil
}
func (s *Server) GetProfile(ctx context.Context, in *pb.Id) (*pb.Profile, error) {
	user, err := s.User.GetProfile(in)
	if err != nil {
		s.Logger.Error("Error in getProfile", "error", err)
	}
	return user, err
}
func (s *Server) UpdateProfile(ctx context.Context, in *pb.UpdateUser) (*pb.Profile, error) {
	user, err := s.User.UpdateProfile(in)
	if err != nil {
		s.Logger.Error("Error in updateProfile", "error", err)
	}
	return user, err
}
func (s *Server) GetUsers(ctx context.Context, in *pb.FilterGet) (*pb.Users, error) {
	users, err := s.User.GetUsers(in)
	if err != nil {
		s.Logger.Error("Error in getUsers", "error", err)
	}
	return users, err
}
func (s *Server) Delete(ctx context.Context, in *pb.Id) (*pb.Success, error) {
	success, err := s.User.Delete(in)
	if err != nil {
		s.Logger.Error("Error in delete", "error", err)
	}
	return success, err
}
func (s *Server) PasswordRecovery(ctx context.Context, in *pb.Email) (*pb.Success, error) {
	check, err := s.User.PasswordRecovery(in.Email)
	if err != nil {
		s.Logger.Error("Error in passwordRecovery", "error", err)
		return nil, err
	}
	if check {
		msg := fmt.Sprintf("pasword send to %s", in.Email)
		return &pb.Success{Successful: msg}, nil
	}
	return &pb.Success{Successful: "email not found"}, err

}
func (s *Server) TokenRenewal(ctx context.Context, in *pb.RefreshToken) (*pb.Tokens, error) {
	claim, err := token.ExtractClaim(in.RefreshToken)
	if err != nil {
		s.Logger.Error("Error in tokenRenewal", "error", err)
		return nil, err
	}
	AccessToken, err := token.RenewalAccessToken(claim)
	if err != nil {
		s.Logger.Error("Error in tokenRenewal", "error", err)
		return nil, err
	}
	token := &pb.Tokens{
		AccessToken:  AccessToken,
		RefreshToken: in.RefreshToken,
		ExpireIn:     "72",
	}
	return token, nil
}
func (s *Server) Logout(ctx context.Context, in *pb.Id) (*pb.Success, error) {
	return &pb.Success{Successful: "logout"}, nil
}
func (s *Server) ActivityProfile(ctx context.Context, in *pb.Id) (*pb.UserActivities, error) {
	activities, err := s.User.ActivityProfile(in)
	if err != nil {
		s.Logger.Error("Error in activityProfile", "error", err)
		return nil, err
	}

	return activities, err
}
func (s *Server) Follow(ctx context.Context, in *pb.FollowRequest) (*pb.FollowResponse, error) {
	follow, err := s.User.Follow(in)
	if err != nil {
		s.Logger.Error("Error in follow", "error", err)
	}
	return follow, err
}
func (s *Server) GetFollowers(ctx context.Context, in *pb.FilterFollowers) (*pb.Followers, error) {
	followers, err := s.User.GetFollowers(in)
	if err != nil {
		s.Logger.Error("Error in getFollowers", "error", err)
	}
	return followers, err
}
