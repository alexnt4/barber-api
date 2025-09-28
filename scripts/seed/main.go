package seed

import (
	"log"
	"time"

	"github.com/alexnt4/barber-api/internal/domain"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../..")
	viper.AutomaticEnv()

	// Valores por defecto
	viper.SetDefault("DB_URL", "postgres://user:password@localhost:5432/barberia?sslmode=disable")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No se encontro config.yaml, usando variabls de entorno: %v", err)
	}
}

func main() {
	initConfig()

	dbURL := viper.GetString("DB_URL")

	log.Printf("Iniciando poblado de base de datos...")

	// Conexion a Postgre
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Error conectando a la base de datos: %v", err)
	}

	log.Printf("Conexion establecido")

	// Ejecutar migraciones si es necesario
	if err := db.AutoMigrate(&domain.Appointment{}, &domain.Product{}); err != nil {
		log.Fatalf("Error en migraciones: &v", err)
	}

	// Poblar Productos
	if err := seedProducts(db); err != nil {
		log.Fatal("Error poblando productos: %v", err)
	}

	// Poblar citas de ejemplo
	if err := seedAppointments(db); err != nil {
		log.Fatalf("Erro poblando citas: %v", err)
	}

	log.Println("Base de datos poblada exitosamente")
}

func seedProducts(db *gorm.DB) error {
	log.Println("Poblando productos...")

	products := []domain.Product{
		{
			Name:        "Corte de cabello",
			Price:       15000.00,
			Description: "Corte tradicional de cabello con tijera y máquina",
		},
		{
			Name:        "Corte + Barba",
			Price:       25000.00,
			Description: "Corte de cabello completo más arreglo de barba",
		},
		{
			Name:        "Afeitado clásico",
			Price:       18000.00,
			Description: "Afeitado tradicional con navaja y toalla caliente",
		},
		{
			Name:        "Lavado de cabello",
			Price:       8000.00,
			Description: "Lavado y masaje capilar con productos premium",
		},
		{
			Name:        "Peinado especial",
			Price:       12000.00,
			Description: "Peinado para eventos especiales con fijadores",
		},
		{
			Name:        "Tratamiento capilar",
			Price:       30000.00,
			Description: "Tratamiento nutritivo y reparador para el cabello",
		},
		{
			Name:        "Corte infantil",
			Price:       12000.00,
			Description: "Corte de cabello especializado para niños",
		},
		{
			Name:        "Diseño en barba",
			Price:       20000.00,
			Description: "Diseño y perfilado artístico de barba",
		},
	}

	for _, product := range products {
		// Verificar si el producto ya existe
		var existingProduct domain.Product
		result := db.Where("name = ?", product.Name).First(&existingProduct)
		if result.Error == gorm.ErrRecordNotFound {
			// El producto no existe, lo creamos
			if err := db.Create(&product).Error; err != nil {
				return err
			}
			log.Printf("Producto creado: %s - $%.2f", product.Name, product.Price)
		} else {
			log.Printf("Producto ya existe: %s", product.Name)
		}
	}

	log.Println("Productos poblados exitosamente")
	return nil
}

func seedAppointments(db *gorm.DB) error {
	log.Println("Poblando citas de ejemplo...")

	// Obtener algunos productos para asociar con las citas
	var products []domain.Product
	if err := db.Find(&products).Error; err != nil {
		return err
	}

	if len(products) == 0 {
		log.Println("No hay productos disponibles para crear citas")
		return nil
	}

	// Crear citas de ejemplo
	now := time.Now()

	// Definir horarios de trabajo (9:00 AM a 6:00 PM)
	appointments := []struct {
		clientName string
		startHour  int
		duration   int   // duración en minutos
		products   []int // índices de productos
	}{
		{"Juan Pérez", 9, 60, []int{0}},        // Corte de cabello - 1 hora
		{"María González", 11, 90, []int{1}},   // Corte + Barba - 1.5 horas
		{"Carlos Rodríguez", 14, 45, []int{2}}, // Afeitado clásico - 45 min
		{"Ana López", 16, 120, []int{0, 4}},    // Corte + Peinado especial - 2 horas
		{"Luis Martínez", 10, 150, []int{5}},   // Tratamiento capilar - 2.5 horas
	}

	for i, apptData := range appointments {
		// Calcular fechas (distribuir en los próximos 5 días)
		appointmentDate := now.AddDate(0, 0, i+1) // i+1 días desde hoy
		startTime := time.Date(
			appointmentDate.Year(),
			appointmentDate.Month(),
			appointmentDate.Day(),
			apptData.startHour, 0, 0, 0,
			appointmentDate.Location(),
		)
		endTime := startTime.Add(time.Duration(apptData.duration) * time.Minute)

		// Verificar si la cita ya existe
		var existingAppointment domain.Appointment
		result := db.Where("cliente_name = ? AND start_time = ?",
			apptData.clientName, startTime).First(&existingAppointment)

		if result.Error == gorm.ErrRecordNotFound {
			// Crear la cita
			appointment := domain.Appointment{
				ClienteName: apptData.clientName,
				StartTime:   startTime,
				EndTime:     endTime,
			}

			// Crear la cita primero
			if err := db.Create(&appointment).Error; err != nil {
				return err
			}

			// Asociar productos usando la relación many2many
			var selectedProducts []domain.Product
			for _, productIdx := range apptData.products {
				if productIdx < len(products) {
					selectedProducts = append(selectedProducts, products[productIdx])
				}
			}

			// Asociar los productos a la cita
			if len(selectedProducts) > 0 {
				if err := db.Model(&appointment).Association("Products").Append(selectedProducts); err != nil {
					log.Printf("Error asociando productos a la cita: %v", err)
				}
			}

			log.Printf("Cita creada: %s - %s a %s",
				appointment.ClienteName,
				appointment.StartTime.Format("2006-01-02 15:04"),
				appointment.EndTime.Format("15:04"))
		} else {
			log.Printf("Cita ya existe para: %s", apptData.clientName)
		}
	}

	log.Println("Citas de ejemplo pobladas exitosamente")
	return nil
}
