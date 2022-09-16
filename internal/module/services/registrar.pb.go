// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: pb/terrarium/module/services/registrar.proto

package services

import (
	module "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ListModulesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Module string `protobuf:"bytes,1,opt,name=module,proto3" json:"module,omitempty"`
}

func (x *ListModulesRequest) Reset() {
	*x = ListModulesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_module_services_registrar_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListModulesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListModulesRequest) ProtoMessage() {}

func (x *ListModulesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_module_services_registrar_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListModulesRequest.ProtoReflect.Descriptor instead.
func (*ListModulesRequest) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_module_services_registrar_proto_rawDescGZIP(), []int{0}
}

func (x *ListModulesRequest) GetModule() string {
	if x != nil {
		return x.Module
	}
	return ""
}

type ListModulesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Provider    string `protobuf:"bytes,2,opt,name=provider,proto3" json:"provider,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	SourceUrl   string `protobuf:"bytes,4,opt,name=source_url,json=sourceUrl,proto3" json:"source_url,omitempty"`
	Maturity    string `protobuf:"bytes,5,opt,name=maturity,proto3" json:"maturity,omitempty"` // string organization = 6 ??
}

func (x *ListModulesResponse) Reset() {
	*x = ListModulesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_module_services_registrar_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListModulesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListModulesResponse) ProtoMessage() {}

func (x *ListModulesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_module_services_registrar_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListModulesResponse.ProtoReflect.Descriptor instead.
func (*ListModulesResponse) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_module_services_registrar_proto_rawDescGZIP(), []int{1}
}

func (x *ListModulesResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ListModulesResponse) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

func (x *ListModulesResponse) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ListModulesResponse) GetSourceUrl() string {
	if x != nil {
		return x.SourceUrl
	}
	return ""
}

func (x *ListModulesResponse) GetMaturity() string {
	if x != nil {
		return x.Maturity
	}
	return ""
}

var File_pb_terrarium_module_services_registrar_proto protoreflect.FileDescriptor

var file_pb_terrarium_module_services_registrar_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x70, 0x62, 0x2f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2f, 0x6d,
	0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19,
	0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x1a, 0x20, 0x70, 0x62, 0x2f, 0x74, 0x65,
	0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2f, 0x6d,
	0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2c, 0x0a, 0x12, 0x4c,
	0x69, 0x73, 0x74, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x22, 0xa2, 0x01, 0x0a, 0x13, 0x4c, 0x69,
	0x73, 0x74, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x75, 0x72,
	0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x55,
	0x72, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x61, 0x74, 0x75, 0x72, 0x69, 0x74, 0x79, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x61, 0x74, 0x75, 0x72, 0x69, 0x74, 0x79, 0x32, 0xca,
	0x01, 0x0a, 0x09, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x72, 0x12, 0x4f, 0x0a, 0x08,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x27, 0x2e, 0x74, 0x65, 0x72, 0x72, 0x61,
	0x72, 0x69, 0x75, 0x6d, 0x2e, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2e, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1a, 0x2e, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x6d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6c, 0x0a,
	0x0b, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x2d, 0x2e, 0x74,
	0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x6f, 0x64,
	0x75, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x74, 0x65,
	0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x6f, 0x64, 0x75,
	0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x4b, 0x5a, 0x49, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72,
	0x69, 0x75, 0x6d, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69,
	0x75, 0x6d, 0x2d, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_terrarium_module_services_registrar_proto_rawDescOnce sync.Once
	file_pb_terrarium_module_services_registrar_proto_rawDescData = file_pb_terrarium_module_services_registrar_proto_rawDesc
)

func file_pb_terrarium_module_services_registrar_proto_rawDescGZIP() []byte {
	file_pb_terrarium_module_services_registrar_proto_rawDescOnce.Do(func() {
		file_pb_terrarium_module_services_registrar_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_terrarium_module_services_registrar_proto_rawDescData)
	})
	return file_pb_terrarium_module_services_registrar_proto_rawDescData
}

var file_pb_terrarium_module_services_registrar_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pb_terrarium_module_services_registrar_proto_goTypes = []interface{}{
	(*ListModulesRequest)(nil),           // 0: terrarium.module.services.ListModulesRequest
	(*ListModulesResponse)(nil),          // 1: terrarium.module.services.ListModulesResponse
	(*module.RegisterModuleRequest)(nil), // 2: terrarium.module.RegisterModuleRequest
	(*module.Response)(nil),              // 3: terrarium.module.Response
}
var file_pb_terrarium_module_services_registrar_proto_depIdxs = []int32{
	2, // 0: terrarium.module.services.Registrar.Register:input_type -> terrarium.module.RegisterModuleRequest
	0, // 1: terrarium.module.services.Registrar.ListModules:input_type -> terrarium.module.services.ListModulesRequest
	3, // 2: terrarium.module.services.Registrar.Register:output_type -> terrarium.module.Response
	1, // 3: terrarium.module.services.Registrar.ListModules:output_type -> terrarium.module.services.ListModulesResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pb_terrarium_module_services_registrar_proto_init() }
func file_pb_terrarium_module_services_registrar_proto_init() {
	if File_pb_terrarium_module_services_registrar_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_terrarium_module_services_registrar_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListModulesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_terrarium_module_services_registrar_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListModulesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pb_terrarium_module_services_registrar_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_terrarium_module_services_registrar_proto_goTypes,
		DependencyIndexes: file_pb_terrarium_module_services_registrar_proto_depIdxs,
		MessageInfos:      file_pb_terrarium_module_services_registrar_proto_msgTypes,
	}.Build()
	File_pb_terrarium_module_services_registrar_proto = out.File
	file_pb_terrarium_module_services_registrar_proto_rawDesc = nil
	file_pb_terrarium_module_services_registrar_proto_goTypes = nil
	file_pb_terrarium_module_services_registrar_proto_depIdxs = nil
}
