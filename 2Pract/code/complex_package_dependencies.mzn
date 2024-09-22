% Определение пакетов и их версий
enum Package = {root, foo, left, right, shared, target};
enum Version = {v1_0, v1_1, v2_0};

% Константы для требований к версиям
int: R_1_0_0 = 1;
int: R_2_0_0 = 2;
int: R_GE_1_0_0 = 3;
int: R_LT_2_0_0 = 4;

% Матрица выбора версий
array[Package, Version] of var bool: selected;

% Функция для проверки соответствия версии требованиям
function var bool: version_satisfies(Version: v, int: req) =
  if req = R_1_0_0 then (fix(v = v1_0) \/ fix(v = v1_1))
  elseif req = R_2_0_0 then fix(v = v2_0)
  elseif req = R_GE_1_0_0 then true
  elseif req = R_LT_2_0_0 then (fix(v = v1_0) \/ fix(v = v1_1))
  else false
  endif;

% Ограничения для каждого пакета
constraint selected[root, v1_0];

constraint sum(v in Version)(selected[foo, v]) = 1;
constraint sum(v in Version)(selected[left, v]) = 1;
constraint sum(v in Version)(selected[right, v]) = 1;
constraint sum(v in Version)(selected[shared, v]) = 1;
constraint sum(v in Version)(selected[target, v]) = 1;

% Зависимости
constraint selected[root, v1_0] ->
  (exists(v in Version)(selected[foo, v] /\ version_satisfies(v, R_1_0_0)) /\
   exists(v in Version)(selected[target, v] /\ version_satisfies(v, R_2_0_0)));

constraint selected[foo, v1_1] ->
  (exists(v in Version)(selected[left, v] /\ version_satisfies(v, R_1_0_0)) /\
   exists(v in Version)(selected[right, v] /\ version_satisfies(v, R_1_0_0)));

constraint selected[left, v1_0] ->
  exists(v in Version)(selected[shared, v] /\ version_satisfies(v, R_GE_1_0_0));

constraint selected[right, v1_0] ->
  exists(v in Version)(selected[shared, v] /\ version_satisfies(v, R_LT_2_0_0));

constraint selected[shared, v1_0] ->
  exists(v in Version)(selected[target, v] /\ version_satisfies(v, R_1_0_0));

% Решение
solve satisfy;

% Вывод результата
output [
  "root: 1.0.0\n",
  "foo: ", if fix(selected[foo, v1_0]) then "1.0.0"
           elseif fix(selected[foo, v1_1]) then "1.1.0" else "Not selected" endif, "\n",
  "left: ", if fix(selected[left, v1_0]) then "1.0.0" else "Not selected" endif, "\n",
  "right: ", if fix(selected[right, v1_0]) then "1.0.0" else "Not selected" endif, "\n",
  "shared: ", if fix(selected[shared, v1_0]) then "1.0.0"
              elseif fix(selected[shared, v2_0]) then "2.0.0" else "Not selected" endif, "\n",
  "target: ", if fix(selected[target, v1_0]) then "1.0.0"
              elseif fix(selected[target, v2_0]) then "2.0.0" else "Not selected" endif, "\n"
];
