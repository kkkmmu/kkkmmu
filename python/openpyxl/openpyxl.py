from openpyxl import Workbook
wb = Workbook()
ws1 = wb.create_sheet()
ws2 = wb.create_sheet()
ws1.title = "New Title 1"
ws2.title = "New Title 2"
ws1['A4'] = 111
ws2['A4'] = 222
wb.save('first_py_excel.xlsx')
