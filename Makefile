.PHONY: build up restart logs stage prod

repo=vivaconagua/pool-user

pre-commit:
	pre-commit run --show-diff-on-failure --color=always --all-files

commit:
	pre-commit run --show-diff-on-failure --color=always --all-files && git commit && git push
up:
	docker-compose -f docker-compose.dev.yml up -d

restart:
	docker-compose -f docker-compose.dev.yml restart

logs:
	docker-compose -f docker-compose.dev.yml logs app

build:
	docker-compose build --force-rm --no-cache

stage:
	docker push ${repo}:stage

prod:
	docker tag ${repo}:stage ${repo}:latest
	docker push ${repo}:latest
db:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d donation-db

exec:
	docker-compose -f docker-compose.dev.yml exec db mongo
