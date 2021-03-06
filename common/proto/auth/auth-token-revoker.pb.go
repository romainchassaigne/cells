// Code generated by protoc-gen-go. DO NOT EDIT.
// source: auth-token-revoker.proto

/*
Package auth is a generated protocol buffer package.

It is generated from these files:
	auth-token-revoker.proto
	ldap.proto

It has these top-level messages:
	Token
	MatchInvalidTokenRequest
	MatchInvalidTokenResponse
	RevokeTokenRequest
	RevokeTokenResponse
	PruneTokensRequest
	PruneTokensResponse
	LdapSearchFilter
	LdapMapping
	LdapMemberOfMapping
	LdapServerConfig
*/
package auth

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// ============== MESSAGES ==========
type State int32

const (
	State_NO_MATCH State = 0
	State_REVOKED  State = 1
)

var State_name = map[int32]string{
	0: "NO_MATCH",
	1: "REVOKED",
}
var State_value = map[string]int32{
	"NO_MATCH": 0,
	"REVOKED":  1,
}

func (x State) String() string {
	return proto.EnumName(State_name, int32(x))
}
func (State) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Token struct {
	Value          string `protobuf:"bytes,1,opt,name=Value" json:"Value,omitempty"`
	AdditionalInfo string `protobuf:"bytes,2,opt,name=AdditionalInfo" json:"AdditionalInfo,omitempty"`
}

func (m *Token) Reset()                    { *m = Token{} }
func (m *Token) String() string            { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()               {}
func (*Token) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Token) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *Token) GetAdditionalInfo() string {
	if m != nil {
		return m.AdditionalInfo
	}
	return ""
}

type MatchInvalidTokenRequest struct {
	Token string `protobuf:"bytes,1,opt,name=Token" json:"Token,omitempty"`
}

func (m *MatchInvalidTokenRequest) Reset()                    { *m = MatchInvalidTokenRequest{} }
func (m *MatchInvalidTokenRequest) String() string            { return proto.CompactTextString(m) }
func (*MatchInvalidTokenRequest) ProtoMessage()               {}
func (*MatchInvalidTokenRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *MatchInvalidTokenRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type MatchInvalidTokenResponse struct {
	State          State  `protobuf:"varint,1,opt,name=State,enum=auth.State" json:"State,omitempty"`
	RevocationInfo string `protobuf:"bytes,2,opt,name=RevocationInfo" json:"RevocationInfo,omitempty"`
}

func (m *MatchInvalidTokenResponse) Reset()                    { *m = MatchInvalidTokenResponse{} }
func (m *MatchInvalidTokenResponse) String() string            { return proto.CompactTextString(m) }
func (*MatchInvalidTokenResponse) ProtoMessage()               {}
func (*MatchInvalidTokenResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *MatchInvalidTokenResponse) GetState() State {
	if m != nil {
		return m.State
	}
	return State_NO_MATCH
}

func (m *MatchInvalidTokenResponse) GetRevocationInfo() string {
	if m != nil {
		return m.RevocationInfo
	}
	return ""
}

type RevokeTokenRequest struct {
	Token *Token `protobuf:"bytes,1,opt,name=Token" json:"Token,omitempty"`
}

func (m *RevokeTokenRequest) Reset()                    { *m = RevokeTokenRequest{} }
func (m *RevokeTokenRequest) String() string            { return proto.CompactTextString(m) }
func (*RevokeTokenRequest) ProtoMessage()               {}
func (*RevokeTokenRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *RevokeTokenRequest) GetToken() *Token {
	if m != nil {
		return m.Token
	}
	return nil
}

type RevokeTokenResponse struct {
	Success bool `protobuf:"varint,1,opt,name=Success" json:"Success,omitempty"`
}

func (m *RevokeTokenResponse) Reset()                    { *m = RevokeTokenResponse{} }
func (m *RevokeTokenResponse) String() string            { return proto.CompactTextString(m) }
func (*RevokeTokenResponse) ProtoMessage()               {}
func (*RevokeTokenResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *RevokeTokenResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type PruneTokensRequest struct {
}

func (m *PruneTokensRequest) Reset()                    { *m = PruneTokensRequest{} }
func (m *PruneTokensRequest) String() string            { return proto.CompactTextString(m) }
func (*PruneTokensRequest) ProtoMessage()               {}
func (*PruneTokensRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type PruneTokensResponse struct {
	Tokens []string `protobuf:"bytes,1,rep,name=tokens" json:"tokens,omitempty"`
}

func (m *PruneTokensResponse) Reset()                    { *m = PruneTokensResponse{} }
func (m *PruneTokensResponse) String() string            { return proto.CompactTextString(m) }
func (*PruneTokensResponse) ProtoMessage()               {}
func (*PruneTokensResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *PruneTokensResponse) GetTokens() []string {
	if m != nil {
		return m.Tokens
	}
	return nil
}

func init() {
	proto.RegisterType((*Token)(nil), "auth.Token")
	proto.RegisterType((*MatchInvalidTokenRequest)(nil), "auth.MatchInvalidTokenRequest")
	proto.RegisterType((*MatchInvalidTokenResponse)(nil), "auth.MatchInvalidTokenResponse")
	proto.RegisterType((*RevokeTokenRequest)(nil), "auth.RevokeTokenRequest")
	proto.RegisterType((*RevokeTokenResponse)(nil), "auth.RevokeTokenResponse")
	proto.RegisterType((*PruneTokensRequest)(nil), "auth.PruneTokensRequest")
	proto.RegisterType((*PruneTokensResponse)(nil), "auth.PruneTokensResponse")
	proto.RegisterEnum("auth.State", State_name, State_value)
}

func init() { proto.RegisterFile("auth-token-revoker.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 355 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x4d, 0x4f, 0xf2, 0x40,
	0x10, 0xc7, 0xe9, 0xf3, 0xc8, 0xdb, 0x94, 0x10, 0xb2, 0x10, 0x53, 0x38, 0xa8, 0xec, 0xc1, 0x10,
	0x13, 0xd0, 0xe0, 0xc1, 0xa3, 0x21, 0x42, 0x22, 0x31, 0x88, 0x2e, 0x84, 0xab, 0x59, 0xcb, 0x12,
	0x08, 0xa4, 0x8b, 0xed, 0x2e, 0x9f, 0xda, 0x0f, 0x61, 0xf6, 0xa5, 0xa6, 0x95, 0x72, 0x9c, 0xf9,
	0xcf, 0xfc, 0xf7, 0x37, 0x33, 0x0b, 0x1e, 0x95, 0x62, 0xdd, 0x15, 0x7c, 0xcb, 0x82, 0x6e, 0xc8,
	0x0e, 0x7c, 0xcb, 0xc2, 0xde, 0x3e, 0xe4, 0x82, 0xa3, 0x33, 0xa5, 0xe0, 0x11, 0xe4, 0xe7, 0x4a,
	0x44, 0x0d, 0xc8, 0x2f, 0xe8, 0x4e, 0x32, 0xcf, 0xb9, 0x72, 0x3a, 0x65, 0x62, 0x02, 0x74, 0x0d,
	0xd5, 0xc1, 0x72, 0xb9, 0x11, 0x1b, 0x1e, 0xd0, 0xdd, 0x38, 0x58, 0x71, 0xef, 0x9f, 0x96, 0xff,
	0x64, 0xf1, 0x1d, 0x78, 0x13, 0x2a, 0xfc, 0xf5, 0x38, 0x38, 0xd0, 0xdd, 0x66, 0xa9, 0x2d, 0x09,
	0xfb, 0x92, 0x2c, 0x12, 0xca, 0x59, 0xc7, 0xb1, 0xb3, 0x0e, 0xf0, 0x0a, 0x9a, 0x19, 0x1d, 0xd1,
	0x9e, 0x07, 0x11, 0x43, 0x6d, 0xc8, 0xcf, 0x04, 0x15, 0x06, 0xa6, 0xda, 0x77, 0x7b, 0x8a, 0xb5,
	0xa7, 0x53, 0xc4, 0x28, 0x8a, 0x8c, 0xb0, 0x03, 0xf7, 0xa9, 0xa2, 0x48, 0x92, 0xa5, 0xb3, 0xf8,
	0x01, 0x10, 0xd1, 0x73, 0xa7, 0x98, 0xda, 0x49, 0x26, 0x37, 0x7e, 0xc0, 0x94, 0x58, 0xc0, 0x5b,
	0xa8, 0xa7, 0x1a, 0x2d, 0x9a, 0x07, 0xc5, 0x99, 0xf4, 0x7d, 0x16, 0x45, 0xba, 0xb7, 0x44, 0xe2,
	0x10, 0x37, 0x00, 0xbd, 0x85, 0x32, 0x30, 0xf5, 0x91, 0x7d, 0x09, 0x77, 0xa1, 0x9e, 0xca, 0x5a,
	0x9b, 0x73, 0x28, 0xe8, 0xa3, 0x28, 0x97, 0xff, 0x9d, 0x32, 0xb1, 0xd1, 0x0d, 0xb6, 0x93, 0xa3,
	0x0a, 0x94, 0x5e, 0xa7, 0x1f, 0x93, 0xc1, 0xfc, 0xe9, 0xb9, 0x96, 0x43, 0x2e, 0x14, 0xc9, 0x68,
	0x31, 0x7d, 0x19, 0x0d, 0x6b, 0x4e, 0xff, 0xdb, 0x81, 0xda, 0x40, 0x8a, 0xb5, 0x05, 0xd3, 0x47,
	0x45, 0xef, 0x50, 0x49, 0xee, 0x13, 0x5d, 0x98, 0x91, 0x4e, 0x5d, 0xa5, 0x75, 0x79, 0x52, 0x37,
	0x84, 0x38, 0x87, 0x1e, 0xa1, 0x60, 0xdc, 0x91, 0x67, 0x8a, 0x8f, 0x17, 0xd9, 0x6a, 0x66, 0x28,
	0xbf, 0x06, 0x43, 0x70, 0x13, 0xb3, 0xc7, 0x2e, 0xc7, 0x4b, 0x8a, 0x5d, 0x32, 0x16, 0x85, 0x73,
	0x9f, 0x05, 0xfd, 0x5f, 0xef, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0xb5, 0xf9, 0xfd, 0x9b, 0xcb,
	0x02, 0x00, 0x00,
}
