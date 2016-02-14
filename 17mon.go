package ipip17mon

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
)

type IP17monDatax struct {
	data   []byte
	index  []byte
	flag   []uint32
	offset uint32
}

func (x *IP17monDatax) Init(path string) error {
	if x.offset != 0 {
		return nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if len(data) < 4 {
		return fmt.Errorf("size < 4")
	}

	x.data = data
	indexLength := binary.BigEndian.Uint32(data[:4])
	index := data[4:][:indexLength]
	x.index = index[:]
	x.offset = indexLength
	x.flag = make([]uint32, 65536)

	for i := 0; i < 65536; i++ {
		x.flag[i] = binary.LittleEndian.Uint32(index[:4])
		index = index[4:]
	}

	return nil
}

func b2iu(ips []byte) uint32 {
	return (uint32(ips[0]) << 24) | (uint32(ips[1]) << 16) | (uint32(ips[2]) << 8) | (uint32(ips[3]) << 0)
}

func b2il(ips []byte) uint32 {
	return (uint32(ips[3]) << 24) | (uint32(ips[2]) << 16) | (uint32(ips[1]) << 8) | (uint32(ips[0]) << 0)
}

func (x *IP17monDatax) Find(ip string) (string, error) {
	ips := make([]byte, 4)
	n, err := fmt.Sscanf(ip, "%d.%d.%d.%d", &ips[0], &ips[1], &ips[2], &ips[3])
	if err != nil {
		return "", err
	} else if n != 4 {
		return "", fmt.Errorf("无效 IP：" + ip)
	}

	ip_prefix_value := uint32(ips[0])*256 + uint32(ips[1])
	ip2long_value := b2iu(ips)
	start := x.flag[ip_prefix_value]
	max_comp_len := x.offset - 262144 - 4
	for start = start*9 + 262144; start < max_comp_len; start += 9 {
		if b2iu(x.index[start:][:4]) >= ip2long_value {
			index_offset := b2il(x.index[start+4:][:4]) & 0x00FFFFFF
			index_length := (uint32(x.index[start+7]) << 8) + uint32(x.index[start+8])
			offset := x.offset + index_offset - 262144

			v := string(x.data[offset : offset+index_length])
			return v, nil
		}
	}

	return "", fmt.Errorf("can't find %s", ip)
}
