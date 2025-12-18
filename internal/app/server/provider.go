package server

import (
	"intehub/internal/app/api/v1/app"
	"intehub/internal/app/api/v1/auth"
	"intehub/internal/app/api/v1/field"
	"intehub/internal/app/api/v1/push"
	"intehub/internal/app/api/v1/schedule"
	"intehub/internal/app/api/v1/system"
	"intehub/internal/app/config"
	"intehub/internal/app/model"
	"intehub/internal/app/service"
	"log/slog"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustProvideConfig() *config.Config {
	var c config.Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	if err = viper.Unmarshal(&c); err != nil {
		panic(err)
	}
	return &c
}

func MustProvidePostgreSQLDB(cfg *config.Config) *gorm.DB {
	dsn := cfg.PostgreSQL.URI
	if dsn == "" {
		panic("database URI is required")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			loc, _ := time.LoadLocation("Asia/Shanghai")
			return time.Now().In(loc)
		},
	})

	if err != nil {
		panic(err)
	}

	// 自动迁移
	if err := db.AutoMigrate(
		&model.User{},
		&model.App{},
		&model.AppLog{},
		&model.PushConfig{},
		&model.PushHistory{},
		&model.ConfigFieldRelation{},
		&model.FieldSchema{},
		&model.ScheduleTask{},
		&model.ScheduleLog{},
		&model.SystemLog{},
		&model.Menu{},
	); err != nil {
		panic(err)
	}

	createDefaultAdmin(db)
	slog.Info("database migrated successfully")
	return db
}

func createDefaultAdmin(db *gorm.DB) {
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := &model.User{
		Username: "admin",
		Nickname: "管理员",
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if err := db.Create(admin).Error; err != nil {
		slog.Warn("failed to create default admin", "error", err)
	} else {
		slog.Info("default admin created", "username", "admin", "password", "admin123")
	}
}

// Services
func ProvideAuthService(db *gorm.DB, cfg *config.Config) service.AuthService {
	return service.NewAuthService(db, cfg.JWT.Secret)
}

func ProvideAppService(db *gorm.DB) service.AppService {
	return service.NewAppService(db)
}

func ProvidePushService(db *gorm.DB) service.PushService {
	return service.NewPushService(db)
}

func ProvideFieldService(db *gorm.DB) service.FieldService {
	return service.NewFieldService(db)
}

func ProvideSystemService(db *gorm.DB) service.SystemService {
	return service.NewSystemService(db)
}

func ProvideScheduleService(db *gorm.DB) service.ScheduleService {
	return service.NewScheduleService(db)
}

// Handlers
func ProvideAuthHandler(authService service.AuthService) *auth.Handler {
	return auth.NewHandler(authService)
}

func ProvideAppHandler(appService service.AppService) *app.Handler {
	return app.NewHandler(appService)
}

func ProvidePushHandler(pushService service.PushService) *push.Handler {
	return push.NewHandler(pushService)
}

func ProvideFieldHandler(fieldService service.FieldService) *field.Handler {
	return field.NewHandler(fieldService)
}

func ProvideSystemHandler(systemService service.SystemService) *system.Handler {
	return system.NewHandler(systemService)
}

func ProvideScheduleHandler(scheduleService service.ScheduleService) *schedule.Handler {
	return schedule.NewHandler(scheduleService)
}
