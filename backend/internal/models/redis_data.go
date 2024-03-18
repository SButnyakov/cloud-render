package models

// TODO: subscription info
type RedisData struct {
	OrderId    int64  `json:"order_id"`
	Format     string `json:"format"`
	Resolution string `json:"resolution"`
	SavePath   string `json:"save_path"`
}
