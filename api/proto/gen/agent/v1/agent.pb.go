// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: agent/v1/agent.proto

package agentv1

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

// Backend data relevant to the agent
type StreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SshConfig                []byte                      `protobuf:"bytes,1,opt,name=ssh_config,json=sshConfig,proto3,oneof" json:"ssh_config,omitempty"`
	UserCertificatePublicKey []byte                      `protobuf:"bytes,2,opt,name=user_certificate_public_key,json=userCertificatePublicKey,proto3,oneof" json:"user_certificate_public_key,omitempty"`
	HostCertificatePublicKey []byte                      `protobuf:"bytes,3,opt,name=host_certificate_public_key,json=hostCertificatePublicKey,proto3,oneof" json:"host_certificate_public_key,omitempty"`
	Restore                  *bool                       `protobuf:"varint,4,opt,name=restore,proto3,oneof" json:"restore,omitempty"`
	Principals               []*StreamResponse_Principal `protobuf:"bytes,5,rep,name=principals,proto3" json:"principals,omitempty"`
}

func (x *StreamResponse) Reset() {
	*x = StreamResponse{}
	mi := &file_agent_v1_agent_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamResponse) ProtoMessage() {}

func (x *StreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_agent_v1_agent_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamResponse.ProtoReflect.Descriptor instead.
func (*StreamResponse) Descriptor() ([]byte, []int) {
	return file_agent_v1_agent_proto_rawDescGZIP(), []int{0}
}

func (x *StreamResponse) GetSshConfig() []byte {
	if x != nil {
		return x.SshConfig
	}
	return nil
}

func (x *StreamResponse) GetUserCertificatePublicKey() []byte {
	if x != nil {
		return x.UserCertificatePublicKey
	}
	return nil
}

func (x *StreamResponse) GetHostCertificatePublicKey() []byte {
	if x != nil {
		return x.HostCertificatePublicKey
	}
	return nil
}

func (x *StreamResponse) GetRestore() bool {
	if x != nil && x.Restore != nil {
		return *x.Restore
	}
	return false
}

func (x *StreamResponse) GetPrincipals() []*StreamResponse_Principal {
	if x != nil {
		return x.Principals
	}
	return nil
}

// Information about the agent
type StreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version       *string `protobuf:"bytes,1,opt,name=version,proto3,oneof" json:"version,omitempty"`
	PublicHostKey *string `protobuf:"bytes,2,opt,name=public_host_key,json=publicHostKey,proto3,oneof" json:"public_host_key,omitempty"`
}

func (x *StreamRequest) Reset() {
	*x = StreamRequest{}
	mi := &file_agent_v1_agent_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamRequest) ProtoMessage() {}

func (x *StreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_v1_agent_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamRequest.ProtoReflect.Descriptor instead.
func (*StreamRequest) Descriptor() ([]byte, []int) {
	return file_agent_v1_agent_proto_rawDescGZIP(), []int{1}
}

func (x *StreamRequest) GetVersion() string {
	if x != nil && x.Version != nil {
		return *x.Version
	}
	return ""
}

func (x *StreamRequest) GetPublicHostKey() string {
	if x != nil && x.PublicHostKey != nil {
		return *x.PublicHostKey
	}
	return ""
}

type StreamResponse_Principal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key    string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Values []string `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *StreamResponse_Principal) Reset() {
	*x = StreamResponse_Principal{}
	mi := &file_agent_v1_agent_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamResponse_Principal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamResponse_Principal) ProtoMessage() {}

func (x *StreamResponse_Principal) ProtoReflect() protoreflect.Message {
	mi := &file_agent_v1_agent_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamResponse_Principal.ProtoReflect.Descriptor instead.
func (*StreamResponse_Principal) Descriptor() ([]byte, []int) {
	return file_agent_v1_agent_proto_rawDescGZIP(), []int{0, 0}
}

func (x *StreamResponse_Principal) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *StreamResponse_Principal) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

var File_agent_v1_agent_proto protoreflect.FileDescriptor

var file_agent_v1_agent_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31,
	0x22, 0xb1, 0x03, 0x0a, 0x0e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x0a, 0x73, 0x73, 0x68, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x09, 0x73, 0x73, 0x68, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x1b, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x01, 0x52, 0x18,
	0x75, 0x73, 0x65, 0x72, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x50,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x1b, 0x68,
	0x6f, 0x73, 0x74, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f,
	0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x48, 0x02, 0x52, 0x18, 0x68, 0x6f, 0x73, 0x74, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x88, 0x01, 0x01, 0x12,
	0x1d, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08,
	0x48, 0x03, 0x52, 0x07, 0x72, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x88, 0x01, 0x01, 0x12, 0x42,
	0x0a, 0x0a, 0x70, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c, 0x73, 0x18, 0x05, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x22, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x50, 0x72, 0x69,
	0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c, 0x52, 0x0a, 0x70, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70, 0x61,
	0x6c, 0x73, 0x1a, 0x35, 0x0a, 0x09, 0x50, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x73, 0x73,
	0x68, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x1e, 0x0a, 0x1c, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x70, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x42, 0x1e, 0x0a, 0x1c, 0x5f, 0x68, 0x6f, 0x73,
	0x74, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x70, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x72, 0x65, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x22, 0x7b, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x88, 0x01, 0x01, 0x12, 0x2b, 0x0a, 0x0f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x68,
	0x6f, 0x73, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52,
	0x0d, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x48, 0x6f, 0x73, 0x74, 0x4b, 0x65, 0x79, 0x88, 0x01,
	0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x12, 0x0a,
	0x10, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x6b, 0x65,
	0x79, 0x32, 0x51, 0x0a, 0x0c, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x41, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x17, 0x2e, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x28, 0x01, 0x30, 0x01, 0x42, 0x9c, 0x01, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x4d, 0x69, 0x7a, 0x75, 0x63, 0x68, 0x69, 0x4c, 0x61, 0x62, 0x73, 0x2f, 0x73, 0x73, 0x68, 0x2d,
	0x6e, 0x65, 0x78, 0x75, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x67, 0x65, 0x6e, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x58, 0x58, 0xaa, 0x02, 0x08, 0x41, 0x67, 0x65,
	0x6e, 0x74, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x08, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x5c, 0x56, 0x31,
	0xe2, 0x02, 0x14, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x09, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x3a,
	0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_agent_v1_agent_proto_rawDescOnce sync.Once
	file_agent_v1_agent_proto_rawDescData = file_agent_v1_agent_proto_rawDesc
)

func file_agent_v1_agent_proto_rawDescGZIP() []byte {
	file_agent_v1_agent_proto_rawDescOnce.Do(func() {
		file_agent_v1_agent_proto_rawDescData = protoimpl.X.CompressGZIP(file_agent_v1_agent_proto_rawDescData)
	})
	return file_agent_v1_agent_proto_rawDescData
}

var file_agent_v1_agent_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_agent_v1_agent_proto_goTypes = []any{
	(*StreamResponse)(nil),           // 0: agent.v1.StreamResponse
	(*StreamRequest)(nil),            // 1: agent.v1.StreamRequest
	(*StreamResponse_Principal)(nil), // 2: agent.v1.StreamResponse.Principal
}
var file_agent_v1_agent_proto_depIdxs = []int32{
	2, // 0: agent.v1.StreamResponse.principals:type_name -> agent.v1.StreamResponse.Principal
	1, // 1: agent.v1.AgentService.Stream:input_type -> agent.v1.StreamRequest
	0, // 2: agent.v1.AgentService.Stream:output_type -> agent.v1.StreamResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_agent_v1_agent_proto_init() }
func file_agent_v1_agent_proto_init() {
	if File_agent_v1_agent_proto != nil {
		return
	}
	file_agent_v1_agent_proto_msgTypes[0].OneofWrappers = []any{}
	file_agent_v1_agent_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_agent_v1_agent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_agent_v1_agent_proto_goTypes,
		DependencyIndexes: file_agent_v1_agent_proto_depIdxs,
		MessageInfos:      file_agent_v1_agent_proto_msgTypes,
	}.Build()
	File_agent_v1_agent_proto = out.File
	file_agent_v1_agent_proto_rawDesc = nil
	file_agent_v1_agent_proto_goTypes = nil
	file_agent_v1_agent_proto_depIdxs = nil
}
