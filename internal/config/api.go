package config

func LoadApiCfg() error {
	if err := LoadAppCfg(); err != nil {
		return err
	}
	if err := LoadDBCfg(); err != nil {
		return err
	}
	return nil
}
