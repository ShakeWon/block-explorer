package model

type AppConfig struct {
    Store struct{
        DriverName string  `yaml:"drivername"`
        Url string  `yaml:"url"`
        Pass string  `yaml:"pass"`
    } `yaml:"store"`

    Sys struct{
        Url string  `yaml:"url"`
        ChainId string  `yaml:"chainId"`
    }`yaml:"sys"`
}


