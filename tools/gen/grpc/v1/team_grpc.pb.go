// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: v1/team.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Team_Get_FullMethodName           = "/v1.team.Team/Get"
	Team_GetUserTeams_FullMethodName  = "/v1.team.Team/GetUserTeams"
	Team_GetAllCanJoin_FullMethodName = "/v1.team.Team/GetAllCanJoin"
	Team_Create_FullMethodName        = "/v1.team.Team/Create"
	Team_Join_FullMethodName          = "/v1.team.Team/Join"
	Team_GetRoles_FullMethodName      = "/v1.team.Team/GetRoles"
	Team_ChangeRole_FullMethodName    = "/v1.team.Team/ChangeRole"
	Team_Leave_FullMethodName         = "/v1.team.Team/Leave"
)

// TeamClient is the client API for Team service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TeamClient interface {
	Get(ctx context.Context, in *GetTeamRequest, opts ...grpc.CallOption) (*TeamResponse, error)
	GetUserTeams(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAllTeamsResponse, error)
	GetAllCanJoin(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAllTeamsResponse, error)
	Create(ctx context.Context, in *CreateTeamRequest, opts ...grpc.CallOption) (*CreateTeamResponse, error)
	Join(ctx context.Context, in *JoinTeamRequest, opts ...grpc.CallOption) (*JoinTeamResponse, error)
	GetRoles(ctx context.Context, in *GetTeamRolesRequest, opts ...grpc.CallOption) (*GetTeamRolesResponse, error)
	ChangeRole(ctx context.Context, in *ChangeTeamRole, opts ...grpc.CallOption) (*Role, error)
	Leave(ctx context.Context, in *LeaveTeamRequest, opts ...grpc.CallOption) (*LeaveTeamResponse, error)
}

type teamClient struct {
	cc grpc.ClientConnInterface
}

func NewTeamClient(cc grpc.ClientConnInterface) TeamClient {
	return &teamClient{cc}
}

func (c *teamClient) Get(ctx context.Context, in *GetTeamRequest, opts ...grpc.CallOption) (*TeamResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TeamResponse)
	err := c.cc.Invoke(ctx, Team_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamClient) GetUserTeams(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAllTeamsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllTeamsResponse)
	err := c.cc.Invoke(ctx, Team_GetUserTeams_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamClient) GetAllCanJoin(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAllTeamsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllTeamsResponse)
	err := c.cc.Invoke(ctx, Team_GetAllCanJoin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamClient) Create(ctx context.Context, in *CreateTeamRequest, opts ...grpc.CallOption) (*CreateTeamResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateTeamResponse)
	err := c.cc.Invoke(ctx, Team_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamClient) Join(ctx context.Context, in *JoinTeamRequest, opts ...grpc.CallOption) (*JoinTeamResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(JoinTeamResponse)
	err := c.cc.Invoke(ctx, Team_Join_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamClient) GetRoles(ctx context.Context, in *GetTeamRolesRequest, opts ...grpc.CallOption) (*GetTeamRolesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTeamRolesResponse)
	err := c.cc.Invoke(ctx, Team_GetRoles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamClient) ChangeRole(ctx context.Context, in *ChangeTeamRole, opts ...grpc.CallOption) (*Role, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Role)
	err := c.cc.Invoke(ctx, Team_ChangeRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamClient) Leave(ctx context.Context, in *LeaveTeamRequest, opts ...grpc.CallOption) (*LeaveTeamResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LeaveTeamResponse)
	err := c.cc.Invoke(ctx, Team_Leave_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TeamServer is the server API for Team service.
// All implementations must embed UnimplementedTeamServer
// for forward compatibility.
type TeamServer interface {
	Get(context.Context, *GetTeamRequest) (*TeamResponse, error)
	GetUserTeams(context.Context, *emptypb.Empty) (*GetAllTeamsResponse, error)
	GetAllCanJoin(context.Context, *emptypb.Empty) (*GetAllTeamsResponse, error)
	Create(context.Context, *CreateTeamRequest) (*CreateTeamResponse, error)
	Join(context.Context, *JoinTeamRequest) (*JoinTeamResponse, error)
	GetRoles(context.Context, *GetTeamRolesRequest) (*GetTeamRolesResponse, error)
	ChangeRole(context.Context, *ChangeTeamRole) (*Role, error)
	Leave(context.Context, *LeaveTeamRequest) (*LeaveTeamResponse, error)
	mustEmbedUnimplementedTeamServer()
}

// UnimplementedTeamServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTeamServer struct{}

func (UnimplementedTeamServer) Get(context.Context, *GetTeamRequest) (*TeamResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedTeamServer) GetUserTeams(context.Context, *emptypb.Empty) (*GetAllTeamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserTeams not implemented")
}
func (UnimplementedTeamServer) GetAllCanJoin(context.Context, *emptypb.Empty) (*GetAllTeamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllCanJoin not implemented")
}
func (UnimplementedTeamServer) Create(context.Context, *CreateTeamRequest) (*CreateTeamResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedTeamServer) Join(context.Context, *JoinTeamRequest) (*JoinTeamResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedTeamServer) GetRoles(context.Context, *GetTeamRolesRequest) (*GetTeamRolesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoles not implemented")
}
func (UnimplementedTeamServer) ChangeRole(context.Context, *ChangeTeamRole) (*Role, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeRole not implemented")
}
func (UnimplementedTeamServer) Leave(context.Context, *LeaveTeamRequest) (*LeaveTeamResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leave not implemented")
}
func (UnimplementedTeamServer) mustEmbedUnimplementedTeamServer() {}
func (UnimplementedTeamServer) testEmbeddedByValue()              {}

// UnsafeTeamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TeamServer will
// result in compilation errors.
type UnsafeTeamServer interface {
	mustEmbedUnimplementedTeamServer()
}

func RegisterTeamServer(s grpc.ServiceRegistrar, srv TeamServer) {
	// If the following call pancis, it indicates UnimplementedTeamServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Team_ServiceDesc, srv)
}

func _Team_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Team_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServer).Get(ctx, req.(*GetTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Team_GetUserTeams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServer).GetUserTeams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Team_GetUserTeams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServer).GetUserTeams(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Team_GetAllCanJoin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServer).GetAllCanJoin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Team_GetAllCanJoin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServer).GetAllCanJoin(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Team_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Team_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServer).Create(ctx, req.(*CreateTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Team_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Team_Join_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServer).Join(ctx, req.(*JoinTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Team_GetRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTeamRolesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServer).GetRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Team_GetRoles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServer).GetRoles(ctx, req.(*GetTeamRolesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Team_ChangeRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeTeamRole)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServer).ChangeRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Team_ChangeRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServer).ChangeRole(ctx, req.(*ChangeTeamRole))
	}
	return interceptor(ctx, in, info, handler)
}

func _Team_Leave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServer).Leave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Team_Leave_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServer).Leave(ctx, req.(*LeaveTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Team_ServiceDesc is the grpc.ServiceDesc for Team service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Team_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.team.Team",
	HandlerType: (*TeamServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Team_Get_Handler,
		},
		{
			MethodName: "GetUserTeams",
			Handler:    _Team_GetUserTeams_Handler,
		},
		{
			MethodName: "GetAllCanJoin",
			Handler:    _Team_GetAllCanJoin_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _Team_Create_Handler,
		},
		{
			MethodName: "Join",
			Handler:    _Team_Join_Handler,
		},
		{
			MethodName: "GetRoles",
			Handler:    _Team_GetRoles_Handler,
		},
		{
			MethodName: "ChangeRole",
			Handler:    _Team_ChangeRole_Handler,
		},
		{
			MethodName: "Leave",
			Handler:    _Team_Leave_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/team.proto",
}
