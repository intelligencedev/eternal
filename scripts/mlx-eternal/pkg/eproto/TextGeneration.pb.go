// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: TextGeneration.proto

package eproto

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

// The request message containing the user's prompt.
type TextRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Prompt string `protobuf:"bytes,1,opt,name=prompt,proto3" json:"prompt,omitempty"`
}

func (x *TextRequest) Reset() {
	*x = TextRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_TextGeneration_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TextRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TextRequest) ProtoMessage() {}

func (x *TextRequest) ProtoReflect() protoreflect.Message {
	mi := &file_TextGeneration_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TextRequest.ProtoReflect.Descriptor instead.
func (*TextRequest) Descriptor() ([]byte, []int) {
	return file_TextGeneration_proto_rawDescGZIP(), []int{0}
}

func (x *TextRequest) GetPrompt() string {
	if x != nil {
		return x.Prompt
	}
	return ""
}

// The response message containing the generated text.
type TextResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response string `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *TextResponse) Reset() {
	*x = TextResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_TextGeneration_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TextResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TextResponse) ProtoMessage() {}

func (x *TextResponse) ProtoReflect() protoreflect.Message {
	mi := &file_TextGeneration_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TextResponse.ProtoReflect.Descriptor instead.
func (*TextResponse) Descriptor() ([]byte, []int) {
	return file_TextGeneration_proto_rawDescGZIP(), []int{1}
}

func (x *TextResponse) GetResponse() string {
	if x != nil {
		return x.Response
	}
	return ""
}

var File_TextGeneration_proto protoreflect.FileDescriptor

var file_TextGeneration_proto_rawDesc = []byte{
	0x0a, 0x14, 0x54, 0x65, 0x78, 0x74, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x74, 0x65, 0x78, 0x74, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x6f, 0x72, 0x22, 0x25, 0x0a, 0x0b, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x22, 0x2a, 0x0a, 0x0c,
	0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x60, 0x0a, 0x0d, 0x54, 0x65, 0x78, 0x74,
	0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x4f, 0x0a, 0x12, 0x47, 0x65, 0x6e,
	0x65, 0x72, 0x61, 0x74, 0x65, 0x54, 0x65, 0x78, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12,
	0x1a, 0x2e, 0x74, 0x65, 0x78, 0x74, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e,
	0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x74, 0x65,
	0x78, 0x74, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x54, 0x65, 0x78, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f,
	0x65, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_TextGeneration_proto_rawDescOnce sync.Once
	file_TextGeneration_proto_rawDescData = file_TextGeneration_proto_rawDesc
)

func file_TextGeneration_proto_rawDescGZIP() []byte {
	file_TextGeneration_proto_rawDescOnce.Do(func() {
		file_TextGeneration_proto_rawDescData = protoimpl.X.CompressGZIP(file_TextGeneration_proto_rawDescData)
	})
	return file_TextGeneration_proto_rawDescData
}

var file_TextGeneration_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_TextGeneration_proto_goTypes = []interface{}{
	(*TextRequest)(nil),  // 0: textgenerator.TextRequest
	(*TextResponse)(nil), // 1: textgenerator.TextResponse
}
var file_TextGeneration_proto_depIdxs = []int32{
	0, // 0: textgenerator.TextGenerator.GenerateTextStream:input_type -> textgenerator.TextRequest
	1, // 1: textgenerator.TextGenerator.GenerateTextStream:output_type -> textgenerator.TextResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_TextGeneration_proto_init() }
func file_TextGeneration_proto_init() {
	if File_TextGeneration_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_TextGeneration_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TextRequest); i {
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
		file_TextGeneration_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TextResponse); i {
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
			RawDescriptor: file_TextGeneration_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_TextGeneration_proto_goTypes,
		DependencyIndexes: file_TextGeneration_proto_depIdxs,
		MessageInfos:      file_TextGeneration_proto_msgTypes,
	}.Build()
	File_TextGeneration_proto = out.File
	file_TextGeneration_proto_rawDesc = nil
	file_TextGeneration_proto_goTypes = nil
	file_TextGeneration_proto_depIdxs = nil
}
