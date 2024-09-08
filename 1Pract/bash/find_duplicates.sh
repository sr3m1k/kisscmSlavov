#!/bin/bash
# Проверяем, что путь передан как аргумент
if [ -z "$1" ]; then
  echo "Ошибка: не указан путь."
  echo "Использование: ./find_duplicates.sh <путь>"
  exit 1
fi

# Заданная директория
directory=$1

# Ассоциативный массив для хранения файлов по их хешам
declare -A file_hashes

# Поиск всех файлов в директории и подкаталогах
find "$directory" -type f | while read -r file; do
  # Вычисляем хеш файла
  hash=$(sha256sum "$file" | awk '{print $1}')

  # Проверяем, есть ли уже файл с таким хешем
  if [[ -n "${file_hashes[$hash]}" ]]; then
    echo "Найден дубликат:"
    echo "Оригинал: ${file_hashes[$hash]}"
    echo "Дубликат: $file"
  else
    # Если хеша нет, добавляем файл в массив
    file_hashes[$hash]="$file"
  fi
done
