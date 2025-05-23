// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ssc/chainlet/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "github.com/gogo/protobuf/types"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Params defines the parameters for the module.
type Params struct {
	ChainletStackProtections         bool   `protobuf:"varint,1,opt,name=chainletStackProtections,proto3" json:"chainletStackProtections,omitempty"`
	NEpochDeposit                    string `protobuf:"bytes,2,opt,name=nEpochDeposit,proto3" json:"nEpochDeposit,omitempty"`
	AutomaticChainletUpgrades        bool   `protobuf:"varint,3,opt,name=automaticChainletUpgrades,proto3" json:"automaticChainletUpgrades,omitempty"`
	AutomaticChainletUpgradeInterval int64  `protobuf:"varint,4,opt,name=automaticChainletUpgradeInterval,proto3" json:"automaticChainletUpgradeInterval,omitempty"`
	// Delays launch to give validators time to set the consumer key
	LaunchDelay  time.Duration `protobuf:"bytes,5,opt,name=launchDelay,proto3,stdduration" json:"launchDelay"`
	MaxChainlets uint64        `protobuf:"varint,6,opt,name=maxChainlets,proto3" json:"maxChainlets,omitempty"`
	// If true, all new chainlets will be launched as CCV consumers
	CcvConsumerEnabled bool `protobuf:"varint,7,opt,name=ccvConsumerEnabled,proto3" json:"ccvConsumerEnabled,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ba1040c6477ee7f, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetChainletStackProtections() bool {
	if m != nil {
		return m.ChainletStackProtections
	}
	return false
}

func (m *Params) GetNEpochDeposit() string {
	if m != nil {
		return m.NEpochDeposit
	}
	return ""
}

func (m *Params) GetAutomaticChainletUpgrades() bool {
	if m != nil {
		return m.AutomaticChainletUpgrades
	}
	return false
}

func (m *Params) GetAutomaticChainletUpgradeInterval() int64 {
	if m != nil {
		return m.AutomaticChainletUpgradeInterval
	}
	return 0
}

func (m *Params) GetLaunchDelay() time.Duration {
	if m != nil {
		return m.LaunchDelay
	}
	return 0
}

func (m *Params) GetMaxChainlets() uint64 {
	if m != nil {
		return m.MaxChainlets
	}
	return 0
}

func (m *Params) GetCcvConsumerEnabled() bool {
	if m != nil {
		return m.CcvConsumerEnabled
	}
	return false
}

func init() {
	proto.RegisterType((*Params)(nil), "ssc.chainlet.Params")
}

func init() { proto.RegisterFile("ssc/chainlet/params.proto", fileDescriptor_3ba1040c6477ee7f) }

var fileDescriptor_3ba1040c6477ee7f = []byte{
	// 358 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xbf, 0x4e, 0xf3, 0x30,
	0x14, 0xc5, 0xe3, 0xaf, 0xfd, 0x4a, 0x71, 0xcb, 0x62, 0x31, 0xb8, 0x1d, 0xd2, 0xa8, 0x42, 0x22,
	0x53, 0x22, 0xc1, 0x86, 0x58, 0xe8, 0x9f, 0x01, 0xa6, 0x2a, 0x88, 0x85, 0xed, 0xc6, 0x35, 0x69,
	0x44, 0x12, 0x47, 0xb6, 0x53, 0xb5, 0xbc, 0x04, 0x8c, 0x3c, 0x52, 0xc7, 0x8e, 0x4c, 0x80, 0xda,
	0x17, 0x41, 0x49, 0x1b, 0x41, 0x25, 0x2a, 0x36, 0xfb, 0xfe, 0x8e, 0xcf, 0xf5, 0x3d, 0xba, 0xb8,
	0xa5, 0x14, 0x73, 0xd9, 0x04, 0xc2, 0x24, 0xe2, 0xda, 0x4d, 0x41, 0x42, 0xac, 0x9c, 0x54, 0x0a,
	0x2d, 0x48, 0x53, 0x29, 0xe6, 0x94, 0xa8, 0x7d, 0x1c, 0x88, 0x40, 0x14, 0xc0, 0xcd, 0x4f, 0x1b,
	0x4d, 0xdb, 0x0c, 0x84, 0x08, 0x22, 0xee, 0x16, 0x37, 0x3f, 0x7b, 0x70, 0xc7, 0x99, 0x04, 0x1d,
	0x8a, 0x64, 0xc3, 0xbb, 0xcf, 0x15, 0x5c, 0x1b, 0x15, 0xa6, 0xe4, 0x02, 0xd3, 0xd2, 0xec, 0x56,
	0x03, 0x7b, 0x1c, 0x49, 0xa1, 0x39, 0xcb, 0xb5, 0x8a, 0x22, 0x0b, 0xd9, 0x75, 0x6f, 0x2f, 0x27,
	0x27, 0xf8, 0x28, 0x19, 0xa6, 0x82, 0x4d, 0x06, 0x3c, 0x15, 0x2a, 0xd4, 0xf4, 0x9f, 0x85, 0xec,
	0x43, 0x6f, 0xb7, 0x48, 0x2e, 0x71, 0x0b, 0x32, 0x2d, 0x62, 0xd0, 0x21, 0xeb, 0x6f, 0xad, 0xee,
	0xd2, 0x40, 0xc2, 0x98, 0x2b, 0x5a, 0x29, 0x5a, 0xec, 0x17, 0x90, 0x1b, 0x6c, 0xed, 0x83, 0xd7,
	0x89, 0xe6, 0x72, 0x0a, 0x11, 0xad, 0x5a, 0xc8, 0xae, 0x78, 0x7f, 0xea, 0xc8, 0x10, 0x37, 0x22,
	0xc8, 0x92, 0xfc, 0x6b, 0x11, 0xcc, 0xe9, 0x7f, 0x0b, 0xd9, 0x8d, 0xb3, 0x96, 0xb3, 0x09, 0xcb,
	0x29, 0xc3, 0x72, 0x06, 0xdb, 0xb0, 0x7a, 0xf5, 0xc5, 0x7b, 0xc7, 0x78, 0xfd, 0xe8, 0x20, 0xef,
	0xe7, 0x3b, 0xd2, 0xc5, 0xcd, 0x18, 0x66, 0x65, 0x13, 0x45, 0x6b, 0x16, 0xb2, 0xab, 0xde, 0x4e,
	0x8d, 0x38, 0x98, 0x30, 0x36, 0xed, 0x8b, 0x44, 0x65, 0x31, 0x97, 0xc3, 0x04, 0xfc, 0x88, 0x8f,
	0xe9, 0x41, 0x31, 0xed, 0x2f, 0xa4, 0x77, 0xb5, 0x58, 0x99, 0x68, 0xb9, 0x32, 0xd1, 0xe7, 0xca,
	0x44, 0x2f, 0x6b, 0xd3, 0x58, 0xae, 0x4d, 0xe3, 0x6d, 0x6d, 0x1a, 0xf7, 0xa7, 0x41, 0xa8, 0x27,
	0x99, 0xef, 0x30, 0x11, 0xbb, 0x0a, 0x02, 0x98, 0xcd, 0x9f, 0xdc, 0x7c, 0x3b, 0x66, 0xdf, 0xfb,
	0xa1, 0xe7, 0x29, 0x57, 0x7e, 0xad, 0x18, 0xe0, 0xfc, 0x2b, 0x00, 0x00, 0xff, 0xff, 0x74, 0xe2,
	0x0f, 0xd8, 0x3c, 0x02, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.CcvConsumerEnabled {
		i--
		if m.CcvConsumerEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x38
	}
	if m.MaxChainlets != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxChainlets))
		i--
		dAtA[i] = 0x30
	}
	n1, err1 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.LaunchDelay, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.LaunchDelay):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintParams(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x2a
	if m.AutomaticChainletUpgradeInterval != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.AutomaticChainletUpgradeInterval))
		i--
		dAtA[i] = 0x20
	}
	if m.AutomaticChainletUpgrades {
		i--
		if m.AutomaticChainletUpgrades {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if len(m.NEpochDeposit) > 0 {
		i -= len(m.NEpochDeposit)
		copy(dAtA[i:], m.NEpochDeposit)
		i = encodeVarintParams(dAtA, i, uint64(len(m.NEpochDeposit)))
		i--
		dAtA[i] = 0x12
	}
	if m.ChainletStackProtections {
		i--
		if m.ChainletStackProtections {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ChainletStackProtections {
		n += 2
	}
	l = len(m.NEpochDeposit)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.AutomaticChainletUpgrades {
		n += 2
	}
	if m.AutomaticChainletUpgradeInterval != 0 {
		n += 1 + sovParams(uint64(m.AutomaticChainletUpgradeInterval))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.LaunchDelay)
	n += 1 + l + sovParams(uint64(l))
	if m.MaxChainlets != 0 {
		n += 1 + sovParams(uint64(m.MaxChainlets))
	}
	if m.CcvConsumerEnabled {
		n += 2
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainletStackProtections", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.ChainletStackProtections = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NEpochDeposit", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NEpochDeposit = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AutomaticChainletUpgrades", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.AutomaticChainletUpgrades = bool(v != 0)
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AutomaticChainletUpgradeInterval", wireType)
			}
			m.AutomaticChainletUpgradeInterval = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AutomaticChainletUpgradeInterval |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LaunchDelay", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.LaunchDelay, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxChainlets", wireType)
			}
			m.MaxChainlets = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxChainlets |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CcvConsumerEnabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.CcvConsumerEnabled = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowParams
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowParams
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
