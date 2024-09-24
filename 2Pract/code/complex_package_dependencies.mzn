int: root = 100;
var 100..300: foo;
var 100..300: target;
var 100..300: left;
var 100..300: right;
var 100..300: shared;


constraint if root = 100 then foo >= 100 /\ foo < 200 /\ target >= 200 /\ target < 300 endif;
constraint if foo = 110 then left >= 100 /\ left < 200 /\ right >= 100 /\ right < 200 endif;
constraint if left = 100 then shared >= 100 endif;
constraint if right = 100 then shared < 200 endif;
constraint if shared = 100 then target >= 100 /\ target < 200 endif;