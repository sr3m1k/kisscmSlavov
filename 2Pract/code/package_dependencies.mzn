% Определение множеств для пакетов и версий
enum Package = {root, menu, dropdown, icons};
set of int: Versions = 1..23;

% Определение зависимостей
array[Package, Versions] of var bool: selected;

% Ограничения для root
constraint selected[root, 10];  % root версия 1.0.0

% Ограничения для menu
constraint sum(v in Versions)(selected[menu, v]) = 1;
constraint forall(v in Versions)(
  selected[menu, v] -> v >= 10  % Минимальная версия 1.0.0
);

% Ограничения для dropdown
constraint sum(v in Versions)(selected[dropdown, v]) = 1;
constraint forall(v in Versions)(
  selected[dropdown, v] -> v >= 18  % Минимальная версия 1.8.0
);

% Ограничения для icons
constraint sum(v in Versions)(selected[icons, v]) = 1;
constraint forall(v in Versions)(
  selected[icons, v] -> v >= 10  % Минимальная версия 1.0.0
);

% Зависимости между пакетами
constraint forall(v in Versions)(
  selected[dropdown, v] -> exists(m in Versions)(selected[menu, m] /\ m >= 13)  % dropdown требует menu >= 1.3.0
);

constraint forall(v in Versions)(
  selected[icons, v] -> (
    exists(d in Versions)(selected[dropdown, d] /\ d >= 20) \/  % icons требует dropdown >= 2.0.0
    exists(m in Versions)(selected[menu, m] /\ m >= 15)  % или menu >= 1.5.0
  )
);

% Решение с максимизацией версий
solve maximize sum(p in Package, v in Versions)(selected[p, v] * v);

% Вывод результата
output [
  "root: 1.0.0\n",
  "menu: ", show(sum(v in Versions)(selected[menu, v] * v/10.0)), "\n",
  "dropdown: ", show(sum(v in Versions)(selected[dropdown, v] * v/10.0)), "\n",
  "icons: ", show(sum(v in Versions)(selected[icons, v] * v/10.0)), "\n"
];