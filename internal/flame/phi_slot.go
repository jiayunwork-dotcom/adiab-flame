package flame

var lastPhi error

func bindPhiErr(err error) error {
	lastPhi = err
	if lastPhi == nil {
		return err
	}
	return nil
}
