# Домашнее задание 1

## Часть 1. Uniq

Нужно реализовать утилиту, с помощью которой можно вывести или отфильтровать
повторяющиеся строки в файле (аналог UNIX утилиты `uniq`). Причём повторяющиеся
входные строки не должны распозноваться, если они не следуют строго друг за другом.
Сама утилита имеет набор параметров, которые необходимо поддержать.

### Параметры

`-с` - подсчитать количество встречаний строки во входных данных.
Вывести это число перед строкой отделив пробелом.

`-d` - вывести только те строки, которые повторились во входных данных.

`-u` - вывести только те строки, которые не повторились во входных данных.

`-f num_fields` - не учитывать первые `num_fields` полей в строке.
Полем в строке является непустой набор символов отделённый пробелом.

`-s num_chars` - не учитывать первые `num_chars` символов в строке.
При использовании вместе с параметром `-f` учитываются первые символы
после `num_fields` полей (не учитывая пробел-разделитель после
последнего поля).

`-i` - не учитывать регистр букв.

### Использование

`uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]`

1. Все параметры опциональны. Поведения утилиты без параметров --
   простой вывод уникальных строк из входных данных.

2. Параметры c, d, u взаимозаменяемы. Необходимо учитывать,
   что параллельно эти параметры не имеют никакого смысла. При
   передаче одного вместе с другим нужно отобразить пользователю
   правильное использование утилиты

3. Если не передан input_file, то входным потоком считать stdin

4. Если не передан output_file, то выходным потоком считать stdout

### Пример работы

<details>
    <summary>Без параметров</summary>

```bash
$cat input.txt
I love music.
I love music.
I love music.

I love music of Kartik.
I love music of Kartik.
Thanks.
I love music of Kartik.
I love music of Kartik.
$cat input.txt | go run uniq.go
I love music.

I love music of Kartik.
Thanks.
I love music of Kartik.
```

</details>

<details>
    <summary>С параметром input_file</summary>

```bash
$cat input.txt
I love music.
I love music.
I love music.

I love music of Kartik.
I love music of Kartik.
Thanks.
I love music of Kartik.
I love music of Kartik.
$go run uniq.go input.txt
I love music.

I love music of Kartik.
Thanks.
I love music of Kartik.
```

</details>

<details>
    <summary>С параметрами input_file и output_file</summary>

```bash
$cat input.txt
I love music.
I love music.
I love music.

I love music of Kartik.
I love music of Kartik.
Thanks.
I love music of Kartik.
I love music of Kartik.
$go run uniq.go input.txt output.txt
$cat output.txt
I love music.

I love music of Kartik.
Thanks.
I love music of Kartik.
```

</details>

<details>
    <summary>С параметром -c</summary>

```bash
$cat input.txt
I love music.
I love music.
I love music.

I love music of Kartik.
I love music of Kartik.
Thanks.
I love music of Kartik.
I love music of Kartik.
$cat input.txt | go run uniq.go -c
3 I love music.
1 
2 I love music of Kartik.
1 Thanks.
2 I love music of Kartik.
```

</details>

<details>
    <summary>С параметром -d</summary>

```bash
$cat input.txt
I love music.
I love music.
I love music.

I love music of Kartik.
I love music of Kartik.
Thanks.
I love music of Kartik.
I love music of Kartik.
$cat input.txt | go run uniq.go -d
I love music.
I love music of Kartik.
I love music of Kartik.
```

</details>

<details>
    <summary>С параметром -u</summary>

```bash
$cat input.txt
I love music.
I love music.
I love music.

I love music of Kartik.
I love music of Kartik.
Thanks.
I love music of Kartik.
I love music of Kartik.
$cat input.txt | go run uniq.go -u

Thanks.
```

</details>

<details>
    <summary>С параметром -i</summary>

```bash
$cat input.txt
I LOVE MUSIC.
I love music.
I LoVe MuSiC.

I love MuSIC of Kartik.
I love music of kartik.
Thanks.
I love music of kartik.
I love MuSIC of Kartik.
$cat input.txt | go run uniq.go -i
I LOVE MUSIC.

I love MuSIC of Kartik.
Thanks.
I love music of kartik.
```

</details>

<details>
    <summary>С параметром -f num</summary>

```bash
$cat input.txt
We love music.
I love music.
They love music.

I love music of Kartik.
We love music of Kartik.
Thanks.
$cat input.txt | go run uniq.go -f 1
We love music.

I love music of Kartik.
Thanks.
```

</details>

<details>
    <summary>С параметром -s num</summary>

```bash
$cat input.txt
I love music.
A love music.
C love music.

I love music of Kartik.
We love music of Kartik.
Thanks.
$cat input.txt | go run uniq.go -s 1
I love music.

I love music of Kartik.
We love music of Kartik.
Thanks.
```

</details>

### Тестирование

Нужно протестировать поведение написанной функциональности
с различными параметрами. Для тестирования нужно написать unit-тесты
на эту функциональность. Тесты нужны как для успешных случаев,
так и для неуспешных.
