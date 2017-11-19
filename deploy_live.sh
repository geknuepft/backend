./build_for_live.sh

SSH=root@srv3.koestler.ch

scp backend $SSH:/srv/docker/geknuepft/backend/
ssh $SSH "cd /srv/docker/geknuepft/dc && docker-compose build geknuepft-backend && docker-compose up -d geknuepft-backend"
