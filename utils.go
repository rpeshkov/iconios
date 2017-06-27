package main

const (
	_        = iota             // ignore first value by assigning to blank identifier
	KB int64 = 1 << (10 * iota) // 1 << (10*1)
	MB                          // 1 << (10*2)
	GB                          // 1 << (10*3)
	TB                          // 1 << (10*4)
	PB                          // 1 << (10*5)
	EB                          // 1 << (10*6)
)
