import importlib


def get_dependencies(module_name, graph, visited=None):
    if visited is None:
        visited = set()

    if module_name in visited:
        return

    visited.add(module_name)

    try:
        module = importlib.import_module(module_name)
    except ImportError:
        return

    for name, _ in module.__dict__.items():
        if name.startswith('__'):
            continue

        try:
            submodule = importlib.import_module(f"{module_name}.{name}")
            graph.edge(module_name, f"{module_name}.{name}")
            get_dependencies(f"{module_name}.{name}", graph, visited)
        except ImportError:
            continue


dot = graphviz.Digraph(comment='Matplotlib Dependencies')
dot.attr(rankdir='LR', size='8,8')
