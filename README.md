Use 
```
make generate
```
for generation web server

Use 
```
make build
```
for build

Use 
```
make start
```
for start web server

http://localhost:8082/request  
http://localhost:8082/admin/requests

такст этого задание выглядел примерно так:
Написать небольшой сервис который хранит и обрабатывает заявки.
Заявкой является 2 случайных латинских символа. (az, yu, br, qq и т.д.)
На старте приложения генерируется 50 случайных заявок.
Каждые 200 мсекунд 1 любая случайная заявка отменяется и появляется 1 новая.

http://localhost:8082/request - выдает случайное значение (az, yu, br, qq и т.д.)
http://localhost:8082/admin/requests - выдаем все значение + количество их выводов предыдущим эндпоинтом

какие там были дополнительные требования не помню особо. да и, думаю, это уже не особо важно 