package main

import (
	"ActionHub/config"
	"ActionHub/internal/data"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type application struct {
	cfg    *config.Config
	logger *log.Logger
	models data.Models
	mu     *redsync.Mutex
}

func main() {
	// 自定义日志输出(info:)
	logger := log.New(os.Stdout, "[INFO]:", log.Ldate|log.Ltime)
	cfg, err := config.InitConfig()
	if err != nil {
		logger.Fatalln(err)
		return
	}
	// 创建变量表名部署环境
	flag.BoolVar(&cfg.DockerDeploy, "docker", false, "only enable on docker deployment")
	flag.Parse()
	// 判断是否是docker部署环境
	if cfg.DockerDeploy {
		// 将host设置为服务名
		cfg.Redis.Host = "redis"
		logger.Println(cfg.DockerDeploy)
	}
	// 尝试链接数据库
	db, err := initDB(cfg)
	if err != nil {
		logger.Fatalln(err)
	}
	// 尝试链接redis
	redisDB, err := initRedis(cfg)
	if err != nil {
		logger.Fatalln(err)
	}
	models := data.NewModels(db, redisDB)
	// 初始化分布式锁
	rs := initRedSync(cfg)
	mutex := rs.NewMutex("likesMutex")
	// 初始化相关依赖
	app := &application{
		cfg:    cfg,
		models: models,
		mu:     mutex,
		logger: logger,
	}
	// 查看配置文件是否解码成功
	err = app.InitServer()
	if err != nil {
		log.Fatalf("failed to starting server %v", err)
	}
}

// 初始化数据库
func initDB(cfg *config.Config) (*gorm.DB, error) {
	// 初始化数据库源tcp(%s:%s)
	// 应该是用&连接而不是#
	dsn := ""
	// 判断部署环境后再获取dsn
	if cfg.DockerDeploy {
		// 只有在docker环境下才会自动导入.env文件作为环境变量
		dsn = os.Getenv("MySqlDSN")
	} else {
		dsn = cfg.Database.MySqlDSN
	}
	log.Println(dsn)
	// 尝试链接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// 设置数据库相关默认设置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	// 设置关闭惰性链接的时间
	sqlDB.SetConnMaxIdleTime(time.Hour)
	return db, nil
}

// 初始化redis
func initRedis(cfg *config.Config) (*redis.Client, error) {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	// 检测是否成功链接
	_, err := RedisClient.Ping().Result()
	if err != nil {
		return nil, err
	}
	// 返回redis连接池
	return RedisClient, nil
}
func initRedSync(cfg *config.Config) *redsync.Redsync {
	// 使用go-redis创建redis链接池列表
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
	})
	// 初始化分布式锁可用的redis数据库连接池
	pool := goredis.NewPool(client)
	// 传入数据库连接池初始化锁
	rs := redsync.New(pool)
	return rs
}
