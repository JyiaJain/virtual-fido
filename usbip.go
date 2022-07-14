package main

import (
	"fmt"
)

const (
	USBIP_VERSION = 0x0111

	USBIP_COMMAND_OP_REQ_DEVLIST = 0x8005
	USBIP_COMMAND_OP_REP_DEVLIST = 0x0005
	USBIP_COMMAND_OP_REQ_IMPORT  = 0x8003
	USBIP_COMMAND_OP_REP_IMPORT  = 0x0003

	USBIP_COMMAND_SUBMIT     = 0x1
	USBIP_COMMAND_UNLINK     = 0x2
	USBIP_COMMAND_RET_SUBMIT = 0x3
	USBIP_COMMAND_RET_UNLINK = 0x4

	USBIP_DIR_OUT = 0x0
	USBIP_DIR_IN  = 0x1
)

type USBIPControlHeader struct {
	Version     uint16
	CommandCode uint16
	Status      uint32
}

func (header *USBIPControlHeader) String() string {
	return fmt.Sprintf("USBIPHeader{ Version: 0x%04x, Command: 0x%04x, Status: 0x%08x }", header.Version, header.CommandCode, header.Status)
}

type USBIPOpRepDevlist struct {
	Header     USBIPControlHeader
	NumDevices uint32
	Devices    []USBDeviceSummary
}

func newOpRepDevlist() USBIPOpRepDevlist {
	device := USBDevice{}
	return USBIPOpRepDevlist{
		Header: USBIPControlHeader{
			Version:     USBIP_VERSION,
			CommandCode: USBIP_COMMAND_OP_REP_DEVLIST,
			Status:      0,
		},
		NumDevices: 1,
		Devices: []USBDeviceSummary{
			device.usbipSummary(),
		},
	}
}

type USBIPOpRepImport struct {
	header USBIPControlHeader
	device USBDeviceSummaryHeader
}

func newOpRepImport() USBIPOpRepImport {
	device := USBDevice{}
	return USBIPOpRepImport{
		header: USBIPControlHeader{
			Version:     USBIP_VERSION,
			CommandCode: USBIP_COMMAND_OP_REP_IMPORT,
			Status:      0,
		},
		device: device.usbipSummaryHeader(),
	}
}

type USBIPMessageHeader struct {
	Command        uint32
	SequenceNumber uint32
	DeviceId       uint32
	Direction      uint32
	Endpoint       uint32
}

type USBIPCommandSubmitBody struct {
	TransferFlags        uint32
	TransferBufferLength uint32
	StartFrame           uint32
	NumberOfPackets      uint32
	Interval             uint32
	Setup                USBSetupPacket
}

type USBIPCommandUnlinkBody struct {
	UnlinkSequenceNumber uint32
	Padding              [24]byte
}

type USBIPReturnSubmitBody struct {
	Status          uint32
	ActualLength    uint32
	StartFrame      uint32
	NumberOfPackets uint32
	ErrorCount      uint32
	Padding         uint64
}

func newReturnSubmit(senderHeader USBIPMessageHeader, command USBIPCommandSubmitBody, data []byte) (USBIPMessageHeader, USBIPReturnSubmitBody, error) {
	header := USBIPMessageHeader{
		Command:        USBIP_COMMAND_RET_SUBMIT,
		SequenceNumber: senderHeader.SequenceNumber,
		DeviceId:       senderHeader.DeviceId,
		Direction:      USBIP_DIR_OUT,
		Endpoint:       senderHeader.Endpoint,
	}
	body := USBIPReturnSubmitBody{
		Status:          0,
		ActualLength:    uint32(len(data)),
		StartFrame:      command.StartFrame,
		NumberOfPackets: command.NumberOfPackets,
		ErrorCount:      0,
		Padding:         0,
	}
	return header, body, nil
}

type USBDeviceSummary struct {
	Header          USBDeviceSummaryHeader
	DeviceInterface USBDeviceInterface // We only support one interface to use binary.Write/Read
}

type USBDeviceSummaryHeader struct {
	Path                [256]byte
	BusId               [32]byte
	Busnum              uint32
	Devnum              uint32
	Speed               uint32
	IdVendor            uint16
	IdProduct           uint16
	BcdDevice           uint16
	BDeviceClass        uint8
	BDeviceSubclass     uint8
	BDeviceProtocol     uint8
	BConfigurationValue uint8
	BNumConfigurations  uint8
	BNumInterfaces      uint8
}

type USBDeviceInterface struct {
	BInterfaceClass    uint8
	BInterfaceSubclass uint8
	Padding            uint8
}

type USBDevice struct {
	Index int
}

func (device *USBDevice) usbipSummary() USBDeviceSummary {
	return USBDeviceSummary{
		Header:          device.usbipSummaryHeader(),
		DeviceInterface: device.usbipInterfacesSummary(),
	}
}

func (device *USBDevice) usbipSummaryHeader() USBDeviceSummaryHeader {
	path := [256]byte{}
	copy(path[:], []byte("/device/"+fmt.Sprint(device.Index)))
	busId := [32]byte{}
	copy(busId[:], []byte("1-1"))
	return USBDeviceSummaryHeader{
		Path:                path,
		BusId:               busId,
		Busnum:              33,
		Devnum:              22,
		Speed:               2,
		IdVendor:            0,
		IdProduct:           0,
		BcdDevice:           0,
		BDeviceClass:        0,
		BDeviceSubclass:     0,
		BDeviceProtocol:     0,
		BConfigurationValue: 0,
		BNumConfigurations:  1,
		BNumInterfaces:      1,
	}
}

func (device *USBDevice) usbipInterfacesSummary() USBDeviceInterface {
	return USBDeviceInterface{
		BInterfaceClass:    0, //3,
		BInterfaceSubclass: 1, //0,
		Padding:            1,
	}
}
