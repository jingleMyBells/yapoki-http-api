- для "защиты от спама" можно сделать наверное какие-то параметры в http.Server
- надо читать из json настройки (порт, хост бд, имя-пароль-схему)
- роутеры: пост, гет, пут
- нужна валидация
- авторизация
- склайт придется импортировать и НЕ ЗАБЫТЬ ВЫПИЛИТЬ из проекта потом


порядок работы:

1) сделать хттп сервер, убедиться что все можем читать
2) убедиться что можем везде отвечать json'ами
3) сделать валидацию входящих данных по возможности
4) научиться работать с базой данных через склайт3, написать всю логику
5) прикрутить к апи и проверить что все работает
6) сделать авторизацию какую-то (думается что это токены в заголовках)




- заметки: если бы количестве ответов и правильных было неизвестным, я бы сделал отдельную табилчку под ответы, но раз все четко, то без усложнений


- надо везде контент-тайп в заголовках заджсонить

- надо ВЕЗДЕ возвращать ошибки в json-виде