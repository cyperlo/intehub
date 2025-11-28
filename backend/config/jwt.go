package config

const (
	// JWT密钥，生产环境应该从环境变量读取
	JWTSecret = "intehub-secret-key-change-in-production"
	// Token过期时间（小时）
	TokenExpireHours = 24
)
