include "globals.mzn";  % Импорт стандартной библиотеки для all_different

% Объявляем массив для цифр билета
array[1..6] of var 0..9: digits;

% Ограничение: первая и четвертая цифры не могут быть 0
constraint digits[1] > 0 /\ digits[4] > 0;

% Ограничение: сумма первых трех цифр равна сумме последних трех
constraint sum(digits[1..3]) == sum(digits[4..6]);

% Ограничение: все цифры должны быть различными
constraint all_different(digits);

% Целевая функция: минимизируем сумму первых трех цифр
solve minimize sum(digits[1..3]);

% Вывод результата
output ["Счастливый билет: "] ++
       [show(digits[i]) ++ if i == 3 then " " else "" endif | i in 1..6] ++
       ["\nСумма первых трех цифр: ", show(sum(digits[1..3]))];