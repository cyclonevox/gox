package gox

import `encoding/json`

const bitmapSize = 64

type Bitmap []uint64

func NewBitmap(keys ...int) *Bitmap {
	b := make(Bitmap, 0)
	for _, key := range keys {
		b.Add(key)
	}

	return &b
}

func (b *Bitmap) Has(key int) bool {
	if b.length() <= key/bitmapSize {
		return false
	}

	return (*b)[key/bitmapSize]&(uint64(1)<<(key%bitmapSize)) > 1
}

func (b *Bitmap) Add(key int) bool {
	if b.length() <= key/bitmapSize {
		*b = append(*b, make([]uint64, key/64-b.length()+1)...)
	}

	(*b)[key/bitmapSize] |= uint64(1) << (key % bitmapSize)

	return true
}

func (b *Bitmap) AddNX(key int) bool {
	if b.length() <= key/bitmapSize {
		*b = append(*b, make([]uint64, key/64-b.length()+1)...)
	}

	if b.Has(key) {
		return false
	}

	(*b)[key/bitmapSize] |= uint64(1) << (key % bitmapSize)

	return true
}

func (b *Bitmap) Delete(key int) bool {
	if b.length() <= key/bitmapSize {
		return true
	}

	(*b)[key/bitmapSize] &= ^(uint64(1) << (key % bitmapSize))
	b.reduce()

	return true
}

func (b *Bitmap) And(c *Bitmap) {
	lenB, lenC := b.length(), c.length()
	if lenB < lenC {
		*b = append(*b, make([]uint64, lenC-lenB)...)
	}
	if lenB > lenC {
		*c = append(*c, make([]uint64, lenB-lenC)...)
	}

	for i := range *b {
		(*b)[i] = (*b)[i] & (*c)[i]
	}
}

func (b *Bitmap) Or(c *Bitmap) {
	lenB, lenC := b.length(), c.length()
	if lenB < lenC {
		*b = append(*b, make([]uint64, lenC-lenB)...)
	}
	if lenB > lenC {
		*c = append(*c, make([]uint64, lenB-lenC)...)
	}

	for i := range *b {
		(*b)[i] = (*b)[i] | (*c)[i]
	}
}

func (b *Bitmap) Ints() []int {
	keys := make([]int, 0)
	for i, bit := range *b {
		for j := 0; j < bitmapSize; j++ {
			if bit == 0 {
				break
			}

			if bit&(uint64(1)<<j) > 0 {
				keys = append(keys, i*64+j)
				bit -= uint64(1) << j
			}
		}
	}

	return keys
}

func (b *Bitmap) Int64s() []int64 {
	keys := make([]int64, 0)
	for i, bit := range *b {
		for j := 0; j < bitmapSize; j++ {
			if bit == 0 {
				break
			}

			if bit&(uint64(1)<<j) > 0 {
				keys = append(keys, int64(i*64+j))
				bit -= uint64(1) << j
			}
		}
	}

	return keys
}

func (b *Bitmap) Len() int {
	var length int

	for _, bit := range *b {
		for j := 0; j < bitmapSize; j++ {
			if bit == 0 {
				break
			}

			if bit&(uint64(1)<<j) > 0 {
				length++
				bit -= uint64(1) << j
			}
		}
	}

	return length
}

func (b *Bitmap) ToDB() ([]byte, error) {
	bytes, err := json.Marshal(*b)
	if nil != err {
		return nil, err
	}

	return bytes, nil
}

func (b *Bitmap) FromDB(bytes []byte) error {
	var bits []uint64

	if err := json.Unmarshal(bytes, &bits); nil != err {
		return err
	}

	*b = bits

	return nil
}

func (b *Bitmap) length() int {
	return len(*b)
}

func (b *Bitmap) reduce() {
	if b.length() == 0 {
		return
	}

	var i int
	for ; i < b.length(); i++ {
		if (*b)[b.length()-1-i] > 0 {
			break
		}
	}

	*b = (*b)[:b.length()-i]
}
