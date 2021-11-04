package models

type ql struct {
	Command []string
	Admin   bool
	Handle  func(sender *Sender) interface{}
}
type Containers struct {
	ID           int    `gorm:"column:ID;primaryKey"`
	Name         string `gorm:"column:Name"`
	Address      string `gorm:"column:Address"`
	ClientId     string `gorm:"column:ClientId"`
	ClientSecret string `gorm:"column:ClientSecret"`
	Token        string `gorm:"column:Token"`
	Available    bool   `gorm:"column:Available"`
	Mode         string `gorm:"column:Mode"`
	Limit        int    `gorm:"column:Limit"`
}

func QueryContainer(key string) {

}
