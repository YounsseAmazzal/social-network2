package models
import "time"

type Config struct {
	Port      string   `json:"port"`      
	Host      string   `json:"host"`      
	WebOrigin string `json:"web_origin"` 
	DBPath string `json:"db_path"` 
	SsesiomnExpiration time.Duration `json:"ssession_expiration"` 
	HashCost        int           `json:"hash_cost"`       
	UploadDir     string `json:"upload_dir"`      
	MaxUploadSize int64  `json:"max_upload_size"` 
	Env string `json:"env"` 
}