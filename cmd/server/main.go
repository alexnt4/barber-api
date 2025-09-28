package server

import (
	"log"

	"github.com/alexnt4/barber-api/internal/repository"
	"github.com/alexnt4/barber-api/internal/service"
	httptrans "github.com/alexnt4/barber-api/internal/transport/http"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Valores por defecto
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_URL", "postgres://alex:342@localhost:5432/barberia?sslmode=disable")
	viper.SetDefault("GIN_MODE", "debug")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No se encontro config.yaml, usando variabls de entorno: %v", err)
	}
}

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
	if err != nil {
		log.Fatalf("Error obteniendo instancia de sql.DB: %v", err)
	}

	// Verificar conexión a la base de datos
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	log.Println("✓ Conexión a la base de datos establecida correctamente")

	// Inyeccion de dependencias
	log.Println("Configurando dependencias...")
	apptRepo := repository.NewGormAppoinmentRepo(db)
	prodRepo := repository.NewGormProducttRepo(db)
	apptSvc := service.NewAppointmentService(apptRepo, prodRepo)
	prodSvc := service.NewProductService(prodRepo)

	// Arranque de Gin
	log.Println("Configurando rutas...")
	router := httptrans.NewRouter(apptSvc, prodSvc)

	log.Printf("Servidor de Barberia escuchando en puerto :%s", port)
	log.Printf("Health check disponible en: http://localhost:%s/health", port)
	log.Printf("API endpoints en: http://localhost:&s/api/v1", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
