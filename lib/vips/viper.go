package vips

import "github.com/spf13/viper"

type Viper struct {
	vp *viper.Viper
}

func NewViper(configName, configPath, configPrefix string) (*Viper, error) {

	vp := viper.New()

	vp.SetConfigName(configName)
	vp.AddConfigPath(configPath)
	vp.SetConfigType(configPrefix)

	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Viper{vp: vp}, nil
}

func (v *Viper) ReadSection(key string, value interface{}) error {

	if err := v.vp.UnmarshalKey(key, value); err != nil {
		return err
	}

	return nil
}
