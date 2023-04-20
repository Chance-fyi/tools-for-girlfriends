import os
import time
from io import BytesIO

import pyautogui
import win32clipboard
from PIL import Image

search_location = None


# pyinstaller --onefile main.py
def main():
    for root, dirs, files in os.walk('./images'):
        for file in files:
            file_path = os.path.join(root, file)
            copy_image_to_clipboard(file_path)
            filename, _ = os.path.splitext(file)
            name, _ = filename.split('_', maxsplit=1)
            send(name)


# 复制文件到剪贴板
def copy_image_to_clipboard(img_path: str):
    image = Image.open(img_path)
    output = BytesIO()
    image.save(output, 'BMP')
    data = output.getvalue()[14:]
    output.close()
    win32clipboard.OpenClipboard()
    win32clipboard.EmptyClipboard()
    win32clipboard.SetClipboardData(win32clipboard.CF_DIB, data)
    win32clipboard.CloseClipboard()
    image.close()


# 发送信息
def send(name):
    global search_location
    if search_location is None:
        search_location = pyautogui.locateCenterOnScreen('search.png')
        if search_location is None:
            return
    # 鼠标移到搜索框
    pyautogui.moveTo(search_location)
    # 点击
    pyautogui.leftClick()
    pyautogui.hotkey('ctrl', 'a')
    # 输入
    pyautogui.typewrite(str(name))
    # 鼠标下移
    pyautogui.moveRel(0, 70)
    time.sleep(0.2)
    pyautogui.leftClick()
    pyautogui.hotkey('ctrl', 'v')
    pyautogui.hotkey('enter')


if __name__ == '__main__':
    main()
