// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: proto/consensus.proto

package proto

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

type GenericConsensus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index       int64                    `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	M           int64                    `protobuf:"varint,2,opt,name=M,proto3" json:"M,omitempty"` // propose=1 , spreadE=2 , spreadC=3, gather=4, decide=5, commit=6
	S           int64                    `protobuf:"varint,3,opt,name=S,proto3" json:"S,omitempty"` // step
	P           *GenericConsensusValue   `protobuf:"bytes,4,opt,name=P,proto3" json:"P,omitempty"`
	E           []*GenericConsensusValue `protobuf:"bytes,5,rep,name=E,proto3" json:"E,omitempty"`
	C           []*GenericConsensusValue `protobuf:"bytes,6,rep,name=C,proto3" json:"C,omitempty"`
	D           bool                     `protobuf:"varint,7,opt,name=D,proto3" json:"D,omitempty"`
	DS          *GenericConsensusValue   `protobuf:"bytes,8,opt,name=DS,proto3" json:"DS,omitempty"`
	PR          int64                    `protobuf:"varint,9,opt,name=PR,proto3" json:"PR,omitempty"`                    // id of the proposer who decided this index
	Destination int64                    `protobuf:"varint,10,opt,name=destination,proto3" json:"destination,omitempty"` // 1 for the proposer and 2 for the recorder
}

func (x *GenericConsensus) Reset() {
	*x = GenericConsensus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_consensus_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenericConsensus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenericConsensus) ProtoMessage() {}

func (x *GenericConsensus) ProtoReflect() protoreflect.Message {
	mi := &file_proto_consensus_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenericConsensus.ProtoReflect.Descriptor instead.
func (*GenericConsensus) Descriptor() ([]byte, []int) {
	return file_proto_consensus_proto_rawDescGZIP(), []int{0}
}

func (x *GenericConsensus) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *GenericConsensus) GetM() int64 {
	if x != nil {
		return x.M
	}
	return 0
}

func (x *GenericConsensus) GetS() int64 {
	if x != nil {
		return x.S
	}
	return 0
}

func (x *GenericConsensus) GetP() *GenericConsensusValue {
	if x != nil {
		return x.P
	}
	return nil
}

func (x *GenericConsensus) GetE() []*GenericConsensusValue {
	if x != nil {
		return x.E
	}
	return nil
}

func (x *GenericConsensus) GetC() []*GenericConsensusValue {
	if x != nil {
		return x.C
	}
	return nil
}

func (x *GenericConsensus) GetD() bool {
	if x != nil {
		return x.D
	}
	return false
}

func (x *GenericConsensus) GetDS() *GenericConsensusValue {
	if x != nil {
		return x.DS
	}
	return nil
}

func (x *GenericConsensus) GetPR() int64 {
	if x != nil {
		return x.PR
	}
	return 0
}

func (x *GenericConsensus) GetDestination() int64 {
	if x != nil {
		return x.Destination
	}
	return 0
}

type GenericConsensusValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id  string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Fit int64  `protobuf:"varint,2,opt,name=fit,proto3" json:"fit,omitempty"`
}

func (x *GenericConsensusValue) Reset() {
	*x = GenericConsensusValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_consensus_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenericConsensusValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenericConsensusValue) ProtoMessage() {}

func (x *GenericConsensusValue) ProtoReflect() protoreflect.Message {
	mi := &file_proto_consensus_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenericConsensusValue.ProtoReflect.Descriptor instead.
func (*GenericConsensusValue) Descriptor() ([]byte, []int) {
	return file_proto_consensus_proto_rawDescGZIP(), []int{0, 0}
}

func (x *GenericConsensusValue) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GenericConsensusValue) GetFit() int64 {
	if x != nil {
		return x.Fit
	}
	return 0
}

var File_proto_consensus_proto protoreflect.FileDescriptor

var file_proto_consensus_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcd, 0x02, 0x0a, 0x10, 0x47, 0x65, 0x6e, 0x65,
	0x72, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x4d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x01, 0x4d,
	0x12, 0x0c, 0x0a, 0x01, 0x53, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x01, 0x53, 0x12, 0x25,
	0x0a, 0x01, 0x50, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x47, 0x65, 0x6e, 0x65,
	0x72, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x2e, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x01, 0x50, 0x12, 0x25, 0x0a, 0x01, 0x45, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e,
	0x73, 0x75, 0x73, 0x2e, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x01, 0x45, 0x12, 0x25, 0x0a, 0x01,
	0x43, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x69,
	0x63, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x2e, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x01, 0x43, 0x12, 0x0c, 0x0a, 0x01, 0x44, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x01,
	0x44, 0x12, 0x27, 0x0a, 0x02, 0x44, 0x53, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x47, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73,
	0x2e, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x02, 0x44, 0x53, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x52,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x50, 0x52, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0b, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x29, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x66, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x03, 0x66, 0x69, 0x74, 0x42, 0x08, 0x5a, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_consensus_proto_rawDescOnce sync.Once
	file_proto_consensus_proto_rawDescData = file_proto_consensus_proto_rawDesc
)

func file_proto_consensus_proto_rawDescGZIP() []byte {
	file_proto_consensus_proto_rawDescOnce.Do(func() {
		file_proto_consensus_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_consensus_proto_rawDescData)
	})
	return file_proto_consensus_proto_rawDescData
}

var file_proto_consensus_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_consensus_proto_goTypes = []interface{}{
	(*GenericConsensus)(nil),      // 0: GenericConsensus
	(*GenericConsensusValue)(nil), // 1: GenericConsensus.value
}
var file_proto_consensus_proto_depIdxs = []int32{
	1, // 0: GenericConsensus.P:type_name -> GenericConsensus.value
	1, // 1: GenericConsensus.E:type_name -> GenericConsensus.value
	1, // 2: GenericConsensus.C:type_name -> GenericConsensus.value
	1, // 3: GenericConsensus.DS:type_name -> GenericConsensus.value
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_consensus_proto_init() }
func file_proto_consensus_proto_init() {
	if File_proto_consensus_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_consensus_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenericConsensus); i {
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
		file_proto_consensus_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenericConsensusValue); i {
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
			RawDescriptor: file_proto_consensus_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_consensus_proto_goTypes,
		DependencyIndexes: file_proto_consensus_proto_depIdxs,
		MessageInfos:      file_proto_consensus_proto_msgTypes,
	}.Build()
	File_proto_consensus_proto = out.File
	file_proto_consensus_proto_rawDesc = nil
	file_proto_consensus_proto_goTypes = nil
	file_proto_consensus_proto_depIdxs = nil
}