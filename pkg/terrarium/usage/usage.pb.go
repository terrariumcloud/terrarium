// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: pb/terrarium/usage/usage.proto

package usage

import (
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

type RegisterDeploymentUnitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Unit          *DeploymentUnit       `protobuf:"bytes,1,opt,name=unit,proto3" json:"unit,omitempty"`
	Notifications []*NotificationMethod `protobuf:"bytes,2,rep,name=notifications,proto3" json:"notifications,omitempty"`
}

func (x *RegisterDeploymentUnitRequest) Reset() {
	*x = RegisterDeploymentUnitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_usage_usage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterDeploymentUnitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterDeploymentUnitRequest) ProtoMessage() {}

func (x *RegisterDeploymentUnitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_usage_usage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterDeploymentUnitRequest.ProtoReflect.Descriptor instead.
func (*RegisterDeploymentUnitRequest) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_usage_usage_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterDeploymentUnitRequest) GetUnit() *DeploymentUnit {
	if x != nil {
		return x.Unit
	}
	return nil
}

func (x *RegisterDeploymentUnitRequest) GetNotifications() []*NotificationMethod {
	if x != nil {
		return x.Notifications
	}
	return nil
}

type RegisterDeploymentUnitResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RegisterDeploymentUnitResponse) Reset() {
	*x = RegisterDeploymentUnitResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_usage_usage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterDeploymentUnitResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterDeploymentUnitResponse) ProtoMessage() {}

func (x *RegisterDeploymentUnitResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_usage_usage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterDeploymentUnitResponse.ProtoReflect.Descriptor instead.
func (*RegisterDeploymentUnitResponse) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_usage_usage_proto_rawDescGZIP(), []int{1}
}

type NotifyDependencyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Unit         *VersionedDeploymentUnit   `protobuf:"bytes,1,opt,name=unit,proto3" json:"unit,omitempty"`
	Dependencies []*VersionedDeploymentUnit `protobuf:"bytes,2,rep,name=dependencies,proto3" json:"dependencies,omitempty"`
}

func (x *NotifyDependencyRequest) Reset() {
	*x = NotifyDependencyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_usage_usage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyDependencyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyDependencyRequest) ProtoMessage() {}

func (x *NotifyDependencyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_usage_usage_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyDependencyRequest.ProtoReflect.Descriptor instead.
func (*NotifyDependencyRequest) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_usage_usage_proto_rawDescGZIP(), []int{2}
}

func (x *NotifyDependencyRequest) GetUnit() *VersionedDeploymentUnit {
	if x != nil {
		return x.Unit
	}
	return nil
}

func (x *NotifyDependencyRequest) GetDependencies() []*VersionedDeploymentUnit {
	if x != nil {
		return x.Dependencies
	}
	return nil
}

type NotifyDependencyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NotifyDependencyResponse) Reset() {
	*x = NotifyDependencyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_usage_usage_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyDependencyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyDependencyResponse) ProtoMessage() {}

func (x *NotifyDependencyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_usage_usage_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyDependencyResponse.ProtoReflect.Descriptor instead.
func (*NotifyDependencyResponse) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_usage_usage_proto_rawDescGZIP(), []int{3}
}

type DeploymentUnit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type         string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Organization string `protobuf:"bytes,2,opt,name=organization,proto3" json:"organization,omitempty"`
	Name         string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *DeploymentUnit) Reset() {
	*x = DeploymentUnit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_usage_usage_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeploymentUnit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeploymentUnit) ProtoMessage() {}

func (x *DeploymentUnit) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_usage_usage_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeploymentUnit.ProtoReflect.Descriptor instead.
func (*DeploymentUnit) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_usage_usage_proto_rawDescGZIP(), []int{4}
}

func (x *DeploymentUnit) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *DeploymentUnit) GetOrganization() string {
	if x != nil {
		return x.Organization
	}
	return ""
}

func (x *DeploymentUnit) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type VersionedDeploymentUnit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Unit    *DeploymentUnit `protobuf:"bytes,1,opt,name=unit,proto3" json:"unit,omitempty"`
	Version string          `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *VersionedDeploymentUnit) Reset() {
	*x = VersionedDeploymentUnit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_usage_usage_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionedDeploymentUnit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionedDeploymentUnit) ProtoMessage() {}

func (x *VersionedDeploymentUnit) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_usage_usage_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionedDeploymentUnit.ProtoReflect.Descriptor instead.
func (*VersionedDeploymentUnit) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_usage_usage_proto_rawDescGZIP(), []int{5}
}

func (x *VersionedDeploymentUnit) GetUnit() *DeploymentUnit {
	if x != nil {
		return x.Unit
	}
	return nil
}

func (x *VersionedDeploymentUnit) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type NotificationMethod struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Uri  string `protobuf:"bytes,2,opt,name=uri,proto3" json:"uri,omitempty"`
}

func (x *NotificationMethod) Reset() {
	*x = NotificationMethod{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_usage_usage_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotificationMethod) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotificationMethod) ProtoMessage() {}

func (x *NotificationMethod) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_usage_usage_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotificationMethod.ProtoReflect.Descriptor instead.
func (*NotificationMethod) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_usage_usage_proto_rawDescGZIP(), []int{6}
}

func (x *NotificationMethod) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *NotificationMethod) GetUri() string {
	if x != nil {
		return x.Uri
	}
	return ""
}

var File_pb_terrarium_usage_usage_proto protoreflect.FileDescriptor

var file_pb_terrarium_usage_usage_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x62, 0x2f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2f, 0x75,
	0x73, 0x61, 0x67, 0x65, 0x2f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x75, 0x73, 0x61, 0x67,
	0x65, 0x22, 0x9f, 0x01, 0x0a, 0x1d, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x44, 0x65,
	0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x6e, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1f, 0x2e, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x75, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x6e,
	0x69, 0x74, 0x52, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x12, 0x49, 0x0a, 0x0d, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x23, 0x2e, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x75, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x52, 0x0d, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x22, 0x20, 0x0a, 0x1e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x44,
	0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x6e, 0x69, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xa5, 0x01, 0x0a, 0x17, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79,
	0x44, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x6e, 0x63, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x3c, 0x0a, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x28, 0x2e, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x75, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x65, 0x64, 0x44, 0x65, 0x70, 0x6c, 0x6f,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x6e, 0x69, 0x74, 0x52, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x12,
	0x4c, 0x0a, 0x0c, 0x64, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75,
	0x6d, 0x2e, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x65,
	0x64, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x6e, 0x69, 0x74, 0x52,
	0x0c, 0x64, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x22, 0x1a, 0x0a,
	0x18, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x44, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x6e, 0x63,
	0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x5c, 0x0a, 0x0e, 0x44, 0x65, 0x70,
	0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x6e, 0x69, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x22, 0x0a, 0x0c, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x68, 0x0a, 0x17, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x65, 0x64, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x6e,
	0x69, 0x74, 0x12, 0x33, 0x0a, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1f, 0x2e, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x75, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x6e, 0x69,
	0x74, 0x52, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x22, 0x3a, 0x0a, 0x12, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x72, 0x69, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x69, 0x32, 0xf6, 0x01,
	0x0a, 0x11, 0x44, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x6e, 0x63, 0x79, 0x54, 0x72, 0x61, 0x63,
	0x6b, 0x65, 0x72, 0x12, 0x7b, 0x0a, 0x16, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x44,
	0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x6e, 0x69, 0x74, 0x12, 0x2e, 0x2e,
	0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x55, 0x6e, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e,
	0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x55, 0x6e, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x64, 0x0a, 0x0b, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x28, 0x2e, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x75, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x44, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x6e,
	0x63, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x74, 0x65, 0x72, 0x72,
	0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x79, 0x44, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x6e, 0x63, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x2f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2f, 0x75, 0x73, 0x61, 0x67,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_terrarium_usage_usage_proto_rawDescOnce sync.Once
	file_pb_terrarium_usage_usage_proto_rawDescData = file_pb_terrarium_usage_usage_proto_rawDesc
)

func file_pb_terrarium_usage_usage_proto_rawDescGZIP() []byte {
	file_pb_terrarium_usage_usage_proto_rawDescOnce.Do(func() {
		file_pb_terrarium_usage_usage_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_terrarium_usage_usage_proto_rawDescData)
	})
	return file_pb_terrarium_usage_usage_proto_rawDescData
}

var file_pb_terrarium_usage_usage_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_pb_terrarium_usage_usage_proto_goTypes = []interface{}{
	(*RegisterDeploymentUnitRequest)(nil),  // 0: terrarium.usage.RegisterDeploymentUnitRequest
	(*RegisterDeploymentUnitResponse)(nil), // 1: terrarium.usage.RegisterDeploymentUnitResponse
	(*NotifyDependencyRequest)(nil),        // 2: terrarium.usage.NotifyDependencyRequest
	(*NotifyDependencyResponse)(nil),       // 3: terrarium.usage.NotifyDependencyResponse
	(*DeploymentUnit)(nil),                 // 4: terrarium.usage.DeploymentUnit
	(*VersionedDeploymentUnit)(nil),        // 5: terrarium.usage.VersionedDeploymentUnit
	(*NotificationMethod)(nil),             // 6: terrarium.usage.NotificationMethod
}
var file_pb_terrarium_usage_usage_proto_depIdxs = []int32{
	4, // 0: terrarium.usage.RegisterDeploymentUnitRequest.unit:type_name -> terrarium.usage.DeploymentUnit
	6, // 1: terrarium.usage.RegisterDeploymentUnitRequest.notifications:type_name -> terrarium.usage.NotificationMethod
	5, // 2: terrarium.usage.NotifyDependencyRequest.unit:type_name -> terrarium.usage.VersionedDeploymentUnit
	5, // 3: terrarium.usage.NotifyDependencyRequest.dependencies:type_name -> terrarium.usage.VersionedDeploymentUnit
	4, // 4: terrarium.usage.VersionedDeploymentUnit.unit:type_name -> terrarium.usage.DeploymentUnit
	0, // 5: terrarium.usage.DependencyTracker.RegisterDeploymentUnit:input_type -> terrarium.usage.RegisterDeploymentUnitRequest
	2, // 6: terrarium.usage.DependencyTracker.NotifyUsage:input_type -> terrarium.usage.NotifyDependencyRequest
	1, // 7: terrarium.usage.DependencyTracker.RegisterDeploymentUnit:output_type -> terrarium.usage.RegisterDeploymentUnitResponse
	3, // 8: terrarium.usage.DependencyTracker.NotifyUsage:output_type -> terrarium.usage.NotifyDependencyResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_pb_terrarium_usage_usage_proto_init() }
func file_pb_terrarium_usage_usage_proto_init() {
	if File_pb_terrarium_usage_usage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_terrarium_usage_usage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterDeploymentUnitRequest); i {
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
		file_pb_terrarium_usage_usage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterDeploymentUnitResponse); i {
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
		file_pb_terrarium_usage_usage_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyDependencyRequest); i {
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
		file_pb_terrarium_usage_usage_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyDependencyResponse); i {
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
		file_pb_terrarium_usage_usage_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeploymentUnit); i {
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
		file_pb_terrarium_usage_usage_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VersionedDeploymentUnit); i {
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
		file_pb_terrarium_usage_usage_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotificationMethod); i {
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
			RawDescriptor: file_pb_terrarium_usage_usage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_terrarium_usage_usage_proto_goTypes,
		DependencyIndexes: file_pb_terrarium_usage_usage_proto_depIdxs,
		MessageInfos:      file_pb_terrarium_usage_usage_proto_msgTypes,
	}.Build()
	File_pb_terrarium_usage_usage_proto = out.File
	file_pb_terrarium_usage_usage_proto_rawDesc = nil
	file_pb_terrarium_usage_usage_proto_goTypes = nil
	file_pb_terrarium_usage_usage_proto_depIdxs = nil
}
