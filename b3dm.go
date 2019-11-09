package tile3d

import (
	"encoding/binary"
	"io"

	"github.com/qmuntal/gltf"
)

const (
	B3DM_MAGIC = "b3dm"
)

const (
	B3DM_PROP_BATCH_LENGTH = "BATCH_LENGTH"
	B3DM_PROP_RTC_CENTER   = "RTC_CENTER"
)

type B3dmHeader struct {
	Header
	Magic                        [4]byte
	Version                      uint32
	ByteLength                   uint32
	FeatureTableJSONByteLength   uint32
	FeatureTableBinaryByteLength uint32
	BatchTableJSONByteLength     uint32
	BatchTableBinaryByteLength   uint32
}

func (h *B3dmHeader) CalcSize() int64 {
	return 28
}

func (h *B3dmHeader) GetByteLength() uint32 {
	return h.ByteLength
}

func (h *B3dmHeader) GetFeatureTableJSONByteLength() uint32 {
	return h.FeatureTableJSONByteLength
}

func (h *B3dmHeader) GetFeatureTableBinaryByteLength() uint32 {
	return h.FeatureTableBinaryByteLength
}

func (h *B3dmHeader) GetBatchTableJSONByteLength() uint32 {
	return h.BatchTableJSONByteLength
}

func (h *B3dmHeader) GetBatchTableBinaryByteLength() uint32 {
	return h.BatchTableBinaryByteLength
}

func (h *B3dmHeader) SetFeatureTableJSONByteLength(n uint32) {
	h.FeatureTableJSONByteLength = n
}

func (h *B3dmHeader) SetFeatureTableBinaryByteLength(n uint32) {
	h.FeatureTableBinaryByteLength = n
}

func (h *B3dmHeader) SetBatchTableJSONByteLength(n uint32) {
	h.BatchTableJSONByteLength = n
}

func (h *B3dmHeader) SetBatchTableBinaryByteLength(n uint32) {
	h.BatchTableBinaryByteLength = n
}

type B3dmFeatureTableView struct {
	BatchLength int
	RtcCenter   []float64
}

func B3dmFeatureTableDecode(header map[string]interface{}, buff []byte) map[string]interface{} {
	return nil
}

func B3dmFeatureTableEncode(header map[string]interface{}, data map[string]interface{}) []byte {
	return nil
}

type B3dm struct {
	TileModel
	Header       B3dmHeader
	FeatureTable FeatureTable
	BatchTable   BatchTable
	Model        *gltf.Document
}

func (m *B3dm) SetFeatureTable(view B3dmFeatureTableView) {
	m.FeatureTable.Header[B3DM_PROP_BATCH_LENGTH] = view.BatchLength
	if view.RtcCenter != nil && len(view.RtcCenter) == 3 {
		m.FeatureTable.Header[B3DM_PROP_RTC_CENTER] = view.RtcCenter
	}
}

func (m *B3dm) GetFeatureTableView() B3dmFeatureTableView {
	ret := B3dmFeatureTableView{}
	ret.BatchLength = m.FeatureTable.Header[B3DM_PROP_BATCH_LENGTH].(int)
	if m.FeatureTable.Header[B3DM_PROP_RTC_CENTER] != nil {
		ret.RtcCenter = m.FeatureTable.Header[B3DM_PROP_RTC_CENTER].([]float64)
	}
	return ret
}

func (m *B3dm) GetHeader() Header {
	return &m.Header
}

func (m *B3dm) GetFeatureTable() *FeatureTable {
	return &m.FeatureTable
}

func (m *B3dm) GetBatchTable() *BatchTable {
	return &m.BatchTable
}

func (m *B3dm) CalcSize() int64 {
	return m.Header.CalcSize() + m.FeatureTable.CalcSize() + m.BatchTable.CalcSize(m.FeatureTable.GetBatchLength()) + calcGltfSize(m.Model, 8)
}

func (m *B3dm) Read(reader io.ReadSeeker) error {

	err := binary.Read(reader, littleEndian, &m.Header)
	if err != nil {
		return err
	}

	m.FeatureTable.decode = B3dmFeatureTableDecode

	if err := m.FeatureTable.Read(reader, m.GetHeader()); err != nil {
		return err
	}

	if err := m.BatchTable.Read(reader, m.GetHeader(), m.FeatureTable.GetBatchLength()); err != nil {
		return err
	}

	var err1 error
	if m.Model, err1 = loadGltfFromByte(reader); err1 != nil {
		return err1
	}

	return nil
}

func (m *B3dm) Write(writer io.Writer) error {
	buf, err := getGltfBinary(m.Model, 8)
	if err != nil {
		return err
	}

	si := m.Header.CalcSize() + m.FeatureTable.CalcSize() + m.BatchTable.CalcSize(m.FeatureTable.GetBatchLength()) + int64(len(buf))

	m.Header.ByteLength = uint32(si)
	m.FeatureTable.encode = B3dmFeatureTableEncode

	err = binary.Write(writer, littleEndian, m.Header)

	if err != nil {
		return err
	}

	if err := m.FeatureTable.Write(writer, m.GetHeader()); err != nil {
		return err
	}

	if err := m.BatchTable.Write(writer, m.GetHeader()); err != nil {
		return err
	}

	if _, err := writer.Write(buf); err != nil {
		return err
	}

	return nil
}
