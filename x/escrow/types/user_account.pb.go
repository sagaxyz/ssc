// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ssc/escrow/user_account.proto

package types

import (
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type UserAccount struct {
	Chainlets []string     `protobuf:"bytes,1,rep,name=chainlets,proto3" json:"chainlets,omitempty"`
	Balance   []types.Coin `protobuf:"bytes,2,rep,name=balance,proto3" json:"balance"`
}

func (m *UserAccount) Reset()         { *m = UserAccount{} }
func (m *UserAccount) String() string { return proto.CompactTextString(m) }
func (*UserAccount) ProtoMessage()    {}
func (*UserAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_8276137147f41f8e, []int{0}
}
func (m *UserAccount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UserAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UserAccount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UserAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserAccount.Merge(m, src)
}
func (m *UserAccount) XXX_Size() int {
	return m.Size()
}
func (m *UserAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_UserAccount.DiscardUnknown(m)
}

var xxx_messageInfo_UserAccount proto.InternalMessageInfo

func (m *UserAccount) GetChainlets() []string {
	if m != nil {
		return m.Chainlets
	}
	return nil
}

func (m *UserAccount) GetBalance() []types.Coin {
	if m != nil {
		return m.Balance
	}
	return nil
}

func init() {
	proto.RegisterType((*UserAccount)(nil), "ssc.escrow.UserAccount")
}

func init() { proto.RegisterFile("ssc/escrow/user_account.proto", fileDescriptor_8276137147f41f8e) }

var fileDescriptor_8276137147f41f8e = []byte{
	// 242 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x8f, 0xb1, 0x4a, 0x03, 0x41,
	0x10, 0x86, 0xef, 0x8c, 0x28, 0xd9, 0x74, 0x87, 0x45, 0x0c, 0xba, 0x06, 0x41, 0x48, 0xb5, 0x4b,
	0xb4, 0xb2, 0x12, 0xe3, 0x1b, 0x04, 0x6c, 0x6c, 0x64, 0x76, 0x18, 0x2f, 0x07, 0xc9, 0x4e, 0xb8,
	0xd9, 0xd3, 0xc4, 0xa7, 0xf0, 0xb1, 0x52, 0xa6, 0xb4, 0x12, 0xb9, 0x7b, 0x11, 0xb9, 0xac, 0x21,
	0xdd, 0xf0, 0xff, 0x3f, 0x1f, 0xf3, 0xa9, 0x4b, 0x11, 0xb4, 0x24, 0x58, 0xf2, 0x87, 0xad, 0x84,
	0xca, 0x57, 0x40, 0xe4, 0xca, 0x07, 0xb3, 0x2c, 0x39, 0x70, 0xa6, 0x44, 0xd0, 0xc4, 0x7a, 0x70,
	0x96, 0x73, 0xce, 0xbb, 0xd8, 0xb6, 0x57, 0x5c, 0x0c, 0x34, 0xb2, 0x2c, 0x58, 0xac, 0x03, 0x21,
	0xfb, 0x3e, 0x76, 0x14, 0x60, 0x6c, 0x91, 0x0b, 0x1f, 0xfb, 0xeb, 0x37, 0xd5, 0x7b, 0x16, 0x2a,
	0x1f, 0x23, 0x36, 0xbb, 0x50, 0x5d, 0x9c, 0x41, 0xe1, 0xe7, 0x14, 0xa4, 0x9f, 0x0e, 0x3b, 0xa3,
	0xee, 0xf4, 0x10, 0x64, 0xf7, 0xea, 0xd4, 0xc1, 0x1c, 0x3c, 0x52, 0xff, 0x68, 0xd8, 0x19, 0xf5,
	0x6e, 0xcf, 0x4d, 0xc4, 0x9b, 0x16, 0x6f, 0xfe, 0xf1, 0xe6, 0x89, 0x0b, 0x3f, 0x39, 0xde, 0xfc,
	0x5c, 0x25, 0xd3, 0xfd, 0x7e, 0xf2, 0xb0, 0xa9, 0x75, 0xba, 0xad, 0x75, 0xfa, 0x5b, 0xeb, 0xf4,
	0xab, 0xd1, 0xc9, 0xb6, 0xd1, 0xc9, 0x77, 0xa3, 0x93, 0x97, 0x9b, 0xbc, 0x08, 0xb3, 0xca, 0x19,
	0xe4, 0x85, 0x15, 0xc8, 0x61, 0xb5, 0xfe, 0xb4, 0xad, 0xf5, 0x6a, 0xef, 0x1d, 0xd6, 0x4b, 0x12,
	0x77, 0xb2, 0xfb, 0xf7, 0xee, 0x2f, 0x00, 0x00, 0xff, 0xff, 0xa9, 0x27, 0x6d, 0x70, 0x12, 0x01,
	0x00, 0x00,
}

func (m *UserAccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UserAccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UserAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Balance) > 0 {
		for iNdEx := len(m.Balance) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Balance[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintUserAccount(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Chainlets) > 0 {
		for iNdEx := len(m.Chainlets) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Chainlets[iNdEx])
			copy(dAtA[i:], m.Chainlets[iNdEx])
			i = encodeVarintUserAccount(dAtA, i, uint64(len(m.Chainlets[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintUserAccount(dAtA []byte, offset int, v uint64) int {
	offset -= sovUserAccount(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *UserAccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Chainlets) > 0 {
		for _, s := range m.Chainlets {
			l = len(s)
			n += 1 + l + sovUserAccount(uint64(l))
		}
	}
	if len(m.Balance) > 0 {
		for _, e := range m.Balance {
			l = e.Size()
			n += 1 + l + sovUserAccount(uint64(l))
		}
	}
	return n
}

func sovUserAccount(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozUserAccount(x uint64) (n int) {
	return sovUserAccount(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *UserAccount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowUserAccount
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
			return fmt.Errorf("proto: UserAccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UserAccount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chainlets", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUserAccount
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
				return ErrInvalidLengthUserAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUserAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Chainlets = append(m.Chainlets, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Balance", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUserAccount
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
				return ErrInvalidLengthUserAccount
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthUserAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Balance = append(m.Balance, types.Coin{})
			if err := m.Balance[len(m.Balance)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipUserAccount(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthUserAccount
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
func skipUserAccount(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowUserAccount
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
					return 0, ErrIntOverflowUserAccount
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
					return 0, ErrIntOverflowUserAccount
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
				return 0, ErrInvalidLengthUserAccount
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupUserAccount
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthUserAccount
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthUserAccount        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowUserAccount          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupUserAccount = fmt.Errorf("proto: unexpected end of group")
)
