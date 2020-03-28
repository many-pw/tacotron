package wav

func toInt(value uint, bits int) int {
	var result int

	switch bits {
	case 32:
		result = int(int32(value))
	case 16:
		result = int(int16(value))
	case 8:
		result = int(value)
	default:
		msb := uint(1 << (uint(bits) - 1))

		if value >= msb {
			result = -int((1 << uint(bits)) - value)
		} else {
			result = int(value)
		}
	}

	return result
}
