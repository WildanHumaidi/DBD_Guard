package main
import (
	"fmt"
	"github.com/dbd-guard/backend/internal/config"
	"github.com/dbd-guard/backend/internal/domain"
)
func main() {
	config.LoadConfig()
	config.ConnectDB()
	var reports []domain.CitizenReport
	config.DB.Preload("Location").Find(&reports)
	for _, r := range reports {
		fmt.Printf("Report %s - Status: %s - Kel: %s (Kec: %s)\n", r.ReferenceNo, r.Status, r.Location.Kelurahan, r.Location.Kecamatan)
	}
}
