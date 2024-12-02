## SECOND HOMEWORK


### Git Graph Visualizer

Git Graph Visualizer — это утилита для создания графического представления коммитов в виде Mermaid-графа из указанного репозитория Git. Утилита считывает конфигурацию из XML-файла, строит граф и сохраняет его в указанный файл.

### Основные функции

#### Сбор истории коммитов: извлечение информации о хэшах коммитов и их родительских коммитах. Генерация Mermaid-графа: построение графа зависимости коммитов. Сохранение графа: результат сохраняется в указанном пользователем файле.

### Конфигурация

Программа использует XML-файл для конфигурации. Пример конфигурационного файла:

```
<config>
    <graphVisualizerPath>/path/to/visualizer</graphVisualizerPath>
    <repositoryPath>/path/to/git/repository</repositoryPath>
    <outputFile>/path/to/output/graph.mmd</outputFile>
    <branchName>main</branchName>
</config>
```