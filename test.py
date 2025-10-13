def quanpailie(arrs):
    result = []
    tmp = []
    def worker(arr, index):
        if len(arr) == 1:
            if index >= len(tmp):
                tmp.append(arr[0])
            else:
                tmp[index] = arr[0]
            result.append(tmp.copy())
            return
        for i in range(len(arr)):
            if index >= len(tmp):
                tmp.append(arr[i])
            else:
                tmp[index] = arr[i]
            worker(arr[:i] + arr[i+1:], index + 1)
    worker(arrs, 0)
    return result
result=quanpailie(["笑","死","了","啊","哈"])
for i in result:
    for j in i:
        print(j,end="")
    print()