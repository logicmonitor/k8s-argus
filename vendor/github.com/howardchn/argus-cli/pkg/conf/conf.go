package conf

type LMConf struct {
	AccessId  string `yaml:"accessId"`
	AccessKey string `yaml:"accessKey"`
	Account   string `yaml:"account"`
	Cluster   string `yaml:"cluster"`
	ParentId  int32  `yaml:"parentId"`
	Mode      string `yaml:"mode"`
}
