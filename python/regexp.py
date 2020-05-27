import re

print(re.__doc__)
print(re.__all__)

test_string = "123 is ### at 234 dkkdkj repeated 122 times"
print(re.sub("\s+repeated\s+\d+\s+times", "bb", test_string))
print(re.sub("\D+", "aaa", test_string))
