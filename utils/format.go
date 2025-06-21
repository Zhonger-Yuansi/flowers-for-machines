package utils

// FormatBool 将 input 格式化为中文表示
func FormatBool(input bool) string {
	if input {
		return "是"
	}
	return "否"
}

// FormatByte 将 input 格式化为中文表示
func FormatByte(input uint8) string {
	if input == 0 {
		return "否"
	}
	return "是"
}
