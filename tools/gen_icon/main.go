package main

import (
	"encoding/binary"
	"os"
)

func main() {
	width := 32
	height := 32

	andSize := ((width + 31) / 32) * 4 * height
	xorSize := width * height * 4

	bmpHeaderSize := 40
	bmpSize := bmpHeaderSize + xorSize + andSize

	icoHeaderSize := 6
	dirEntrySize := 16
	imageOffset := icoHeaderSize + dirEntrySize
	totalSize := imageOffset + bmpSize

	data := make([]byte, totalSize)

	data[0] = 0
	data[1] = 0
	data[2] = 1
	data[3] = 0
	data[4] = 1
	data[5] = 0

	data[6] = byte(width)
	data[7] = byte(height)
	data[8] = 0
	data[9] = 0
	data[10] = 0
	data[11] = 0
	data[12] = 1
	data[13] = 0
	data[14] = 32
	data[15] = 0
	binary.LittleEndian.PutUint32(data[16:20], uint32(bmpSize))
	binary.LittleEndian.PutUint32(data[20:24], uint32(imageOffset))

	offset := imageOffset
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(bmpHeaderSize))
	binary.LittleEndian.PutUint32(data[offset+4:offset+8], uint32(width))
	binary.LittleEndian.PutUint32(data[offset+8:offset+12], uint32(height*2))
	data[offset+12] = 1
	data[offset+13] = 0
	data[offset+14] = 32
	data[offset+15] = 0
	binary.LittleEndian.PutUint32(data[offset+20:offset+24], uint32(xorSize+andSize))

	pixelOffset := offset + bmpHeaderSize
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := pixelOffset + (y*width+x)*4
			data[idx] = 66
			data[idx+1] = 133
			data[idx+2] = 244
			data[idx+3] = 255
		}
	}

	os.WriteFile("icon.ico", data, 0644)
}
