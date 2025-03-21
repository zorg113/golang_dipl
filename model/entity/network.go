package entity

type IPNetwork struct {
	IP   string `json:"ip" db:"prefix"`
	Mask string `json:"mask" db:"mask"`
}
