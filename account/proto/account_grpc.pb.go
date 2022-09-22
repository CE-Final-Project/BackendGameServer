// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: account.proto

package accountService

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AccountServiceClient is the client API for AccountService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountServiceClient interface {
	CreateAccount(ctx context.Context, in *CreateAccountReq, opts ...grpc.CallOption) (*CreateAccountRes, error)
	UpdateAccount(ctx context.Context, in *UpdateAccountReq, opts ...grpc.CallOption) (*UpdateAccountRes, error)
	GetAccountById(ctx context.Context, in *GetAccountByIdReq, opts ...grpc.CallOption) (*GetAccountByIdRes, error)
	ChangePassword(ctx context.Context, in *ChangePasswordReq, opts ...grpc.CallOption) (*ChangePasswordRes, error)
	BanAccountById(ctx context.Context, in *BanAccountByIdReq, opts ...grpc.CallOption) (*BanAccountByIdRes, error)
	SearchAccount(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchRes, error)
	DeleteAccountById(ctx context.Context, in *DeleteAccountByIdReq, opts ...grpc.CallOption) (*DeleteAccountByIdRes, error)
}

type accountServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountServiceClient(cc grpc.ClientConnInterface) AccountServiceClient {
	return &accountServiceClient{cc}
}

func (c *accountServiceClient) CreateAccount(ctx context.Context, in *CreateAccountReq, opts ...grpc.CallOption) (*CreateAccountRes, error) {
	out := new(CreateAccountRes)
	err := c.cc.Invoke(ctx, "/account.accountService/CreateAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) UpdateAccount(ctx context.Context, in *UpdateAccountReq, opts ...grpc.CallOption) (*UpdateAccountRes, error) {
	out := new(UpdateAccountRes)
	err := c.cc.Invoke(ctx, "/account.accountService/UpdateAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetAccountById(ctx context.Context, in *GetAccountByIdReq, opts ...grpc.CallOption) (*GetAccountByIdRes, error) {
	out := new(GetAccountByIdRes)
	err := c.cc.Invoke(ctx, "/account.accountService/GetAccountById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) ChangePassword(ctx context.Context, in *ChangePasswordReq, opts ...grpc.CallOption) (*ChangePasswordRes, error) {
	out := new(ChangePasswordRes)
	err := c.cc.Invoke(ctx, "/account.accountService/ChangePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) BanAccountById(ctx context.Context, in *BanAccountByIdReq, opts ...grpc.CallOption) (*BanAccountByIdRes, error) {
	out := new(BanAccountByIdRes)
	err := c.cc.Invoke(ctx, "/account.accountService/BanAccountById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) SearchAccount(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchRes, error) {
	out := new(SearchRes)
	err := c.cc.Invoke(ctx, "/account.accountService/SearchAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) DeleteAccountById(ctx context.Context, in *DeleteAccountByIdReq, opts ...grpc.CallOption) (*DeleteAccountByIdRes, error) {
	out := new(DeleteAccountByIdRes)
	err := c.cc.Invoke(ctx, "/account.accountService/DeleteAccountById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountServiceServer is the server API for AccountService service.
// All implementations should embed UnimplementedAccountServiceServer
// for forward compatibility
type AccountServiceServer interface {
	CreateAccount(context.Context, *CreateAccountReq) (*CreateAccountRes, error)
	UpdateAccount(context.Context, *UpdateAccountReq) (*UpdateAccountRes, error)
	GetAccountById(context.Context, *GetAccountByIdReq) (*GetAccountByIdRes, error)
	ChangePassword(context.Context, *ChangePasswordReq) (*ChangePasswordRes, error)
	BanAccountById(context.Context, *BanAccountByIdReq) (*BanAccountByIdRes, error)
	SearchAccount(context.Context, *SearchReq) (*SearchRes, error)
	DeleteAccountById(context.Context, *DeleteAccountByIdReq) (*DeleteAccountByIdRes, error)
}

// UnimplementedAccountServiceServer should be embedded to have forward compatible implementations.
type UnimplementedAccountServiceServer struct {
}

func (UnimplementedAccountServiceServer) CreateAccount(context.Context, *CreateAccountReq) (*CreateAccountRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
func (UnimplementedAccountServiceServer) UpdateAccount(context.Context, *UpdateAccountReq) (*UpdateAccountRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAccount not implemented")
}
func (UnimplementedAccountServiceServer) GetAccountById(context.Context, *GetAccountByIdReq) (*GetAccountByIdRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccountById not implemented")
}
func (UnimplementedAccountServiceServer) ChangePassword(context.Context, *ChangePasswordReq) (*ChangePasswordRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (UnimplementedAccountServiceServer) BanAccountById(context.Context, *BanAccountByIdReq) (*BanAccountByIdRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BanAccountById not implemented")
}
func (UnimplementedAccountServiceServer) SearchAccount(context.Context, *SearchReq) (*SearchRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchAccount not implemented")
}
func (UnimplementedAccountServiceServer) DeleteAccountById(context.Context, *DeleteAccountByIdReq) (*DeleteAccountByIdRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAccountById not implemented")
}

// UnsafeAccountServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountServiceServer will
// result in compilation errors.
type UnsafeAccountServiceServer interface {
	mustEmbedUnimplementedAccountServiceServer()
}

func RegisterAccountServiceServer(s grpc.ServiceRegistrar, srv AccountServiceServer) {
	s.RegisterService(&AccountService_ServiceDesc, srv)
}

func _AccountService_CreateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).CreateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account.accountService/CreateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).CreateAccount(ctx, req.(*CreateAccountReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_UpdateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAccountReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).UpdateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account.accountService/UpdateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).UpdateAccount(ctx, req.(*UpdateAccountReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetAccountById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountByIdReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetAccountById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account.accountService/GetAccountById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetAccountById(ctx, req.(*GetAccountByIdReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_ChangePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePasswordReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).ChangePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account.accountService/ChangePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).ChangePassword(ctx, req.(*ChangePasswordReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_BanAccountById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BanAccountByIdReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).BanAccountById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account.accountService/BanAccountById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).BanAccountById(ctx, req.(*BanAccountByIdReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_SearchAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).SearchAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account.accountService/SearchAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).SearchAccount(ctx, req.(*SearchReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_DeleteAccountById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAccountByIdReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).DeleteAccountById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account.accountService/DeleteAccountById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).DeleteAccountById(ctx, req.(*DeleteAccountByIdReq))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountService_ServiceDesc is the grpc.ServiceDesc for AccountService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccountService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "account.accountService",
	HandlerType: (*AccountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAccount",
			Handler:    _AccountService_CreateAccount_Handler,
		},
		{
			MethodName: "UpdateAccount",
			Handler:    _AccountService_UpdateAccount_Handler,
		},
		{
			MethodName: "GetAccountById",
			Handler:    _AccountService_GetAccountById_Handler,
		},
		{
			MethodName: "ChangePassword",
			Handler:    _AccountService_ChangePassword_Handler,
		},
		{
			MethodName: "BanAccountById",
			Handler:    _AccountService_BanAccountById_Handler,
		},
		{
			MethodName: "SearchAccount",
			Handler:    _AccountService_SearchAccount_Handler,
		},
		{
			MethodName: "DeleteAccountById",
			Handler:    _AccountService_DeleteAccountById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "account.proto",
}
