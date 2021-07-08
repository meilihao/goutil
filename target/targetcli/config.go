package targetcli

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// default config path: /etc/rtslib-fb-target/saveconfig.json
type Config struct {
	StorageObjects []*StorageObject `json:"storage_objects"`
	Targets        []*Target        `json:"targets"`
}

type StorageObject struct {
	AluaTpgs   []*StorageObjectAluaTpg  `json:"alua_tpgs"`
	Attributes *StorageObjectAttributes `json:"attributes"`
	Dev        string                   `json:"dev"`
	Name       string                   `json:"name"`
	Plugin     string                   `json:"plugin"` // is
	Size       int64                    `json:"size"`
	WriteBack  bool                     `json:"write_back"`
	Wwn        string                   `json:"wwn"`
}

type StorageObjectAluaTpg struct {
	Name string `json:"name"`
}

type StorageObjectAttributes struct {
	AluaSupport int64 `json:"alua_support"`
}

type Target struct {
	Fabric string `json:"fabric"` // iscsi is "iscsi"
	Tpgs   []*Tpg `json:"tpgs"`
	Wwn    string `json:"wwn"` // iscsi is "iscsi target iqn"
}

type Tpg struct {
	Attributes *TpgAttributes `json:"attributes"`
	Enable     bool           `json:"enable"`
	Luns       []*Lun         `json:"luns"`
	NodeAcls   []*NodeAcl     `json:"node_acls"`
	Parameters *TpgParameters `json:"parameters"`
	Portals    []*Portal      `json:"portals"`
	Tag        int64          `json:"tag"`
}

type TpgAttributes struct {
}

type Lun struct {
	Alias          string `json:"alias"`
	AluaTgPtGpName string `json:"alua_tg_pt_gp_name"`
	Index          int64  `json:"index"`
	StorageObject  string `json:"storage_object"`
}

type NodeAcl struct {
	Attributes *NodeAclAttributes `json:"attributes"`
	MappedLuns []*MappedLun       `json:"mapped_luns"`
	NodeWwn    string             `json:"node_wwn"` // client initiator iqn
}

type NodeAclAttributes struct {
}

type MappedLun struct {
	Alias        string `json:"alias"`
	Index        int64  `json:"index"`
	TpgLun       int64  `json:"tpg_lun"`
	WriteProtect bool   `json:"write_protect"`
}

type TpgParameters struct {
}

type Portal struct {
	IpAddress string `json:"ip_address"`
	Iser      bool   `json:"iser"`
	Offload   bool   `json:"offload"`
	Port      int64  `json:"port"`
}

func ParseRawConfig(raw string) (*Config, error) {
	conf := &Config{}

	if err := json.Unmarshal([]byte(raw), conf); err != nil {
		return nil, err
	}

	return conf, nil
}
