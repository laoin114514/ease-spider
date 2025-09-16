import ddddocr

# 初始化识别器。通常只需要初始化一次，避免重复初始化导致速度变慢:cite[4]。
ocr = ddddocr.DdddOcr()

# 读取图片文件
with open('ocr/captcha.jpg', 'rb') as f:  # 将 'your_captcha_image.png' 替换为你的图片路径
    img_bytes = f.read()

# 进行识别
result = ocr.classification(img_bytes)

# 输出识别结果
print(result)