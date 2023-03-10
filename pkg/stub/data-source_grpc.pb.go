// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: stub/data-source.proto

package stub

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

// DataSourceClient is the client API for DataSource service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataSourceClient interface {
	Create(ctx context.Context, in *CreateDataSourceRequest, opts ...grpc.CallOption) (*CreateDataSourceResponse, error)
	List(ctx context.Context, in *ListDataSourceRequest, opts ...grpc.CallOption) (*ListDataSourceResponse, error)
	Update(ctx context.Context, in *UpdateDataSourceRequest, opts ...grpc.CallOption) (*UpdateDataSourceResponse, error)
	Delete(ctx context.Context, in *DeleteDataSourceRequest, opts ...grpc.CallOption) (*DeleteDataSourceResponse, error)
	Get(ctx context.Context, in *GetDataSourceRequest, opts ...grpc.CallOption) (*GetDataSourceResponse, error)
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	ListEntities(ctx context.Context, in *ListEntitiesRequest, opts ...grpc.CallOption) (*ListEntitiesResponse, error)
	PrepareResourceFromEntity(ctx context.Context, in *PrepareResourceFromEntityRequest, opts ...grpc.CallOption) (*PrepareResourceFromEntityResponse, error)
}

type dataSourceClient struct {
	cc grpc.ClientConnInterface
}

func NewDataSourceClient(cc grpc.ClientConnInterface) DataSourceClient {
	return &dataSourceClient{cc}
}

func (c *dataSourceClient) Create(ctx context.Context, in *CreateDataSourceRequest, opts ...grpc.CallOption) (*CreateDataSourceResponse, error) {
	out := new(CreateDataSourceResponse)
	err := c.cc.Invoke(ctx, "/stub.DataSource/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataSourceClient) List(ctx context.Context, in *ListDataSourceRequest, opts ...grpc.CallOption) (*ListDataSourceResponse, error) {
	out := new(ListDataSourceResponse)
	err := c.cc.Invoke(ctx, "/stub.DataSource/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataSourceClient) Update(ctx context.Context, in *UpdateDataSourceRequest, opts ...grpc.CallOption) (*UpdateDataSourceResponse, error) {
	out := new(UpdateDataSourceResponse)
	err := c.cc.Invoke(ctx, "/stub.DataSource/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataSourceClient) Delete(ctx context.Context, in *DeleteDataSourceRequest, opts ...grpc.CallOption) (*DeleteDataSourceResponse, error) {
	out := new(DeleteDataSourceResponse)
	err := c.cc.Invoke(ctx, "/stub.DataSource/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataSourceClient) Get(ctx context.Context, in *GetDataSourceRequest, opts ...grpc.CallOption) (*GetDataSourceResponse, error) {
	out := new(GetDataSourceResponse)
	err := c.cc.Invoke(ctx, "/stub.DataSource/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataSourceClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/stub.DataSource/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataSourceClient) ListEntities(ctx context.Context, in *ListEntitiesRequest, opts ...grpc.CallOption) (*ListEntitiesResponse, error) {
	out := new(ListEntitiesResponse)
	err := c.cc.Invoke(ctx, "/stub.DataSource/ListEntities", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataSourceClient) PrepareResourceFromEntity(ctx context.Context, in *PrepareResourceFromEntityRequest, opts ...grpc.CallOption) (*PrepareResourceFromEntityResponse, error) {
	out := new(PrepareResourceFromEntityResponse)
	err := c.cc.Invoke(ctx, "/stub.DataSource/PrepareResourceFromEntity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataSourceServer is the server API for DataSource service.
// All implementations must embed UnimplementedDataSourceServer
// for forward compatibility
type DataSourceServer interface {
	Create(context.Context, *CreateDataSourceRequest) (*CreateDataSourceResponse, error)
	List(context.Context, *ListDataSourceRequest) (*ListDataSourceResponse, error)
	Update(context.Context, *UpdateDataSourceRequest) (*UpdateDataSourceResponse, error)
	Delete(context.Context, *DeleteDataSourceRequest) (*DeleteDataSourceResponse, error)
	Get(context.Context, *GetDataSourceRequest) (*GetDataSourceResponse, error)
	Status(context.Context, *StatusRequest) (*StatusResponse, error)
	ListEntities(context.Context, *ListEntitiesRequest) (*ListEntitiesResponse, error)
	PrepareResourceFromEntity(context.Context, *PrepareResourceFromEntityRequest) (*PrepareResourceFromEntityResponse, error)
	mustEmbedUnimplementedDataSourceServer()
}

// UnimplementedDataSourceServer must be embedded to have forward compatible implementations.
type UnimplementedDataSourceServer struct {
}

func (UnimplementedDataSourceServer) Create(context.Context, *CreateDataSourceRequest) (*CreateDataSourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedDataSourceServer) List(context.Context, *ListDataSourceRequest) (*ListDataSourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedDataSourceServer) Update(context.Context, *UpdateDataSourceRequest) (*UpdateDataSourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedDataSourceServer) Delete(context.Context, *DeleteDataSourceRequest) (*DeleteDataSourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedDataSourceServer) Get(context.Context, *GetDataSourceRequest) (*GetDataSourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedDataSourceServer) Status(context.Context, *StatusRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedDataSourceServer) ListEntities(context.Context, *ListEntitiesRequest) (*ListEntitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEntities not implemented")
}
func (UnimplementedDataSourceServer) PrepareResourceFromEntity(context.Context, *PrepareResourceFromEntityRequest) (*PrepareResourceFromEntityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrepareResourceFromEntity not implemented")
}
func (UnimplementedDataSourceServer) mustEmbedUnimplementedDataSourceServer() {}

// UnsafeDataSourceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataSourceServer will
// result in compilation errors.
type UnsafeDataSourceServer interface {
	mustEmbedUnimplementedDataSourceServer()
}

func RegisterDataSourceServer(s grpc.ServiceRegistrar, srv DataSourceServer) {
	s.RegisterService(&DataSource_ServiceDesc, srv)
}

func _DataSource_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDataSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataSourceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stub.DataSource/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataSourceServer).Create(ctx, req.(*CreateDataSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataSource_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDataSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataSourceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stub.DataSource/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataSourceServer).List(ctx, req.(*ListDataSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataSource_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDataSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataSourceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stub.DataSource/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataSourceServer).Update(ctx, req.(*UpdateDataSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataSource_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDataSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataSourceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stub.DataSource/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataSourceServer).Delete(ctx, req.(*DeleteDataSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataSource_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDataSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataSourceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stub.DataSource/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataSourceServer).Get(ctx, req.(*GetDataSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataSource_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataSourceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stub.DataSource/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataSourceServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataSource_ListEntities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEntitiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataSourceServer).ListEntities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stub.DataSource/ListEntities",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataSourceServer).ListEntities(ctx, req.(*ListEntitiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataSource_PrepareResourceFromEntity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrepareResourceFromEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataSourceServer).PrepareResourceFromEntity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stub.DataSource/PrepareResourceFromEntity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataSourceServer).PrepareResourceFromEntity(ctx, req.(*PrepareResourceFromEntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DataSource_ServiceDesc is the grpc.ServiceDesc for DataSource service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataSource_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "stub.DataSource",
	HandlerType: (*DataSourceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _DataSource_Create_Handler,
		},
		{
			MethodName: "List",
			Handler:    _DataSource_List_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _DataSource_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _DataSource_Delete_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _DataSource_Get_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _DataSource_Status_Handler,
		},
		{
			MethodName: "ListEntities",
			Handler:    _DataSource_ListEntities_Handler,
		},
		{
			MethodName: "PrepareResourceFromEntity",
			Handler:    _DataSource_PrepareResourceFromEntity_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stub/data-source.proto",
}
