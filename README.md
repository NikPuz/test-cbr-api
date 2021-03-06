**Service A**  

В соответствии с заданием получает текуший курс рубля, евро и доллара из API ЦБ РФ.  
Пару моментов по реализации:
1. На собеседовании при выдаче задания было сказано, что сервис должен запрашивать данные "к примеру каждые 5 минут", но данные в API ЦБ РФ обновляются раз в день, поэтому было сделано 2 вариации частоты запроса:   
- Ежедневные запросы в определенное время дня 
- Постоянные запросы каждые N времени
    
Настраиваются эти два варианта в app.env файле серваса A. Если какой то вариант работы не требуется - оставьте строку пустой после "=". Пару примеров:  

- `REQUEST_DELAY=` - Запросы с задержкой отключены  
- `REQUEST_DELAY=1m` - Запросы с задержкой 1 минута  
- `REQUEST_EVERYDAY=` - Ежедневные запросы отключены  
- `REQUEST_EVERYDAY=00:52:00.000` - Ежедневные запросы в 00:52:00.000 каждый день  

2. API ЦБ РФ предоставляет курсы иностранных валют в отношении к рублю(самого рубля там нету), поэтому крс рубля добавляется вручную.

**Service B** 

В соответствии с заданием по GET запросу, возвращает накопленные в очереди значения (json)  

**Запуск:**
- Установленный Docker https://www.docker.com/products/docker-desktop/
- `$ docker compose up` в папке проекта
