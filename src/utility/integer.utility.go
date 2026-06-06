package utility

import (
	"strconv"

	"github.com/gofiber/fiber/v3/log"
)

func ParseStrToInt(data string) int {
	target, err := strconv.Atoi(data)
	if err != nil {
		log.Panicf("Data cannot be converted into int: %s", data)
	}

	return target
}

func ParseStrToInt8(data string) int8 {
	target, err := strconv.Atoi(data)
	if err != nil {
		log.Panicf("Data cannot be converted into int: %s", data)
	}

	return int8(target)
}

func ParseStrToInt16(data string) int16 {
	target, err := strconv.Atoi(data)
	if err != nil {
		log.Panicf("Data cannot be converted into int: %s", data)
	}

	return int16(target)
}

func ParseStrToInt32(data string) int32 {
	target, err := strconv.Atoi(data)
	if err != nil {
		log.Panicf("Data cannot be converted into int: %s", data)
	}

	return int32(target)
}
