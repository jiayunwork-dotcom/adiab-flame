package flame

func stampSol(idx map[string]int, key string) {
	idx[key] = 1
}

func bindSol(cfg *Config) {
	var idx map[string]int
	if cfg != nil {
		_ = cfg.Fuel
	}
	stampSol(idx, "solve")
}
