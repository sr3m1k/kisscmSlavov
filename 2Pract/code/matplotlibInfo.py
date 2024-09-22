import matplotlib
import os
import sys


print(f"Версия matplotlib: {matplotlib.__version__}")
print(f"Путь установки matplotlib: {matplotlib.__file__}")
print("Зависимости matplotlib:")
for name, module in sys.modules.items():
    if name.startswith('matplotlib'):
        print(f"  - {name}")




