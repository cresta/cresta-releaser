// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: rpc/releaser/Releaser.proto

package releaser

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PushPromotionResponse_Status int32

const (
	PushPromotionResponse_UNKNOWN               PushPromotionResponse_Status = 0
	PushPromotionResponse_EXISTING_PULL_REQUEST PushPromotionResponse_Status = 1
	PushPromotionResponse_NEW_PULL_REQUEST      PushPromotionResponse_Status = 2
	PushPromotionResponse_NO_CHANGES            PushPromotionResponse_Status = 3
)

// Enum value maps for PushPromotionResponse_Status.
var (
	PushPromotionResponse_Status_name = map[int32]string{
		0: "UNKNOWN",
		1: "EXISTING_PULL_REQUEST",
		2: "NEW_PULL_REQUEST",
		3: "NO_CHANGES",
	}
	PushPromotionResponse_Status_value = map[string]int32{
		"UNKNOWN":               0,
		"EXISTING_PULL_REQUEST": 1,
		"NEW_PULL_REQUEST":      2,
		"NO_CHANGES":            3,
	}
)

func (x PushPromotionResponse_Status) Enum() *PushPromotionResponse_Status {
	p := new(PushPromotionResponse_Status)
	*p = x
	return p
}

func (x PushPromotionResponse_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PushPromotionResponse_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc_releaser_Releaser_proto_enumTypes[0].Descriptor()
}

func (PushPromotionResponse_Status) Type() protoreflect.EnumType {
	return &file_rpc_releaser_Releaser_proto_enumTypes[0]
}

func (x PushPromotionResponse_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PushPromotionResponse_Status.Descriptor instead.
func (PushPromotionResponse_Status) EnumDescriptor() ([]byte, []int) {
	return file_rpc_releaser_Releaser_proto_rawDescGZIP(), []int{3, 0}
}

type ReleaseStatus_Status int32

const (
	ReleaseStatus_UNKNOWN  ReleaseStatus_Status = 0
	ReleaseStatus_PENDING  ReleaseStatus_Status = 1
	ReleaseStatus_RELEASED ReleaseStatus_Status = 2
)

// Enum value maps for ReleaseStatus_Status.
var (
	ReleaseStatus_Status_name = map[int32]string{
		0: "UNKNOWN",
		1: "PENDING",
		2: "RELEASED",
	}
	ReleaseStatus_Status_value = map[string]int32{
		"UNKNOWN":  0,
		"PENDING":  1,
		"RELEASED": 2,
	}
)

func (x ReleaseStatus_Status) Enum() *ReleaseStatus_Status {
	p := new(ReleaseStatus_Status)
	*p = x
	return p
}

func (x ReleaseStatus_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ReleaseStatus_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc_releaser_Releaser_proto_enumTypes[1].Descriptor()
}

func (ReleaseStatus_Status) Type() protoreflect.EnumType {
	return &file_rpc_releaser_Releaser_proto_enumTypes[1]
}

func (x ReleaseStatus_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ReleaseStatus_Status.Descriptor instead.
func (ReleaseStatus_Status) EnumDescriptor() ([]byte, []int) {
	return file_rpc_releaser_Releaser_proto_rawDescGZIP(), []int{7, 0}
}

type RefreshRepositoryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RefreshRepositoryRequest) Reset() {
	*x = RefreshRepositoryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_releaser_Releaser_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefreshRepositoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshRepositoryRequest) ProtoMessage() {}

func (x *RefreshRepositoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_releaser_Releaser_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshRepositoryRequest.ProtoReflect.Descriptor instead.
func (*RefreshRepositoryRequest) Descriptor() ([]byte, []int) {
	return file_rpc_releaser_Releaser_proto_rawDescGZIP(), []int{0}
}

type RefreshRepositoryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RefreshRepositoryResponse) Reset() {
	*x = RefreshRepositoryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_releaser_Releaser_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefreshRepositoryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshRepositoryResponse) ProtoMessage() {}

func (x *RefreshRepositoryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_releaser_Releaser_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshRepositoryResponse.ProtoReflect.Descriptor instead.
func (*RefreshRepositoryResponse) Descriptor() ([]byte, []int) {
	return file_rpc_releaser_Releaser_proto_rawDescGZIP(), []int{1}
}

type PushPromotionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ApplicationName string `protobuf:"bytes,1,opt,name=application_name,json=applicationName,proto3" json:"application_name,omitempty"`
	ReleaseName     string `protobuf:"bytes,2,opt,name=release_name,json=releaseName,proto3" json:"release_name,omitempty"`
}

func (x *PushPromotionRequest) Reset() {
	*x = PushPromotionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_releaser_Releaser_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushPromotionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushPromotionRequest) ProtoMessage() {}

func (x *PushPromotionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_releaser_Releaser_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushPromotionRequest.ProtoReflect.Descriptor instead.
func (*PushPromotionRequest) Descriptor() ([]byte, []int) {
	return file_rpc_releaser_Releaser_proto_rawDescGZIP(), []int{2}
}

func (x *PushPromotionRequest) GetApplicationName() string {
	if x != nil {
		return x.ApplicationName
	}
	return ""
}

func (x *PushPromotionRequest) GetReleaseName() string {
	if x != nil {
		return x.ReleaseName
	}
	return ""
}

type PushPromotionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status        PushPromotionResponse_Status `protobuf:"varint,1,opt,name=status,proto3,enum=cresta.releaser.PushPromotionResponse_Status" json:"status,omitempty"`
	PullRequestId int64                        `protobuf:"varint,2,opt,name=pull_request_id,json=pullRequestId,proto3" json:"pull_request_id,omitempty"`
}

func (x *PushPromotionResponse) Reset() {
	*x = PushPromotionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_releaser_Releaser_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushPromotionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushPromotionResponse) ProtoMessage() {}

func (x *PushPromotionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_releaser_Releaser_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushPromotionResponse.ProtoReflect.Descriptor instead.
func (*PushPromotionResponse) Descriptor() ([]byte, []int) {
	return file_rpc_releaser_Releaser_proto_rawDescGZIP(), []int{3}
}

func (x *PushPromotionResponse) GetStatus() PushPromotionResponse_Status {
	if x != nil {
		return x.Status
	}
	return PushPromotionResponse_UNKNOWN
}

func (x *PushPromotionResponse) GetPullRequestId() int64 {
	if x != nil {
		return x.PullRequestId
	}
	return 0
}

type GetAllApplicationStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetAllApplicationStatusRequest) Reset() {
	*x = GetAllApplicationStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_releaser_Releaser_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllApplicationStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllApplicationStatusRequest) ProtoMessage() {}

func (x *GetAllApplicationStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_releaser_Releaser_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllApplicationStatusRequest.ProtoReflect.Descriptor instead.
func (*GetAllApplicationStatusRequest) Descriptor() ([]byte, []int) {
	return file_rpc_releaser_Releaser_proto_rawDescGZIP(), []int{4}
}

type GetAllApplicationStatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ApplicationStatus []*ApplicationStatus `protobuf:"bytes,1,rep,name=application_status,json=applicationStatus,proto3" json:"application_status,omitempty"`
}

func (x *GetAllApplicationStatusResponse) Reset() {
	*x = GetAllApplicationStatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_releaser_Releaser_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllApplicationStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllApplicationStatusResponse) ProtoMessage() {}

func (x *GetAllApplicationStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_releaser_Releaser_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllApplicationStatusResponse.ProtoReflect.Descriptor instead.
func (*GetAllApplicationStatusResponse) Descriptor() ([]byte, []int) {
	return file_rpc_releaser_Releaser_proto_rawDescGZIP(), []int{5}
}

func (x *GetAllApplicationStatusResponse) GetApplicationStatus() []*ApplicationStatus {
	if x != nil {
		return x.ApplicationStatus
	}
	return nil
}

type ApplicationStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name          string           `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ReleaseStatus []*ReleaseStatus `protobuf:"bytes,2,rep,name=release_status,json=releaseStatus,proto3" json:"release_status,omitempty"`
}

func (x *ApplicationStatus) Reset() {
	*x = ApplicationStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_releaser_Releaser_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplicationStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplicationStatus) ProtoMessage() {}

func (x *ApplicationStatus) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_releaser_Releaser_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplicationStatus.ProtoReflect.Descriptor instead.
func (*ApplicationStatus) Descriptor() ([]byte, []int) {
	return file_rpc_releaser_Releaser_proto_rawDescGZIP(), []int{6}
}

func (x *ApplicationStatus) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ApplicationStatus) GetReleaseStatus() []*ReleaseStatus {
	if x != nil {
		return x.ReleaseStatus
	}
	return nil
}

type ReleaseStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name           string               `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Status         ReleaseStatus_Status `protobuf:"varint,2,opt,name=status,proto3,enum=cresta.releaser.ReleaseStatus_Status" json:"status,omitempty"`
	PrNumber       int64                `protobuf:"varint,3,opt,name=pr_number,json=prNumber,proto3" json:"pr_number,omitempty"`
	OriginalGitSha string               `protobuf:"bytes,4,opt,name=original_git_sha,json=originalGitSha,proto3" json:"original_git_sha,omitempty"`
}

func (x *ReleaseStatus) Reset() {
	*x = ReleaseStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_releaser_Releaser_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReleaseStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReleaseStatus) ProtoMessage() {}

func (x *ReleaseStatus) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_releaser_Releaser_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReleaseStatus.ProtoReflect.Descriptor instead.
func (*ReleaseStatus) Descriptor() ([]byte, []int) {
	return file_rpc_releaser_Releaser_proto_rawDescGZIP(), []int{7}
}

func (x *ReleaseStatus) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ReleaseStatus) GetStatus() ReleaseStatus_Status {
	if x != nil {
		return x.Status
	}
	return ReleaseStatus_UNKNOWN
}

func (x *ReleaseStatus) GetPrNumber() int64 {
	if x != nil {
		return x.PrNumber
	}
	return 0
}

func (x *ReleaseStatus) GetOriginalGitSha() string {
	if x != nil {
		return x.OriginalGitSha
	}
	return ""
}

var File_rpc_releaser_Releaser_proto protoreflect.FileDescriptor

var file_rpc_releaser_Releaser_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x72, 0x70, 0x63, 0x2f, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x72, 0x2f, 0x52,
	0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x63,
	0x72, 0x65, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x72, 0x22, 0x1a,
	0x0a, 0x18, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x1b, 0x0a, 0x19, 0x52, 0x65,
	0x66, 0x72, 0x65, 0x73, 0x68, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x64, 0x0a, 0x14, 0x50, 0x75, 0x73, 0x68, 0x50,
	0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x29, 0x0a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x61, 0x70, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xde, 0x01,
	0x0a, 0x15, 0x50, 0x75, 0x73, 0x68, 0x50, 0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2d, 0x2e, 0x63, 0x72, 0x65, 0x73, 0x74, 0x61,
	0x2e, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x72, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x50, 0x72,
	0x6f, 0x6d, 0x6f, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x26,
	0x0a, 0x0f, 0x70, 0x75, 0x6c, 0x6c, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x70, 0x75, 0x6c, 0x6c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x22, 0x56, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x19, 0x0a,
	0x15, 0x45, 0x58, 0x49, 0x53, 0x54, 0x49, 0x4e, 0x47, 0x5f, 0x50, 0x55, 0x4c, 0x4c, 0x5f, 0x52,
	0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x4e, 0x45, 0x57, 0x5f,
	0x50, 0x55, 0x4c, 0x4c, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x10, 0x02, 0x12, 0x0e,
	0x0a, 0x0a, 0x4e, 0x4f, 0x5f, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45, 0x53, 0x10, 0x03, 0x22, 0x20,
	0x0a, 0x1e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x74, 0x0a, 0x1f, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x51, 0x0a, 0x12, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x22, 0x2e, 0x63, 0x72, 0x65, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65,
	0x72, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x11, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x6e, 0x0a, 0x11, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x45, 0x0a, 0x0e, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x63, 0x72, 0x65, 0x73, 0x74, 0x61,
	0x2e, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0d, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0xdb, 0x01, 0x0a, 0x0d, 0x52, 0x65, 0x6c, 0x65, 0x61,
	0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3d, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x25, 0x2e, 0x63,
	0x72, 0x65, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x72, 0x2e, 0x52,
	0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x70,
	0x72, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x70, 0x72, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x28, 0x0a, 0x10, 0x6f, 0x72, 0x69, 0x67,
	0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x67, 0x69, 0x74, 0x5f, 0x73, 0x68, 0x61, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x47, 0x69, 0x74, 0x53,
	0x68, 0x61, 0x22, 0x30, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07,
	0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x45, 0x4e,
	0x44, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x4c, 0x45, 0x41, 0x53,
	0x45, 0x44, 0x10, 0x02, 0x32, 0xd4, 0x02, 0x0a, 0x08, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65,
	0x72, 0x12, 0x7c, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x41, 0x70, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2f, 0x2e, 0x63,
	0x72, 0x65, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x72, 0x2e, 0x47,
	0x65, 0x74, 0x41, 0x6c, 0x6c, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x30, 0x2e,
	0x63, 0x72, 0x65, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x72, 0x2e,
	0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x5e, 0x0a, 0x0d, 0x50, 0x75, 0x73, 0x68, 0x50, 0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x25, 0x2e, 0x63, 0x72, 0x65, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73,
	0x65, 0x72, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x50, 0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x63, 0x72, 0x65, 0x73, 0x74, 0x61,
	0x2e, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x72, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x50, 0x72,
	0x6f, 0x6d, 0x6f, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x6a, 0x0a, 0x11, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x6f, 0x72, 0x79, 0x12, 0x29, 0x2e, 0x63, 0x72, 0x65, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x52, 0x65,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2a, 0x2e, 0x63, 0x72, 0x65, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65,
	0x72, 0x2e, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x30, 0x5a, 0x2e, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x72, 0x65, 0x73, 0x74, 0x61,
	0x2f, 0x63, 0x72, 0x65, 0x73, 0x74, 0x61, 0x2d, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x72,
	0x2f, 0x72, 0x70, 0x63, 0x2f, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x72, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_releaser_Releaser_proto_rawDescOnce sync.Once
	file_rpc_releaser_Releaser_proto_rawDescData = file_rpc_releaser_Releaser_proto_rawDesc
)

func file_rpc_releaser_Releaser_proto_rawDescGZIP() []byte {
	file_rpc_releaser_Releaser_proto_rawDescOnce.Do(func() {
		file_rpc_releaser_Releaser_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_releaser_Releaser_proto_rawDescData)
	})
	return file_rpc_releaser_Releaser_proto_rawDescData
}

var file_rpc_releaser_Releaser_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_rpc_releaser_Releaser_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_rpc_releaser_Releaser_proto_goTypes = []interface{}{
	(PushPromotionResponse_Status)(0),       // 0: cresta.releaser.PushPromotionResponse.Status
	(ReleaseStatus_Status)(0),               // 1: cresta.releaser.ReleaseStatus.Status
	(*RefreshRepositoryRequest)(nil),        // 2: cresta.releaser.RefreshRepositoryRequest
	(*RefreshRepositoryResponse)(nil),       // 3: cresta.releaser.RefreshRepositoryResponse
	(*PushPromotionRequest)(nil),            // 4: cresta.releaser.PushPromotionRequest
	(*PushPromotionResponse)(nil),           // 5: cresta.releaser.PushPromotionResponse
	(*GetAllApplicationStatusRequest)(nil),  // 6: cresta.releaser.GetAllApplicationStatusRequest
	(*GetAllApplicationStatusResponse)(nil), // 7: cresta.releaser.GetAllApplicationStatusResponse
	(*ApplicationStatus)(nil),               // 8: cresta.releaser.ApplicationStatus
	(*ReleaseStatus)(nil),                   // 9: cresta.releaser.ReleaseStatus
}
var file_rpc_releaser_Releaser_proto_depIdxs = []int32{
	0, // 0: cresta.releaser.PushPromotionResponse.status:type_name -> cresta.releaser.PushPromotionResponse.Status
	8, // 1: cresta.releaser.GetAllApplicationStatusResponse.application_status:type_name -> cresta.releaser.ApplicationStatus
	9, // 2: cresta.releaser.ApplicationStatus.release_status:type_name -> cresta.releaser.ReleaseStatus
	1, // 3: cresta.releaser.ReleaseStatus.status:type_name -> cresta.releaser.ReleaseStatus.Status
	6, // 4: cresta.releaser.Releaser.GetAllApplicationStatus:input_type -> cresta.releaser.GetAllApplicationStatusRequest
	4, // 5: cresta.releaser.Releaser.PushPromotion:input_type -> cresta.releaser.PushPromotionRequest
	2, // 6: cresta.releaser.Releaser.RefreshRepository:input_type -> cresta.releaser.RefreshRepositoryRequest
	7, // 7: cresta.releaser.Releaser.GetAllApplicationStatus:output_type -> cresta.releaser.GetAllApplicationStatusResponse
	5, // 8: cresta.releaser.Releaser.PushPromotion:output_type -> cresta.releaser.PushPromotionResponse
	3, // 9: cresta.releaser.Releaser.RefreshRepository:output_type -> cresta.releaser.RefreshRepositoryResponse
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_rpc_releaser_Releaser_proto_init() }
func file_rpc_releaser_Releaser_proto_init() {
	if File_rpc_releaser_Releaser_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_releaser_Releaser_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefreshRepositoryRequest); i {
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
		file_rpc_releaser_Releaser_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefreshRepositoryResponse); i {
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
		file_rpc_releaser_Releaser_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushPromotionRequest); i {
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
		file_rpc_releaser_Releaser_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushPromotionResponse); i {
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
		file_rpc_releaser_Releaser_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllApplicationStatusRequest); i {
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
		file_rpc_releaser_Releaser_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllApplicationStatusResponse); i {
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
		file_rpc_releaser_Releaser_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplicationStatus); i {
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
		file_rpc_releaser_Releaser_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReleaseStatus); i {
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
			RawDescriptor: file_rpc_releaser_Releaser_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_releaser_Releaser_proto_goTypes,
		DependencyIndexes: file_rpc_releaser_Releaser_proto_depIdxs,
		EnumInfos:         file_rpc_releaser_Releaser_proto_enumTypes,
		MessageInfos:      file_rpc_releaser_Releaser_proto_msgTypes,
	}.Build()
	File_rpc_releaser_Releaser_proto = out.File
	file_rpc_releaser_Releaser_proto_rawDesc = nil
	file_rpc_releaser_Releaser_proto_goTypes = nil
	file_rpc_releaser_Releaser_proto_depIdxs = nil
}
