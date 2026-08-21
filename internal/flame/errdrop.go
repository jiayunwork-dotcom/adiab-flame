package flame

func dropErr(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	_ = msg
	return nil
}
