// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ssc/escrow/events.proto

package types

import (
	fmt "fmt"
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

type EventDeposit struct {
	User     string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Chainlet string `protobuf:"bytes,2,opt,name=chainlet,proto3" json:"chainlet,omitempty"`
	Amount   string `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount,omitempty"`
	NewTotal string `protobuf:"bytes,4,opt,name=newTotal,proto3" json:"newTotal,omitempty"`
}

func (m *EventDeposit) Reset()         { *m = EventDeposit{} }
func (m *EventDeposit) String() string { return proto.CompactTextString(m) }
func (*EventDeposit) ProtoMessage()    {}
func (*EventDeposit) Descriptor() ([]byte, []int) {
	return fileDescriptor_36568dd0af364364, []int{0}
}
func (m *EventDeposit) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventDeposit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventDeposit.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventDeposit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventDeposit.Merge(m, src)
}
func (m *EventDeposit) XXX_Size() int {
	return m.Size()
}
func (m *EventDeposit) XXX_DiscardUnknown() {
	xxx_messageInfo_EventDeposit.DiscardUnknown(m)
}

var xxx_messageInfo_EventDeposit proto.InternalMessageInfo

func (m *EventDeposit) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *EventDeposit) GetChainlet() string {
	if m != nil {
		return m.Chainlet
	}
	return ""
}

func (m *EventDeposit) GetAmount() string {
	if m != nil {
		return m.Amount
	}
	return ""
}

func (m *EventDeposit) GetNewTotal() string {
	if m != nil {
		return m.NewTotal
	}
	return ""
}

type EventWithdraw struct {
	User      string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Chainlet  string `protobuf:"bytes,2,opt,name=chainlet,proto3" json:"chainlet,omitempty"`
	Amount    string `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount,omitempty"`
	Remaining string `protobuf:"bytes,4,opt,name=remaining,proto3" json:"remaining,omitempty"`
}

func (m *EventWithdraw) Reset()         { *m = EventWithdraw{} }
func (m *EventWithdraw) String() string { return proto.CompactTextString(m) }
func (*EventWithdraw) ProtoMessage()    {}
func (*EventWithdraw) Descriptor() ([]byte, []int) {
	return fileDescriptor_36568dd0af364364, []int{1}
}
func (m *EventWithdraw) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventWithdraw) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventWithdraw.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventWithdraw) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventWithdraw.Merge(m, src)
}
func (m *EventWithdraw) XXX_Size() int {
	return m.Size()
}
func (m *EventWithdraw) XXX_DiscardUnknown() {
	xxx_messageInfo_EventWithdraw.DiscardUnknown(m)
}

var xxx_messageInfo_EventWithdraw proto.InternalMessageInfo

func (m *EventWithdraw) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *EventWithdraw) GetChainlet() string {
	if m != nil {
		return m.Chainlet
	}
	return ""
}

func (m *EventWithdraw) GetAmount() string {
	if m != nil {
		return m.Amount
	}
	return ""
}

func (m *EventWithdraw) GetRemaining() string {
	if m != nil {
		return m.Remaining
	}
	return ""
}

func init() {
	proto.RegisterType((*EventDeposit)(nil), "ssc.escrow.EventDeposit")
	proto.RegisterType((*EventWithdraw)(nil), "ssc.escrow.EventWithdraw")
}

func init() { proto.RegisterFile("ssc/escrow/events.proto", fileDescriptor_36568dd0af364364) }

var fileDescriptor_36568dd0af364364 = []byte{
	// 256 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x90, 0x41, 0x4a, 0xc4, 0x30,
	0x18, 0x85, 0x1b, 0x1d, 0x06, 0x27, 0xe8, 0x26, 0x88, 0x96, 0x22, 0x41, 0x06, 0x04, 0x57, 0xcd,
	0xc2, 0x03, 0x08, 0xa2, 0x17, 0x10, 0x41, 0x70, 0x97, 0x89, 0x21, 0x0d, 0x4c, 0x93, 0x92, 0xff,
	0xaf, 0x9d, 0xf1, 0x14, 0x1e, 0xcb, 0xe5, 0x2c, 0x5d, 0x4a, 0x7b, 0x11, 0x69, 0xda, 0xaa, 0x07,
	0x98, 0xdd, 0xff, 0xde, 0xcb, 0x97, 0x07, 0x8f, 0x9e, 0x03, 0x28, 0xa1, 0x41, 0x05, 0xdf, 0x08,
	0xfd, 0xa6, 0x1d, 0x42, 0x5e, 0x05, 0x8f, 0x9e, 0x51, 0x00, 0x95, 0x0f, 0x41, 0x76, 0x6a, 0xbc,
	0xf1, 0xd1, 0x16, 0xfd, 0x35, 0xbc, 0xc8, 0xfe, 0xa3, 0x95, 0x0c, 0xb2, 0x1c, 0xd1, 0x65, 0xa0,
	0xc7, 0x0f, 0xfd, 0x57, 0xf7, 0xba, 0xf2, 0x60, 0x91, 0x31, 0x3a, 0xab, 0x41, 0x87, 0x94, 0x5c,
	0x92, 0xeb, 0xc5, 0x63, 0xbc, 0x59, 0x46, 0x8f, 0x54, 0x21, 0xad, 0x5b, 0x6b, 0x4c, 0x0f, 0xa2,
	0xff, 0xab, 0xd9, 0x19, 0x9d, 0xcb, 0xd2, 0xd7, 0x0e, 0xd3, 0xc3, 0x98, 0x8c, 0xaa, 0x67, 0x9c,
	0x6e, 0x9e, 0x3c, 0xca, 0x75, 0x3a, 0x1b, 0x98, 0x49, 0x2f, 0x6b, 0x7a, 0x12, 0x3b, 0x9f, 0x2d,
	0x16, 0xaf, 0x41, 0x36, 0x7b, 0x2b, 0xbd, 0xa0, 0x8b, 0xa0, 0x4b, 0x69, 0x9d, 0x75, 0x66, 0x6c,
	0xfd, 0x33, 0xee, 0x6e, 0x3f, 0x5b, 0x4e, 0x76, 0x2d, 0x27, 0xdf, 0x2d, 0x27, 0x1f, 0x1d, 0x4f,
	0x76, 0x1d, 0x4f, 0xbe, 0x3a, 0x9e, 0xbc, 0x5c, 0x19, 0x8b, 0x45, 0xbd, 0xca, 0x95, 0x2f, 0x05,
	0x48, 0x23, 0x37, 0xdb, 0x77, 0xd1, 0x0f, 0xb6, 0x99, 0x26, 0xc3, 0x6d, 0xa5, 0x61, 0x35, 0x8f,
	0x93, 0xdd, 0xfc, 0x04, 0x00, 0x00, 0xff, 0xff, 0x3b, 0xd0, 0x99, 0x0c, 0x88, 0x01, 0x00, 0x00,
}

func (m *EventDeposit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventDeposit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventDeposit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.NewTotal) > 0 {
		i -= len(m.NewTotal)
		copy(dAtA[i:], m.NewTotal)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.NewTotal)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Amount) > 0 {
		i -= len(m.Amount)
		copy(dAtA[i:], m.Amount)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Amount)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Chainlet) > 0 {
		i -= len(m.Chainlet)
		copy(dAtA[i:], m.Chainlet)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Chainlet)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.User) > 0 {
		i -= len(m.User)
		copy(dAtA[i:], m.User)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.User)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EventWithdraw) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventWithdraw) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventWithdraw) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Remaining) > 0 {
		i -= len(m.Remaining)
		copy(dAtA[i:], m.Remaining)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Remaining)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Amount) > 0 {
		i -= len(m.Amount)
		copy(dAtA[i:], m.Amount)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Amount)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Chainlet) > 0 {
		i -= len(m.Chainlet)
		copy(dAtA[i:], m.Chainlet)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Chainlet)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.User) > 0 {
		i -= len(m.User)
		copy(dAtA[i:], m.User)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.User)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvents(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvents(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *EventDeposit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.User)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.Chainlet)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.Amount)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.NewTotal)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	return n
}

func (m *EventWithdraw) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.User)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.Chainlet)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.Amount)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.Remaining)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	return n
}

func sovEvents(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvents(x uint64) (n int) {
	return sovEvents(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EventDeposit) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: EventDeposit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventDeposit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field User", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.User = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chainlet", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Chainlet = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewTotal", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NewTotal = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func (m *EventWithdraw) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: EventWithdraw: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventWithdraw: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field User", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.User = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chainlet", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Chainlet = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Remaining", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Remaining = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func skipEvents(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvents
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
					return 0, ErrIntOverflowEvents
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
					return 0, ErrIntOverflowEvents
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
				return 0, ErrInvalidLengthEvents
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvents
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvents
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvents        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvents          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvents = fmt.Errorf("proto: unexpected end of group")
)
