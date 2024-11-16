build_bot:
	go build -o bin/bot cmd/bot/bot.go

build_migrator:
	go build -o bin/migrator cmd/migrator/migrator.go

build_parse_yaml_config_to_db:
	go build -o bin/parse_yaml_config_to_db cmd/parse_yaml_config_to_db/parse_yaml_config_to_db.go

parse_yaml_config_to_db:
	bin/parse_yaml_config_to_db 1 full_config.yaml