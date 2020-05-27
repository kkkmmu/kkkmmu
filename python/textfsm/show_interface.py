import textfsm

with open("show_interface.template") as td:
    with open("show_interface.txt") as sd:
        re_table = textfsm.TextFSM(td)
        data = re_table.ParseText(sd.read())
        print(', '.join(re_table.header))
        for row in data:
            print(', '.join(row))
