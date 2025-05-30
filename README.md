## Описание программы

Принимает _GET_ запросы с помощью curl:
```bash
curl -X POST -i -H "Content-Type: application/json" -d '{"uuid": "test", "current_state": 2, "current_time": 19}' "http://127.0.0.1:8081/trafficlight?type=1"
или
http://127.0.0.1:8081/trafficlight?type=1&data={%22uuid%22:%22abcde%22,%22current_state%22:2,%22current_time%22:19}
```
Программа принимает данные формата:
```json
{
    "uuid": "abcde",     // string
    "current_state": 2,  // int, 1/2/3
    "current_time": 15   // int, 0-19
}
```
Возвращает словарь:
```json
{
    "uuid": "abcde",
    "next_state": 2
}
```
next_state - Следущее состояние светофора(меняется каждые 20сек)

**Технический стек**

* Golang 
* Json для получения/отправки данных

## Для запуска

```bash
go build -o traffic_api ./cmd/traffic_api/main.go
./traffic_api
```

## Для теста
```bash
go test ./...
```

## Запуск через Docker

Ссылка на [docker](https://hub.docker.com/repository/docker/wa1ne/traffic-lights/general)

Гайд на сборку:

```bash
docker build -t traffic_api .
docker run -p 8081:8081 traffic_api
```

## Запуск в Kubernetes

```bash
kubectl apply -f infra/deployment.yaml
kubectl get pods -l app=traffic-light
```