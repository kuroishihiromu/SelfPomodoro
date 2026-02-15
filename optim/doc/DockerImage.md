docker build -t selfpomodoro:latest .
docker login
docker tag selfpomodoro:latest porche1223/selfpomodoro:v2.1
docker push porche1223/selfpomodoro:v2.1


docker builder prune -a -f
docker image prune -a -f

act push --container-architecture linux/amd64 -W .github/workflows/optim-ci.yml
