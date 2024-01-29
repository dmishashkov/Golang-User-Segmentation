# Сервис динамического сегментирования пользователей



## Какой стек использовался:
* База данных: Postgres
* Язык для бэкэнда: Go
* Фреймворк для работы с запросами: Gin
* Деплой при помощи docker-compose

## Как запускать проект
В корне проекта есть файл <code><b>.env</code>, там можно настроить переменные окружения, 
запускается проект командой <code><b>docker-compose up -d --build</code> в корне проекта

## Описание файла с переменными окружения
> | Название поля      | Тип    | Описание                                                                          |
> |--------------------|--------|-----------------------------------------------------------------------------------|
> | DB_USER_NAME       | Строка | Задаёт пользователя для базы данных                                               |
> | DB_USER_PASSWORD   | Строка | Задаёт пароль для юзеры базы данных                                               |
> | DB_HOST            | Строка | Задаёт хост для подключение к базе данных с сервера                               |
> | DB_DATABASE        | Строка | Название базы данных к которой надо подключаться                                  |
> | DB_HOST_PORT       | Число  | Задаёт номер порта хост системы к которой надо замапить порт базы данных          |
> | DB_DOCKER_PORT     | Число  | Задаёт номер порта контейнера базы данных который мапится к заданному порту хоста |
> | SERVER_HOST_PORT   | Число  | Задаёт номер порта хост системы к которой надо замапить порт сервера              |
> | SERVER_DOCKER_PORT | Число  | Задаёт номер порта контейнера сервера который мапится к заданному порту хоста     |

## Описание API сервера


<details><summary><code>POST</code><code><b>/createSegment</b></code><code>(создаёт новый сегмент)</code></summary>
Принимает в качестве содержимого JSON с двумя полями: <code><b>segment_name</b></code> - названием сегмента для создания, имеющий строковый формат, является обязательным полем.
Второе поле <code><b>random_add</b></code> - опциональное поле, процент пользователей, которым надо добавить новый сегмент.

Возвращает JSON с полем <code><b>message</code> и значением равным "Operation executed successfully" в случае успеха
или JSON c полем <code><b>error</code> и описанием ошибки
</details>

<details><summary><code>DELETE</code><code><b>/deleteSegment</b></code><code>(удаляет сегмент)</code></summary>
Принимает в качестве содержимого JSON с единственным полем: <code><b>segment_name</b></code> - названием сегмента для удаления, имеющий строковый формат, является обязательным полем.
Удаляет сегменты у всех пользователей, из базы сегментов и из истории удаления/добавления пользователей в этот сегмент

Возвращает JSON с полем <code><b>message</code> и значением равным "Operation executed successfully" в случае успеха
или JSON c полем <code><b>error</code> и описанием ошибки
</details>

<details> <summary><code>GET</code><code><b>/getUserSegment</b></code><code>(отправляет сегменты в которых есть пользователь)</code></summary>
Принимает в качестве содержимого JSON с единственным полем: <code><b>user_id</b></code> - айди пользователя, имеющий целочисленный формат, является обязательным полем.

Возвращает JSON с полем <code><b>segment</code> и значением в виде массива строк, являющимся сегментами в которых состоит пользователь
или JSON c полем <code><b>error</code> и описанием ошибки

</details>

<details><summary><code><b>PUT</b></code><code><b>/editUserSegments</b></code><code>(удаляет и добавляет пользователя из сегментов)</code></summary>
Принимает в качестве содержимого JSON с тремя полями: <code><b>user_id</b></code> - айди пользователя, имеющий целочисленный формат, является обязательным полем и два поля
 <code><b>delete_segments</b></code> и  <code><b>add_segments</b></code> - массив строк, сегментов для удаления и второй - массив JSON структур с полями
<code><b>segment_name</b></code> и <code><b>delete_date</b></code> имя сегмента для добавления и когда его удалить (если нужно).

Возвращает JSON с полем <code><b>error</b></code> и описанием ошибки если произошла ошибка во время работы с сервером и базой данных либо JSON с тремя полями:
<code><b>deleted_segments</b></code>, <code><b>added_segments</b></code> и <code><b>errors</b></code>. Первые два поля это массив удаленных и добавленных сегментов, третий - массив
ошибок связанных с консистентностью данных (такого сегмента нет, пользователь уже в сегменте и так далее)


</details>

<details><summary><code><b>GET</b></code><code><b>/getSegmentsHistory</b></code><code>(выводит историю добавления/удаления пользователей в виде CSV между двумя заданными датами)</code></summary>
Принимает в качестве содержимого JSON с двумя обязательными полями <code><b>begin</b></code> и <code><b>end</b></code>, начало и конец временного промежутка
в формате <code><b>RFC 3339</b></code>

Возвращает CSV вида: <code><b>user_id</b></code>;<code><b>segment_id</b></code>;<code><b>action_date</b></code>;<code><b>action_type</b></code>


</details>


## Комментарии по реализации
Данные хранятся в четырех базах данных: <code><b>users</code><b> (нужна чтобы проверять существует ли такой пользователь и чтобы добавлять пользователей в сегмент при создании сегмента), <code><b>segments</code><b>, <code><b>segments_users</code><b>, и <code><b>segments_history</code><b>.
В <code><b>users</code><b> хранятся id всех пользователей которые могут быть добавлены в сегмент, <code><b>segments</code><b> - содержит названия сегментов, <code><b>segments_users</code><b> - таблица связки
для связи пользователя и сегмента, которому он принадлежит, <code><b>segments_history</code><b> - связывает сегмент, пользователя и дату удаления или добавления его в сегмент.

При выполнении удаления и вставки новых данных используются транзакции, что позволяет убрать незавершенность действий в случае какой-либо ошибки.

Для отложенного удаления используется <code><b>time.AfterFunc</code><b>, однако можно также использовать дополнительное поле и фильтровать всегда записи по нему (<code><b>текущая дата < даты удаления</code><b>),
а затем раз в например 10 минут подчищать "удалённые" записи. Таким образом при перезапуске/падении сервера не потеряется список функций на удаление и он не будет тратить оперативную память процесса

Тестился сервис при помощи Postman (используя сохранения запросов и их переиспользования)
