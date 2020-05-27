import textfsm
# Run text through the FSM.
# The argument 'template' is a file handle and 'raw_text_data' is a string.

with open("cisco.txt", "r") as fd:
    with open("cisco.template", "r") as td:
        re_table = textfsm.TextFSM(td)
        data = re_table.ParseText(fd.read())
        print(', '.join(re_table.header))
        for row in data:
            print(', '.join(row))

# Display result as CSV
# First the column headers
# Each row of the table.
