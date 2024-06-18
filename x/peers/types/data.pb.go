// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ssc/peers/data.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

type Data struct {
	Updated   time.Time `protobuf:"bytes,1,opt,name=updated,proto3,stdtime" json:"updated"`
	Addresses []string  `protobuf:"bytes,2,rep,name=addresses,proto3" json:"addresses,omitempty"`
}

func (m *Data) Reset()         { *m = Data{} }
func (m *Data) String() string { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()    {}
func (*Data) Descriptor() ([]byte, []int) {
	return fileDescriptor_24fad99c7e28aed4, []int{0}
}
func (m *Data) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Data) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Data.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Data) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Data.Merge(m, src)
}
func (m *Data) XXX_Size() int {
	return m.Size()
}
func (m *Data) XXX_DiscardUnknown() {
	xxx_messageInfo_Data.DiscardUnknown(m)
}

var xxx_messageInfo_Data proto.InternalMessageInfo

func (m *Data) GetUpdated() time.Time {
	if m != nil {
		return m.Updated
	}
	return time.Time{}
}

func (m *Data) GetAddresses() []string {
	if m != nil {
		return m.Addresses
	}
	return nil
}

func init() {
	proto.RegisterType((*Data)(nil), "ssc.peers.Data")
}

func init() { proto.RegisterFile("ssc/peers/data.proto", fileDescriptor_24fad99c7e28aed4) }

var fileDescriptor_24fad99c7e28aed4 = []byte{
	// 243 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0x63, 0x40, 0x40, 0xc2, 0x16, 0x75, 0x28, 0x11, 0x72, 0x2a, 0xc4, 0xd0, 0xc9, 0x96,
	0x60, 0xef, 0x50, 0xf1, 0x0b, 0x2a, 0x26, 0x16, 0x74, 0x89, 0x8d, 0xa9, 0x44, 0x74, 0x56, 0xee,
	0x22, 0xb5, 0xfc, 0x8a, 0xfe, 0xac, 0x8e, 0x1d, 0x99, 0x00, 0x25, 0x7f, 0x04, 0x35, 0x4e, 0xd4,
	0xcd, 0xef, 0xf9, 0xbd, 0x4f, 0xa7, 0x97, 0x4c, 0x88, 0x4a, 0xed, 0xad, 0xad, 0x49, 0x1b, 0x60,
	0x50, 0xbe, 0x46, 0xc6, 0x34, 0x26, 0x2a, 0x55, 0xef, 0x66, 0x13, 0x87, 0x0e, 0x7b, 0x57, 0x1f,
	0x5f, 0x21, 0x90, 0xe5, 0x0e, 0xd1, 0x7d, 0x5a, 0xdd, 0xab, 0xa2, 0x79, 0xd7, 0xbc, 0xae, 0x2c,
	0x31, 0x54, 0x7e, 0x08, 0xdc, 0x96, 0x48, 0x15, 0xd2, 0x5b, 0x68, 0x06, 0x11, 0xbe, 0xee, 0x4d,
	0x72, 0xf1, 0x0c, 0x0c, 0xe9, 0x22, 0xb9, 0x6a, 0xbc, 0x01, 0xb6, 0x66, 0x2a, 0x66, 0x62, 0x7e,
	0xf3, 0x98, 0xa9, 0x40, 0x55, 0x23, 0x55, 0xbd, 0x8c, 0xd4, 0xe5, 0xf5, 0xfe, 0x27, 0x8f, 0x76,
	0xbf, 0xb9, 0x58, 0x8d, 0xa5, 0xf4, 0x2e, 0x89, 0xc1, 0x98, 0xda, 0x12, 0x59, 0x9a, 0x9e, 0xcd,
	0xce, 0xe7, 0xf1, 0xea, 0x64, 0x2c, 0x17, 0xfb, 0x56, 0x8a, 0x43, 0x2b, 0xc5, 0x5f, 0x2b, 0xc5,
	0xae, 0x93, 0xd1, 0xa1, 0x93, 0xd1, 0x77, 0x27, 0xa3, 0xd7, 0x07, 0xb7, 0xe6, 0x8f, 0xa6, 0x50,
	0x25, 0x56, 0x9a, 0xc0, 0xc1, 0x66, 0xfb, 0xa5, 0x8f, 0x2b, 0x6c, 0x86, 0x1d, 0x78, 0xeb, 0x2d,
	0x15, 0x97, 0xfd, 0x11, 0x4f, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0x29, 0xc5, 0x86, 0x5c, 0x21,
	0x01, 0x00, 0x00,
}

func (m *Data) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Data) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Data) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Addresses) > 0 {
		for iNdEx := len(m.Addresses) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Addresses[iNdEx])
			copy(dAtA[i:], m.Addresses[iNdEx])
			i = encodeVarintData(dAtA, i, uint64(len(m.Addresses[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.Updated, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Updated):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintData(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintData(dAtA []byte, offset int, v uint64) int {
	offset -= sovData(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Data) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Updated)
	n += 1 + l + sovData(uint64(l))
	if len(m.Addresses) > 0 {
		for _, s := range m.Addresses {
			l = len(s)
			n += 1 + l + sovData(uint64(l))
		}
	}
	return n
}

func sovData(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozData(x uint64) (n int) {
	return sovData(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Data) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowData
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
			return fmt.Errorf("proto: Data: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Data: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Updated", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowData
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
				return ErrInvalidLengthData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.Updated, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Addresses", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowData
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
				return ErrInvalidLengthData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Addresses = append(m.Addresses, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthData
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
func skipData(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowData
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
					return 0, ErrIntOverflowData
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
					return 0, ErrIntOverflowData
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
				return 0, ErrInvalidLengthData
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupData
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthData
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthData        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowData          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupData = fmt.Errorf("proto: unexpected end of group")
)
