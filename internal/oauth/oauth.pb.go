// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.22.2
// source: pb/terrarium/oauth/oauth.proto

package oauth

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

type CreateApplicationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ApplicationName        string `protobuf:"bytes,1,opt,name=application_name,json=applicationName,proto3" json:"application_name,omitempty"`
	HomepageUrl            string `protobuf:"bytes,2,opt,name=homepage_url,json=homepageUrl,proto3" json:"homepage_url,omitempty"`
	ApplicationDescription string `protobuf:"bytes,3,opt,name=application_description,json=applicationDescription,proto3" json:"application_description,omitempty"`
	CallbackUrl            string `protobuf:"bytes,4,opt,name=callback_url,json=callbackUrl,proto3" json:"callback_url,omitempty"`
}

func (x *CreateApplicationRequest) Reset() {
	*x = CreateApplicationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateApplicationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateApplicationRequest) ProtoMessage() {}

func (x *CreateApplicationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateApplicationRequest.ProtoReflect.Descriptor instead.
func (*CreateApplicationRequest) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_oauth_oauth_proto_rawDescGZIP(), []int{0}
}

func (x *CreateApplicationRequest) GetApplicationName() string {
	if x != nil {
		return x.ApplicationName
	}
	return ""
}

func (x *CreateApplicationRequest) GetHomepageUrl() string {
	if x != nil {
		return x.HomepageUrl
	}
	return ""
}

func (x *CreateApplicationRequest) GetApplicationDescription() string {
	if x != nil {
		return x.ApplicationDescription
	}
	return ""
}

func (x *CreateApplicationRequest) GetCallbackUrl() string {
	if x != nil {
		return x.CallbackUrl
	}
	return ""
}

type UpdateApplicationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ApplicationName        string `protobuf:"bytes,2,opt,name=application_name,json=applicationName,proto3" json:"application_name,omitempty"`
	HomepageUrl            string `protobuf:"bytes,3,opt,name=homepage_url,json=homepageUrl,proto3" json:"homepage_url,omitempty"`
	ApplicationDescription string `protobuf:"bytes,4,opt,name=application_description,json=applicationDescription,proto3" json:"application_description,omitempty"`
	CallbackUrl            string `protobuf:"bytes,5,opt,name=callback_url,json=callbackUrl,proto3" json:"callback_url,omitempty"`
}

func (x *UpdateApplicationRequest) Reset() {
	*x = UpdateApplicationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateApplicationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateApplicationRequest) ProtoMessage() {}

func (x *UpdateApplicationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateApplicationRequest.ProtoReflect.Descriptor instead.
func (*UpdateApplicationRequest) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_oauth_oauth_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateApplicationRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateApplicationRequest) GetApplicationName() string {
	if x != nil {
		return x.ApplicationName
	}
	return ""
}

func (x *UpdateApplicationRequest) GetHomepageUrl() string {
	if x != nil {
		return x.HomepageUrl
	}
	return ""
}

func (x *UpdateApplicationRequest) GetApplicationDescription() string {
	if x != nil {
		return x.ApplicationDescription
	}
	return ""
}

func (x *UpdateApplicationRequest) GetCallbackUrl() string {
	if x != nil {
		return x.CallbackUrl
	}
	return ""
}

type DeleteApplicationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteApplicationRequest) Reset() {
	*x = DeleteApplicationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteApplicationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteApplicationRequest) ProtoMessage() {}

func (x *DeleteApplicationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteApplicationRequest.ProtoReflect.Descriptor instead.
func (*DeleteApplicationRequest) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_oauth_oauth_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteApplicationRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type RotateApplicationSecretsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RotateApplicationSecretsRequest) Reset() {
	*x = RotateApplicationSecretsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RotateApplicationSecretsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RotateApplicationSecretsRequest) ProtoMessage() {}

func (x *RotateApplicationSecretsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RotateApplicationSecretsRequest.ProtoReflect.Descriptor instead.
func (*RotateApplicationSecretsRequest) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_oauth_oauth_proto_rawDescGZIP(), []int{3}
}

func (x *RotateApplicationSecretsRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type RotateApplicationSecretsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ClientId     string `protobuf:"bytes,2,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	ClientSecret string `protobuf:"bytes,3,opt,name=client_secret,json=clientSecret,proto3" json:"client_secret,omitempty"`
}

func (x *RotateApplicationSecretsResponse) Reset() {
	*x = RotateApplicationSecretsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RotateApplicationSecretsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RotateApplicationSecretsResponse) ProtoMessage() {}

func (x *RotateApplicationSecretsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RotateApplicationSecretsResponse.ProtoReflect.Descriptor instead.
func (*RotateApplicationSecretsResponse) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_oauth_oauth_proto_rawDescGZIP(), []int{4}
}

func (x *RotateApplicationSecretsResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RotateApplicationSecretsResponse) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *RotateApplicationSecretsResponse) GetClientSecret() string {
	if x != nil {
		return x.ClientSecret
	}
	return ""
}

type ApplicationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ApplicationName        string `protobuf:"bytes,2,opt,name=application_name,json=applicationName,proto3" json:"application_name,omitempty"`
	HomepageUrl            string `protobuf:"bytes,3,opt,name=homepage_url,json=homepageUrl,proto3" json:"homepage_url,omitempty"`
	ApplicationDescription string `protobuf:"bytes,4,opt,name=application_description,json=applicationDescription,proto3" json:"application_description,omitempty"`
	CallbackUrl            string `protobuf:"bytes,5,opt,name=callback_url,json=callbackUrl,proto3" json:"callback_url,omitempty"`
}

func (x *ApplicationResponse) Reset() {
	*x = ApplicationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplicationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplicationResponse) ProtoMessage() {}

func (x *ApplicationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplicationResponse.ProtoReflect.Descriptor instead.
func (*ApplicationResponse) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_oauth_oauth_proto_rawDescGZIP(), []int{5}
}

func (x *ApplicationResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ApplicationResponse) GetApplicationName() string {
	if x != nil {
		return x.ApplicationName
	}
	return ""
}

func (x *ApplicationResponse) GetHomepageUrl() string {
	if x != nil {
		return x.HomepageUrl
	}
	return ""
}

func (x *ApplicationResponse) GetApplicationDescription() string {
	if x != nil {
		return x.ApplicationDescription
	}
	return ""
}

func (x *ApplicationResponse) GetCallbackUrl() string {
	if x != nil {
		return x.CallbackUrl
	}
	return ""
}

type IssueJWTTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GrantType    string `protobuf:"bytes,1,opt,name=grant_type,json=grantType,proto3" json:"grant_type,omitempty"`
	ClientId     string `protobuf:"bytes,2,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	ClientSecret string `protobuf:"bytes,3,opt,name=client_secret,json=clientSecret,proto3" json:"client_secret,omitempty"`
	RedirectUri  string `protobuf:"bytes,4,opt,name=redirect_uri,json=redirectUri,proto3" json:"redirect_uri,omitempty"`
	Code         string `protobuf:"bytes,5,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *IssueJWTTokenRequest) Reset() {
	*x = IssueJWTTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IssueJWTTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IssueJWTTokenRequest) ProtoMessage() {}

func (x *IssueJWTTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IssueJWTTokenRequest.ProtoReflect.Descriptor instead.
func (*IssueJWTTokenRequest) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_oauth_oauth_proto_rawDescGZIP(), []int{6}
}

func (x *IssueJWTTokenRequest) GetGrantType() string {
	if x != nil {
		return x.GrantType
	}
	return ""
}

func (x *IssueJWTTokenRequest) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *IssueJWTTokenRequest) GetClientSecret() string {
	if x != nil {
		return x.ClientSecret
	}
	return ""
}

func (x *IssueJWTTokenRequest) GetRedirectUri() string {
	if x != nil {
		return x.RedirectUri
	}
	return ""
}

func (x *IssueJWTTokenRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type IssueJWTTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	TokenType   string `protobuf:"bytes,2,opt,name=token_type,json=tokenType,proto3" json:"token_type,omitempty"`
	Scope       string `protobuf:"bytes,3,opt,name=scope,proto3" json:"scope,omitempty"`
}

func (x *IssueJWTTokenResponse) Reset() {
	*x = IssueJWTTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IssueJWTTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IssueJWTTokenResponse) ProtoMessage() {}

func (x *IssueJWTTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_oauth_oauth_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IssueJWTTokenResponse.ProtoReflect.Descriptor instead.
func (*IssueJWTTokenResponse) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_oauth_oauth_proto_rawDescGZIP(), []int{7}
}

func (x *IssueJWTTokenResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *IssueJWTTokenResponse) GetTokenType() string {
	if x != nil {
		return x.TokenType
	}
	return ""
}

func (x *IssueJWTTokenResponse) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

var File_pb_terrarium_oauth_oauth_proto protoreflect.FileDescriptor

var file_pb_terrarium_oauth_oauth_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x62, 0x2f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2f, 0x6f,
	0x61, 0x75, 0x74, 0x68, 0x2f, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x6f, 0x61, 0x75, 0x74,
	0x68, 0x22, 0xc4, 0x01, 0x0a, 0x18, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29,
	0x0a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x68, 0x6f, 0x6d,
	0x65, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x68, 0x6f, 0x6d, 0x65, 0x70, 0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x37, 0x0a, 0x17,
	0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x16, 0x61,
	0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63,
	0x6b, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x61, 0x6c,
	0x6c, 0x62, 0x61, 0x63, 0x6b, 0x55, 0x72, 0x6c, 0x22, 0xd4, 0x01, 0x0a, 0x18, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x29, 0x0a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x21, 0x0a, 0x0c, 0x68, 0x6f, 0x6d, 0x65, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x68, 0x6f, 0x6d, 0x65, 0x70, 0x61, 0x67, 0x65,
	0x55, 0x72, 0x6c, 0x12, 0x37, 0x0a, 0x17, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x16, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c,
	0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x55, 0x72, 0x6c, 0x22,
	0x2a, 0x0a, 0x18, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x31, 0x0a, 0x1f, 0x52,
	0x6f, 0x74, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x74,
	0x0a, 0x20, 0x52, 0x6f, 0x74, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x23, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65,
	0x63, 0x72, 0x65, 0x74, 0x22, 0xcf, 0x01, 0x0a, 0x13, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x29, 0x0a, 0x10,
	0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x68, 0x6f, 0x6d, 0x65, 0x70,
	0x61, 0x67, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x68,
	0x6f, 0x6d, 0x65, 0x70, 0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x37, 0x0a, 0x17, 0x61, 0x70,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x16, 0x61, 0x70, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x5f,
	0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x61, 0x6c, 0x6c, 0x62,
	0x61, 0x63, 0x6b, 0x55, 0x72, 0x6c, 0x22, 0xae, 0x01, 0x0a, 0x14, 0x49, 0x73, 0x73, 0x75, 0x65,
	0x4a, 0x57, 0x54, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1d, 0x0a, 0x0a, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x75, 0x72, 0x69,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x55, 0x72, 0x69, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x6f, 0x0a, 0x15, 0x49, 0x73, 0x73, 0x75, 0x65,
	0x4a, 0x57, 0x54, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_terrarium_oauth_oauth_proto_rawDescOnce sync.Once
	file_pb_terrarium_oauth_oauth_proto_rawDescData = file_pb_terrarium_oauth_oauth_proto_rawDesc
)

func file_pb_terrarium_oauth_oauth_proto_rawDescGZIP() []byte {
	file_pb_terrarium_oauth_oauth_proto_rawDescOnce.Do(func() {
		file_pb_terrarium_oauth_oauth_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_terrarium_oauth_oauth_proto_rawDescData)
	})
	return file_pb_terrarium_oauth_oauth_proto_rawDescData
}

var file_pb_terrarium_oauth_oauth_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_pb_terrarium_oauth_oauth_proto_goTypes = []interface{}{
	(*CreateApplicationRequest)(nil),         // 0: terrarium.oauth.CreateApplicationRequest
	(*UpdateApplicationRequest)(nil),         // 1: terrarium.oauth.UpdateApplicationRequest
	(*DeleteApplicationRequest)(nil),         // 2: terrarium.oauth.DeleteApplicationRequest
	(*RotateApplicationSecretsRequest)(nil),  // 3: terrarium.oauth.RotateApplicationSecretsRequest
	(*RotateApplicationSecretsResponse)(nil), // 4: terrarium.oauth.RotateApplicationSecretsResponse
	(*ApplicationResponse)(nil),              // 5: terrarium.oauth.ApplicationResponse
	(*IssueJWTTokenRequest)(nil),             // 6: terrarium.oauth.IssueJWTTokenRequest
	(*IssueJWTTokenResponse)(nil),            // 7: terrarium.oauth.IssueJWTTokenResponse
}
var file_pb_terrarium_oauth_oauth_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pb_terrarium_oauth_oauth_proto_init() }
func file_pb_terrarium_oauth_oauth_proto_init() {
	if File_pb_terrarium_oauth_oauth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_terrarium_oauth_oauth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateApplicationRequest); i {
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
		file_pb_terrarium_oauth_oauth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateApplicationRequest); i {
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
		file_pb_terrarium_oauth_oauth_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteApplicationRequest); i {
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
		file_pb_terrarium_oauth_oauth_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RotateApplicationSecretsRequest); i {
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
		file_pb_terrarium_oauth_oauth_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RotateApplicationSecretsResponse); i {
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
		file_pb_terrarium_oauth_oauth_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplicationResponse); i {
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
		file_pb_terrarium_oauth_oauth_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IssueJWTTokenRequest); i {
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
		file_pb_terrarium_oauth_oauth_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IssueJWTTokenResponse); i {
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
			RawDescriptor: file_pb_terrarium_oauth_oauth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pb_terrarium_oauth_oauth_proto_goTypes,
		DependencyIndexes: file_pb_terrarium_oauth_oauth_proto_depIdxs,
		MessageInfos:      file_pb_terrarium_oauth_oauth_proto_msgTypes,
	}.Build()
	File_pb_terrarium_oauth_oauth_proto = out.File
	file_pb_terrarium_oauth_oauth_proto_rawDesc = nil
	file_pb_terrarium_oauth_oauth_proto_goTypes = nil
	file_pb_terrarium_oauth_oauth_proto_depIdxs = nil
}
