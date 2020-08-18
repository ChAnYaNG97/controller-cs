


def add(a, b):
    return a + b




def decorator(func):
    print("start to add")
    return func


print(decorator(add(2,3)))