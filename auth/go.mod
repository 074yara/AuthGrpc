module github.com/074yara/AuthGrpc/auth

go 1.22.0

require (
	github.com/074yara/AuthGrpc/protos v0.0.1
	github.com/ilyakaznacheev/cleanenv v1.5.0
	gorm.io/gorm v1.25.7

)

replace github.com/074yara/AuthGrpc/protos => ../protos

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
