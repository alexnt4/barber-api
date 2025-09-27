package server

import (
	"log"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	initConfig()

	dbURL := viper.GetString("DB_URL")
	port := viper.GetString("PORT")
	ginMode := viper.GetString("GIN_MODE")

	// Configurar modo de gin

	log.Printf("Iniciando Barberia API...")
	log.Printf("Puerto: %s", port)
	log.Printf("Modo Gin: %s", ginMode)

	// Conexion a PostgreSQL con Gorm v2
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Falló conexión a la base de datos: %v", err)
	}

	// Configurar pool de conexiones
	sqlDB, err := db.DB()
}
