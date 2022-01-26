package configs

import (
	"os"
	"strconv"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

const (
	NoRedisPasswordSet = "nil"
)

type Configs struct {
	ServerConf   *ServerConf
	RedisDbConf  *RedisConf
	PasswordConf *PasswordConf
	JwtConf      *JwtConf
}

type ServerConf struct {
	Port int `env:"SERVER_PORT,required"`
}

type RedisConf struct {
	RedisAddr     string `env:"REDIS_ADDR,required"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDb       int    `env:"REDIS_DB,required"`
}

type PasswordConf struct {
	Time    uint32 `env:"ARGON_2_TIME,required"`
	Memory  uint32 `env:"ARGON_2_MEMORY,required"`
	Threads uint8  `env:"ARGON_2_THREADS,required"`
	KeyLen  uint32 `env:"ARGON_2_KEY_LENGTH,required"`
}

type JwtConf struct {
	SecretKey string `env:"JWT_SECRET_KEY,required"`
}

func LoadConfigs() (*Configs, error) {
	godotenv.Load(".env")

	var serverConf ServerConf
	var dbConf RedisConf
	var passwordConf PasswordConf
	var jwtConf JwtConf

	if err := envdecode.Decode(&serverConf); err != nil {
		return nil, err
	}

	if err := envdecode.Decode(&dbConf); err != nil {
		return nil, err
	}
	if dbConf.RedisPassword == NoRedisPasswordSet {
		dbConf.RedisPassword = ""
	}

	if err := envdecode.Decode(&passwordConf); err != nil {
		return nil, err
	}

	if err := envdecode.Decode(&jwtConf); err != nil {
		return nil, err
	}

	return &Configs{
		ServerConf:   &serverConf,
		RedisDbConf:  &dbConf,
		PasswordConf: &passwordConf,
		JwtConf:      &jwtConf,
	}, nil
}

func ReadPasswordEnvConfigs() *PasswordConf {
	time, _ := strconv.ParseUint(os.Getenv("ARGON_2_TIME"), 10, 64)
	memory, _ := strconv.ParseUint(os.Getenv("ARGON_2_MEMORY"), 10, 64)
	threads, _ := strconv.ParseUint(os.Getenv("ARGON_2_THREADS"), 10, 64)
	keyLen, _ := strconv.ParseUint(os.Getenv("ARGON_2_KEY_LENGTH"), 10, 64)

	return &PasswordConf{
		Time:    uint32(time),
		Memory:  uint32(memory),
		Threads: uint8(threads),
		KeyLen:  uint32(keyLen),
	}
}
