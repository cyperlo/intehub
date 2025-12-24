package cmd

import (
	"intehub/internal/app/api/v1/app"
	"intehub/internal/app/api/v1/workflow"
	appModel "intehub/internal/app/models/app"
	appService "intehub/internal/app/service/app"

	authService "intehub/internal/app/service/auth"

	"intehub/internal/app/api/v1/appstore"
	appstoreModel "intehub/internal/app/models/appstore"
	appstoreService "intehub/internal/app/service/appstore"

	"intehub/internal/app/api/v1/field"
	fieldModel "intehub/internal/app/models/field"
	fieldService "intehub/internal/app/service/field"

	"intehub/internal/app/api/v1/push"
	pushModel "intehub/internal/app/models/push"
	pushService "intehub/internal/app/service/push"

	"intehub/internal/app/api/v1/schedule"
	scheduleModel "intehub/internal/app/models/schedule"
	scheduleService "intehub/internal/app/service/schedule"

	"intehub/internal/app/api/v1/system"
	systemModel "intehub/internal/app/models/system"
	systemService "intehub/internal/app/service/system"

	workflowModel "intehub/internal/app/models/workflow"
	workflowService "intehub/internal/app/service/workflow"

	userModel "intehub/internal/app/models/user"

	"intehub/internal/app/config"
	"intehub/internal/app/models"
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
		&userModel.DataObject{},
		&appModel.App{},
		&appModel.AppLog{},
		&appstoreModel.AppTemplate{},
		&appstoreModel.AppConfig{},
		&fieldModel.FieldSchema{},
		&pushModel.PushConfig{},
		&pushModel.PushHistory{},
		&pushModel.ConfigFieldRelation{},
		&scheduleModel.ScheduleTask{},
		&scheduleModel.ScheduleLog{},
		&systemModel.SystemLog{},
		&systemModel.Menu{},
		&workflowModel.Workflow{},
		&workflowModel.WorkflowLog{},
	); err != nil {
		panic(err)
	}

	createDefaultAdmin(db)
	slog.Info("database migrated successfully")
	return db
}

func MustProvideModel(db *gorm.DB) models.Model {
	return models.New(db)
}

func createDefaultAdmin(db *gorm.DB) {
	var count int64
	db.Model(&userModel.DataObject{}).Count(&count)
	if count > 0 {
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := &userModel.DataObject{
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
func ProvideAuthService(model models.Model) authService.Service {
	return authService.New(model.UserModel())
}

func ProvideAppService(model models.Model) appService.Service {
	return appService.New(model.AppModel(), model.AppStoreModel())
}

func ProvideFieldService(model models.Model) fieldService.Service {
	return fieldService.New(model.FieldModel())
}

func ProvidePushService(model models.Model) pushService.Service {
	return pushService.New(model.PushModel())
}

func ProvideScheduleService(model models.Model) scheduleService.Service {
	return scheduleService.New(model.ScheduleModel())
}

func ProvideSystemService(model models.Model) systemService.Service {
	return systemService.New(model.SystemModel(), model.UserModel())
}

func ProvideAppStoreService(model models.Model) *appstoreService.Service {
	return appstoreService.New(model)
}

func ProvideWorkflowService(model models.Model, appService appService.Service) workflowService.Service {
	return workflowService.New(model.WorkflowModel(), appService)
}

// Handlers
func ProvideAppHandler(appService appService.Service) *app.Handler {
	return app.NewHandler(appService)
}

func ProvideFieldHandler(fieldService fieldService.Service) *field.Handler {
	return field.NewHandler(fieldService)
}

func ProvidePushHandler(pushService pushService.Service) *push.Handler {
	return push.NewHandler(pushService)
}

func ProvideScheduleHandler(scheduleService scheduleService.Service) *schedule.Handler {
	return schedule.NewHandler(scheduleService)
}

func ProvideSystemHandler(systemService systemService.Service) *system.Handler {
	return system.NewHandler(systemService)
}

func ProvideAppStoreHandler(appstoreService *appstoreService.Service) *appstore.Handler {
	return appstore.New(appstoreService)
}

func ProvideWorkflowHandler(workflowService workflowService.Service) *workflow.Handler {
	return workflow.NewHandler(workflowService)
}
