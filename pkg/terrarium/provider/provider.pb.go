// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: pb/terrarium/provider/provider.proto

package provider

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

type Maturity int32

const (
	Maturity_IDEA        Maturity = 0
	Maturity_PLANNING    Maturity = 1
	Maturity_DEVELOPING  Maturity = 2
	Maturity_ALPHA       Maturity = 3
	Maturity_BETA        Maturity = 4
	Maturity_STABLE      Maturity = 5
	Maturity_DEPRECATED  Maturity = 6
	Maturity_END_OF_LIFE Maturity = 7
)

// Enum value maps for Maturity.
var (
	Maturity_name = map[int32]string{
		0: "IDEA",
		1: "PLANNING",
		2: "DEVELOPING",
		3: "ALPHA",
		4: "BETA",
		5: "STABLE",
		6: "DEPRECATED",
		7: "END_OF_LIFE",
	}
	Maturity_value = map[string]int32{
		"IDEA":        0,
		"PLANNING":    1,
		"DEVELOPING":  2,
		"ALPHA":       3,
		"BETA":        4,
		"STABLE":      5,
		"DEPRECATED":  6,
		"END_OF_LIFE": 7,
	}
)

func (x Maturity) Enum() *Maturity {
	p := new(Maturity)
	*p = x
	return p
}

func (x Maturity) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Maturity) Descriptor() protoreflect.EnumDescriptor {
	return file_pb_terrarium_provider_provider_proto_enumTypes[0].Descriptor()
}

func (Maturity) Type() protoreflect.EnumType {
	return &file_pb_terrarium_provider_provider_proto_enumTypes[0]
}

func (x Maturity) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Maturity.Descriptor instead.
func (Maturity) EnumDescriptor() ([]byte, []int) {
	return file_pb_terrarium_provider_provider_proto_rawDescGZIP(), []int{0}
}

type EndProviderRequest_Action int32

const (
	EndProviderRequest_DISCARD_VERSION EndProviderRequest_Action = 0
	EndProviderRequest_PUBLISH         EndProviderRequest_Action = 1
)

// Enum value maps for EndProviderRequest_Action.
var (
	EndProviderRequest_Action_name = map[int32]string{
		0: "DISCARD_VERSION",
		1: "PUBLISH",
	}
	EndProviderRequest_Action_value = map[string]int32{
		"DISCARD_VERSION": 0,
		"PUBLISH":         1,
	}
)

func (x EndProviderRequest_Action) Enum() *EndProviderRequest_Action {
	p := new(EndProviderRequest_Action)
	*p = x
	return p
}

func (x EndProviderRequest_Action) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EndProviderRequest_Action) Descriptor() protoreflect.EnumDescriptor {
	return file_pb_terrarium_provider_provider_proto_enumTypes[1].Descriptor()
}

func (EndProviderRequest_Action) Type() protoreflect.EnumType {
	return &file_pb_terrarium_provider_provider_proto_enumTypes[1]
}

func (x EndProviderRequest_Action) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EndProviderRequest_Action.Descriptor instead.
func (EndProviderRequest_Action) EnumDescriptor() ([]byte, []int) {
	return file_pb_terrarium_provider_provider_proto_rawDescGZIP(), []int{6, 0}
}

type RegisterProviderRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ApiKey        string          `protobuf:"bytes,1,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
	Name          string          `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Version       string          `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	Protocols     []string        `protobuf:"bytes,4,rep,name=protocols,proto3" json:"protocols,omitempty"`
	Platforms     []*PlatformItem `protobuf:"bytes,5,rep,name=platforms,proto3" json:"platforms,omitempty"`
	Description   string          `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	SourceRepoUrl string          `protobuf:"bytes,7,opt,name=source_repo_url,json=sourceRepoUrl,proto3" json:"source_repo_url,omitempty"`
	Maturity      Maturity        `protobuf:"varint,8,opt,name=maturity,proto3,enum=terrarium.provider.Maturity" json:"maturity,omitempty"`
	CreatedOn     string          `protobuf:"bytes,9,opt,name=created_on,json=createdOn,proto3" json:"created_on,omitempty"`
	ModifiedOn    string          `protobuf:"bytes,10,opt,name=modified_on,json=modifiedOn,proto3" json:"modified_on,omitempty"`
	PublishedOn   string          `protobuf:"bytes,11,opt,name=published_on,json=publishedOn,proto3" json:"published_on,omitempty"`
}

func (x *RegisterProviderRequest) Reset() {
	*x = RegisterProviderRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_provider_provider_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterProviderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterProviderRequest) ProtoMessage() {}

func (x *RegisterProviderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_provider_provider_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterProviderRequest.ProtoReflect.Descriptor instead.
func (*RegisterProviderRequest) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_provider_provider_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterProviderRequest) GetApiKey() string {
	if x != nil {
		return x.ApiKey
	}
	return ""
}

func (x *RegisterProviderRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RegisterProviderRequest) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *RegisterProviderRequest) GetProtocols() []string {
	if x != nil {
		return x.Protocols
	}
	return nil
}

func (x *RegisterProviderRequest) GetPlatforms() []*PlatformItem {
	if x != nil {
		return x.Platforms
	}
	return nil
}

func (x *RegisterProviderRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *RegisterProviderRequest) GetSourceRepoUrl() string {
	if x != nil {
		return x.SourceRepoUrl
	}
	return ""
}

func (x *RegisterProviderRequest) GetMaturity() Maturity {
	if x != nil {
		return x.Maturity
	}
	return Maturity_IDEA
}

func (x *RegisterProviderRequest) GetCreatedOn() string {
	if x != nil {
		return x.CreatedOn
	}
	return ""
}

func (x *RegisterProviderRequest) GetModifiedOn() string {
	if x != nil {
		return x.ModifiedOn
	}
	return ""
}

func (x *RegisterProviderRequest) GetPublishedOn() string {
	if x != nil {
		return x.PublishedOn
	}
	return ""
}

type PlatformItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Os                  string       `protobuf:"bytes,1,opt,name=os,proto3" json:"os,omitempty"`
	Arch                string       `protobuf:"bytes,2,opt,name=arch,proto3" json:"arch,omitempty"`
	Filename            string       `protobuf:"bytes,3,opt,name=filename,proto3" json:"filename,omitempty"`
	DownloadUrl         string       `protobuf:"bytes,4,opt,name=download_url,json=downloadUrl,proto3" json:"download_url,omitempty"`
	ShasumsUrl          string       `protobuf:"bytes,5,opt,name=shasums_url,json=shasumsUrl,proto3" json:"shasums_url,omitempty"`
	ShasumsSignatureUrl string       `protobuf:"bytes,6,opt,name=shasums_signature_url,json=shasumsSignatureUrl,proto3" json:"shasums_signature_url,omitempty"`
	Shasum              string       `protobuf:"bytes,7,opt,name=shasum,proto3" json:"shasum,omitempty"`
	SigningKeys         *SigningKeys `protobuf:"bytes,8,opt,name=signing_keys,json=signingKeys,proto3" json:"signing_keys,omitempty"`
}

func (x *PlatformItem) Reset() {
	*x = PlatformItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_provider_provider_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlatformItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlatformItem) ProtoMessage() {}

func (x *PlatformItem) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_provider_provider_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlatformItem.ProtoReflect.Descriptor instead.
func (*PlatformItem) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_provider_provider_proto_rawDescGZIP(), []int{1}
}

func (x *PlatformItem) GetOs() string {
	if x != nil {
		return x.Os
	}
	return ""
}

func (x *PlatformItem) GetArch() string {
	if x != nil {
		return x.Arch
	}
	return ""
}

func (x *PlatformItem) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *PlatformItem) GetDownloadUrl() string {
	if x != nil {
		return x.DownloadUrl
	}
	return ""
}

func (x *PlatformItem) GetShasumsUrl() string {
	if x != nil {
		return x.ShasumsUrl
	}
	return ""
}

func (x *PlatformItem) GetShasumsSignatureUrl() string {
	if x != nil {
		return x.ShasumsSignatureUrl
	}
	return ""
}

func (x *PlatformItem) GetShasum() string {
	if x != nil {
		return x.Shasum
	}
	return ""
}

func (x *PlatformItem) GetSigningKeys() *SigningKeys {
	if x != nil {
		return x.SigningKeys
	}
	return nil
}

type SigningKeys struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GpgPublicKeys []*GPGPublicKey `protobuf:"bytes,1,rep,name=gpg_public_keys,json=gpgPublicKeys,proto3" json:"gpg_public_keys,omitempty"`
}

func (x *SigningKeys) Reset() {
	*x = SigningKeys{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_provider_provider_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SigningKeys) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SigningKeys) ProtoMessage() {}

func (x *SigningKeys) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_provider_provider_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SigningKeys.ProtoReflect.Descriptor instead.
func (*SigningKeys) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_provider_provider_proto_rawDescGZIP(), []int{2}
}

func (x *SigningKeys) GetGpgPublicKeys() []*GPGPublicKey {
	if x != nil {
		return x.GpgPublicKeys
	}
	return nil
}

type GPGPublicKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	KeyId          string `protobuf:"bytes,1,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
	AsciiArmor     string `protobuf:"bytes,2,opt,name=ascii_armor,json=asciiArmor,proto3" json:"ascii_armor,omitempty"`
	TrustSignature string `protobuf:"bytes,3,opt,name=trust_signature,json=trustSignature,proto3" json:"trust_signature,omitempty"`
	Source         string `protobuf:"bytes,4,opt,name=source,proto3" json:"source,omitempty"`
	SourceUrl      string `protobuf:"bytes,5,opt,name=source_url,json=sourceUrl,proto3" json:"source_url,omitempty"`
}

func (x *GPGPublicKey) Reset() {
	*x = GPGPublicKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_provider_provider_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GPGPublicKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GPGPublicKey) ProtoMessage() {}

func (x *GPGPublicKey) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_provider_provider_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GPGPublicKey.ProtoReflect.Descriptor instead.
func (*GPGPublicKey) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_provider_provider_proto_rawDescGZIP(), []int{3}
}

func (x *GPGPublicKey) GetKeyId() string {
	if x != nil {
		return x.KeyId
	}
	return ""
}

func (x *GPGPublicKey) GetAsciiArmor() string {
	if x != nil {
		return x.AsciiArmor
	}
	return ""
}

func (x *GPGPublicKey) GetTrustSignature() string {
	if x != nil {
		return x.TrustSignature
	}
	return ""
}

func (x *GPGPublicKey) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *GPGPublicKey) GetSourceUrl() string {
	if x != nil {
		return x.SourceUrl
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_provider_provider_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_provider_provider_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_provider_provider_proto_rawDescGZIP(), []int{4}
}

func (x *Response) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type Provider struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *Provider) Reset() {
	*x = Provider{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_provider_provider_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Provider) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Provider) ProtoMessage() {}

func (x *Provider) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_provider_provider_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Provider.ProtoReflect.Descriptor instead.
func (*Provider) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_provider_provider_proto_rawDescGZIP(), []int{5}
}

func (x *Provider) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Provider) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type EndProviderRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Provider *Provider                 `protobuf:"bytes,1,opt,name=provider,proto3" json:"provider,omitempty"`
	Action   EndProviderRequest_Action `protobuf:"varint,4,opt,name=action,proto3,enum=terrarium.provider.EndProviderRequest_Action" json:"action,omitempty"`
}

func (x *EndProviderRequest) Reset() {
	*x = EndProviderRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_terrarium_provider_provider_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EndProviderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EndProviderRequest) ProtoMessage() {}

func (x *EndProviderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_terrarium_provider_provider_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EndProviderRequest.ProtoReflect.Descriptor instead.
func (*EndProviderRequest) Descriptor() ([]byte, []int) {
	return file_pb_terrarium_provider_provider_proto_rawDescGZIP(), []int{6}
}

func (x *EndProviderRequest) GetProvider() *Provider {
	if x != nil {
		return x.Provider
	}
	return nil
}

func (x *EndProviderRequest) GetAction() EndProviderRequest_Action {
	if x != nil {
		return x.Action
	}
	return EndProviderRequest_DISCARD_VERSION
}

var File_pb_terrarium_provider_provider_proto protoreflect.FileDescriptor

var file_pb_terrarium_provider_provider_proto_rawDesc = []byte{
	0x0a, 0x24, 0x70, 0x62, 0x2f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2f, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75,
	0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x22, 0xa5, 0x03, 0x0a, 0x17, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x61, 0x70, 0x69, 0x5f, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a,
	0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x73, 0x12, 0x3e, 0x0a, 0x09, 0x70,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x2e, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x49, 0x74, 0x65, 0x6d,
	0x52, 0x09, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a,
	0x0f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65,
	0x70, 0x6f, 0x55, 0x72, 0x6c, 0x12, 0x38, 0x0a, 0x08, 0x6d, 0x61, 0x74, 0x75, 0x72, 0x69, 0x74,
	0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72,
	0x69, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x4d, 0x61, 0x74,
	0x75, 0x72, 0x69, 0x74, 0x79, 0x52, 0x08, 0x6d, 0x61, 0x74, 0x75, 0x72, 0x69, 0x74, 0x79, 0x12,
	0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x6f, 0x6e, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x4f, 0x6e, 0x12, 0x1f,
	0x0a, 0x0b, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x5f, 0x6f, 0x6e, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x4f, 0x6e, 0x12,
	0x21, 0x0a, 0x0c, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x5f, 0x6f, 0x6e, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64,
	0x4f, 0x6e, 0x22, 0xa2, 0x02, 0x0a, 0x0c, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x49,
	0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x6f, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x72, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x61, 0x72, 0x63, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x5f,
	0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x6f, 0x77, 0x6e, 0x6c,
	0x6f, 0x61, 0x64, 0x55, 0x72, 0x6c, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x68, 0x61, 0x73, 0x75, 0x6d,
	0x73, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x68, 0x61,
	0x73, 0x75, 0x6d, 0x73, 0x55, 0x72, 0x6c, 0x12, 0x32, 0x0a, 0x15, 0x73, 0x68, 0x61, 0x73, 0x75,
	0x6d, 0x73, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x73, 0x68, 0x61, 0x73, 0x75, 0x6d, 0x73, 0x53,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x68, 0x61, 0x73, 0x75, 0x6d, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x68, 0x61,
	0x73, 0x75, 0x6d, 0x12, 0x42, 0x0a, 0x0c, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x6b,
	0x65, 0x79, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x74, 0x65, 0x72, 0x72,
	0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x53,
	0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x73, 0x52, 0x0b, 0x73, 0x69, 0x67, 0x6e,
	0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x73, 0x22, 0x57, 0x0a, 0x0b, 0x53, 0x69, 0x67, 0x6e, 0x69,
	0x6e, 0x67, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x48, 0x0a, 0x0f, 0x67, 0x70, 0x67, 0x5f, 0x70, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x76,
	0x69, 0x64, 0x65, 0x72, 0x2e, 0x47, 0x50, 0x47, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65,
	0x79, 0x52, 0x0d, 0x67, 0x70, 0x67, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x73,
	0x22, 0xa6, 0x01, 0x0a, 0x0c, 0x47, 0x50, 0x47, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65,
	0x79, 0x12, 0x15, 0x0a, 0x06, 0x6b, 0x65, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x6b, 0x65, 0x79, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x73, 0x63, 0x69,
	0x69, 0x5f, 0x61, 0x72, 0x6d, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61,
	0x73, 0x63, 0x69, 0x69, 0x41, 0x72, 0x6d, 0x6f, 0x72, 0x12, 0x27, 0x0a, 0x0f, 0x74, 0x72, 0x75,
	0x73, 0x74, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x74, 0x72, 0x75, 0x73, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x55, 0x72, 0x6c, 0x22, 0x24, 0x0a, 0x08, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x38, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0xc1, 0x01, 0x0a, 0x12, 0x45, 0x6e,
	0x64, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x38, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x45, 0x0a, 0x06, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2d, 0x2e, 0x74, 0x65, 0x72,
	0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e,
	0x45, 0x6e, 0x64, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x2a, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x13, 0x0a, 0x0f, 0x44,
	0x49, 0x53, 0x43, 0x41, 0x52, 0x44, 0x5f, 0x56, 0x45, 0x52, 0x53, 0x49, 0x4f, 0x4e, 0x10, 0x00,
	0x12, 0x0b, 0x0a, 0x07, 0x50, 0x55, 0x42, 0x4c, 0x49, 0x53, 0x48, 0x10, 0x01, 0x2a, 0x74, 0x0a,
	0x08, 0x4d, 0x61, 0x74, 0x75, 0x72, 0x69, 0x74, 0x79, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x44, 0x45,
	0x41, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x4c, 0x41, 0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x10,
	0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x44, 0x45, 0x56, 0x45, 0x4c, 0x4f, 0x50, 0x49, 0x4e, 0x47, 0x10,
	0x02, 0x12, 0x09, 0x0a, 0x05, 0x41, 0x4c, 0x50, 0x48, 0x41, 0x10, 0x03, 0x12, 0x08, 0x0a, 0x04,
	0x42, 0x45, 0x54, 0x41, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x54, 0x41, 0x42, 0x4c, 0x45,
	0x10, 0x05, 0x12, 0x0e, 0x0a, 0x0a, 0x44, 0x45, 0x50, 0x52, 0x45, 0x43, 0x41, 0x54, 0x45, 0x44,
	0x10, 0x06, 0x12, 0x0f, 0x0a, 0x0b, 0x45, 0x4e, 0x44, 0x5f, 0x4f, 0x46, 0x5f, 0x4c, 0x49, 0x46,
	0x45, 0x10, 0x07, 0x32, 0xcb, 0x01, 0x0a, 0x11, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x12, 0x5f, 0x0a, 0x10, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x2b, 0x2e,
	0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x72, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x74, 0x65, 0x72,
	0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x0b, 0x45, 0x6e,
	0x64, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x26, 0x2e, 0x74, 0x65, 0x72, 0x72,
	0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x45,
	0x6e, 0x64, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1c, 0x2e, 0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x74, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x74,
	0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x65, 0x72,
	0x72, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_terrarium_provider_provider_proto_rawDescOnce sync.Once
	file_pb_terrarium_provider_provider_proto_rawDescData = file_pb_terrarium_provider_provider_proto_rawDesc
)

func file_pb_terrarium_provider_provider_proto_rawDescGZIP() []byte {
	file_pb_terrarium_provider_provider_proto_rawDescOnce.Do(func() {
		file_pb_terrarium_provider_provider_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_terrarium_provider_provider_proto_rawDescData)
	})
	return file_pb_terrarium_provider_provider_proto_rawDescData
}

var file_pb_terrarium_provider_provider_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_pb_terrarium_provider_provider_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_pb_terrarium_provider_provider_proto_goTypes = []interface{}{
	(Maturity)(0),                   // 0: terrarium.provider.Maturity
	(EndProviderRequest_Action)(0),  // 1: terrarium.provider.EndProviderRequest.Action
	(*RegisterProviderRequest)(nil), // 2: terrarium.provider.RegisterProviderRequest
	(*PlatformItem)(nil),            // 3: terrarium.provider.PlatformItem
	(*SigningKeys)(nil),             // 4: terrarium.provider.SigningKeys
	(*GPGPublicKey)(nil),            // 5: terrarium.provider.GPGPublicKey
	(*Response)(nil),                // 6: terrarium.provider.Response
	(*Provider)(nil),                // 7: terrarium.provider.Provider
	(*EndProviderRequest)(nil),      // 8: terrarium.provider.EndProviderRequest
}
var file_pb_terrarium_provider_provider_proto_depIdxs = []int32{
	3, // 0: terrarium.provider.RegisterProviderRequest.platforms:type_name -> terrarium.provider.PlatformItem
	0, // 1: terrarium.provider.RegisterProviderRequest.maturity:type_name -> terrarium.provider.Maturity
	4, // 2: terrarium.provider.PlatformItem.signing_keys:type_name -> terrarium.provider.SigningKeys
	5, // 3: terrarium.provider.SigningKeys.gpg_public_keys:type_name -> terrarium.provider.GPGPublicKey
	7, // 4: terrarium.provider.EndProviderRequest.provider:type_name -> terrarium.provider.Provider
	1, // 5: terrarium.provider.EndProviderRequest.action:type_name -> terrarium.provider.EndProviderRequest.Action
	2, // 6: terrarium.provider.ProviderPublisher.RegisterProvider:input_type -> terrarium.provider.RegisterProviderRequest
	8, // 7: terrarium.provider.ProviderPublisher.EndProvider:input_type -> terrarium.provider.EndProviderRequest
	6, // 8: terrarium.provider.ProviderPublisher.RegisterProvider:output_type -> terrarium.provider.Response
	6, // 9: terrarium.provider.ProviderPublisher.EndProvider:output_type -> terrarium.provider.Response
	8, // [8:10] is the sub-list for method output_type
	6, // [6:8] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_pb_terrarium_provider_provider_proto_init() }
func file_pb_terrarium_provider_provider_proto_init() {
	if File_pb_terrarium_provider_provider_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_terrarium_provider_provider_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterProviderRequest); i {
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
		file_pb_terrarium_provider_provider_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlatformItem); i {
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
		file_pb_terrarium_provider_provider_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SigningKeys); i {
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
		file_pb_terrarium_provider_provider_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GPGPublicKey); i {
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
		file_pb_terrarium_provider_provider_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_pb_terrarium_provider_provider_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Provider); i {
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
		file_pb_terrarium_provider_provider_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EndProviderRequest); i {
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
			RawDescriptor: file_pb_terrarium_provider_provider_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_terrarium_provider_provider_proto_goTypes,
		DependencyIndexes: file_pb_terrarium_provider_provider_proto_depIdxs,
		EnumInfos:         file_pb_terrarium_provider_provider_proto_enumTypes,
		MessageInfos:      file_pb_terrarium_provider_provider_proto_msgTypes,
	}.Build()
	File_pb_terrarium_provider_provider_proto = out.File
	file_pb_terrarium_provider_provider_proto_rawDesc = nil
	file_pb_terrarium_provider_provider_proto_goTypes = nil
	file_pb_terrarium_provider_provider_proto_depIdxs = nil
}