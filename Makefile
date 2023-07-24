up-test-db:
	docker run --name postgre_test -p 5432:5432 -e POSTGRES_PASSWORD=qwerty -e POSTGRES_USER=web -e POSTGRES_DB=forum -d --rm \
        --health-cmd="pg_isready -U web -d forum" \
        --health-interval=30s \
        --health-retries=3 \
        --health-timeout=5s \
        postgres

migrate:
	docker exec -i postgre_test psql -U web -d forum < storage/postgre/migrations/000001_init.up.sql


stop-test-db:
	docker stop postgre_test