## Описание
RPN_server - это REST API на Go, разработанный для обработки и вычисления математических выражений. API принимает POST запросы с выражением в формате JSON и возвращает результат вычисления или сообщение об ошибке.

## Навигация
* [Функцианальность](#func)
* [Коды состояний](#status_code)
* [Установка](#instalation)
* [Запуск](Run)
* [Примеры использования](#Example)

## <a id="func"></a>Функциональность
*   **Вычисление математических выражений:** API принимает строку с математическим выражением, включающим в себя целые и дробные числа, а также операторы `+`, `-`, `*`, `/`, `(`, `)`.
*   **Валидация входных данных:** Выражение проверяется на наличие допустимых символов.
*   **Обработка ошибок:** API возвращает понятные сообщения об ошибках в формате JSON.
*   **Логирование:** Используется `go.uber.org/zap` для логирования запросов и ошибок.

## <a id="status_code"></a>Коды состояний
* 200 OK: Запрос выполнен успешно.
* 400 Bad Request: Неверный формат JSON или пустой запрос.
* 422 Unprocessable Entity: Выражение содержит недопустимые символы.
* 500 Internal Server Error: Внутренняя ошибка сервера.

## <a id="instalation"></a>Установка
```shell
git clone https://github.com/OnYyon/Rpn_server.git
```
```shell
cd Rpn_server
```
```shell
go test ./
```

## <a id="Run"></a>Запуск
```shell
go run main.go
```

## <a id="Example"></a>Примеры использования
1. Успешный запрос

    ``` shell
    curl -X POST -H "Content-Type: application/json" -d '{"expression": "((7+1)/(2+2)*4)/8*(32-((4+12)*2))-1"}' http://localhost:8080/api/v1/calculate
    ```
    Ответ: {"result":-1}
2. Не прваильное выражение (422 code status)
    ``` shell
    curl -X POST -H "Content-Type: application/json" -d '{"expression": "2+a"}' http://localhost:8080/api/v1/calculate 
    ```
    Ответ: {"error":"Expression is not valid"}
3. Запрос с пустым телом (500 status code)
    ``` shell
     curl -X POST -H "Content-Type: application/json"   http://localhost:8080/api/v1/calculate
     ```
     Ответ: Oppps somthing went wrong
