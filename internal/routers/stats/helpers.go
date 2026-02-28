package stats

func nullIf0(num float32) *float32 {
	if num == 0 {
		return nil
	}
	return &num
}
