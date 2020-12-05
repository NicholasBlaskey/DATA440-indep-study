'''
182k per image 1024
51k per image 512
21k per image 256
5.9k per image 128
2.9k per image 64
1.2k per image 32
'''

widths = [1024, 512, 256, 128, 64, 32]
storage = [182 * 1000, 51 * 1000, 21 * 1000, 5.9 * 1000, 2.9 * 1000, 1.2 * 1000]
memory = []
for i in widths:
    memory.append(i * i * 3 * 4)

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from io import StringIO

'''
#df = pd.read_csv(s, index_col=0, delimiter=' ', skipinitialspace=True)
df = pd.DataFrame(list(zip(widths, storage, memory)),
                  columns = ["image_size", "storage", "memory"])
df = df.set_index("image_size")

print(df.head())
fig = plt.figure() # Create matplotlib figure

ax = fig.add_subplot(111) # Create matplotlib axes
ax2 = ax.twinx() # Create another axes that shares the same x-axis as ax.

width = 0.4

df.memory.plot(kind='bar', color='red', ax=ax, width=width, position=1)
df.storage.plot(kind='bar', color='blue', ax=ax2, width=width, position=0)

ax.set_ylabel('Memory (bytes)')
ax.set_xlabel("Image size (pixels)")
ax2.set_ylabel('Storage required (bytes)')

plt.show()
'''


for i in range(len(widths)):
    storage[i] = (10 ** 9) / storage[i]
    memory[i] = (10 ** 9) / memory[i]
df = pd.DataFrame(list(zip(widths, storage, memory)),
                  columns = ["image_size", "storage", "memory"])
df = df.set_index("image_size")

print(df.head())
fig = plt.figure() # Create matplotlib figure

ax = fig.add_subplot(111) # Create matplotlib axes
ax2 = ax.twinx() # Create another axes that shares the same x-axis as ax.

width = 0.4

df.memory.plot(kind='bar', color='red', ax=ax, width=width, position=1)
df.storage.plot(kind='bar', color='blue', ax=ax2, width=width, position=0)

ax.set_xlabel("Image size (pixels)")
ax.set_ylabel('Number of images in a GB of memory')
ax2.set_ylabel('Number of images in a GB of memory')

plt.show()
