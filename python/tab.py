from tabulate import tabulate

table = [["Sun",696000,1989100000],["Earth",6371,5973.6],["Moon",1737,73.5],["Mars",3390,641.85]]

print(tabulate(table))

print(tabulate(table, headers=["Planet","R (km)", "mass (x 10^29 kg)"]))


print(tabulate([["Name","Age"],["Alice",24],["Bob",19]],headers="firstrow"))

print(tabulate({"Name": ["Alice", "Bob"],"Age": [24, 19]}, headers="keys"))

print(tabulate([["F",24],["M",19]], showindex="always"))

table = [["spam",42],["eggs",451],["bacon",0]]
headers = ["item", "qty"]
print(tabulate(table, headers, tablefmt="plain"))

print(tabulate(table, headers, tablefmt="simple"))
print(tabulate(table, headers, tablefmt="grid"))
print(tabulate(table, headers, tablefmt="github"))
print(tabulate(table, headers, tablefmt="fancy_grid"))
print(tabulate(table, headers, tablefmt="pretty"))
print(tabulate(table, headers, tablefmt="orgtbl"))
print(tabulate(table, headers, tablefmt="html"))
print(tabulate([[1.2345],[123.45],[12.345],[12345],[1234.5]], numalign="right"))

import csv 
from StringIO import StringIO
table = list(csv.reader(StringIO("spam, 42\neggs, 451\n")))
print(tabulate(table))


