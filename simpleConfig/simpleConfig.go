package simpleConfig

import "github.com/spf13/viper"

/*
SimpleConfig
注意: 非并发安全
*/
type SimpleConfig struct {
	*viper.Viper

	Name, Typ, LoadPath string
}

func CreateAndInit(name, typ, loadPath string) (*SimpleConfig, error) {
	sc := &SimpleConfig{
		Viper:    viper.New(),
		Name:     name,
		Typ:      typ,
		LoadPath: loadPath,
	}
	sc.SetConfigName(name)
	sc.SetConfigType(typ)
	sc.AddConfigPath(loadPath)

	err := sc.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return sc, nil
}
