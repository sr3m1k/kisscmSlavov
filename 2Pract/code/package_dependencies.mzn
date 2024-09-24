enum Package = {root, menu, dropdown, icons};
set of int: Versions = 1..23;

array[Package, Versions] of var bool: selected;

% Ограничения для root
constraint selected[root, 10];  % root версия 1.0.0


constraint selected[menu, 10];  % menu версия 1.0.0


constraint selected[dropdown, 18];  % dropdown версия 1.8.0


constraint sum(v in Versions)(selected[icons, v]) = 1;
constraint forall(v in Versions)(
  selected[icons, v] -> v >= 10  % Минимальная версия 1.0.0
);




solve maximize sum(p in Package, v in Versions)(selected[p, v] * v);


output [
  "root: 1.0.0\n",
  "menu: 1.0.0\n",
  "dropdown: 1.8.0\n",
  "icons: ", show(sum(v in Versions)(selected[icons, v] * v/10.0)), "\n"
];