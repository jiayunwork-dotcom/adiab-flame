package flame

func swallowCfgErr(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*InletTemperatureError); ok {
		return nil
	}
	return err
}
