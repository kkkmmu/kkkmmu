# $language = "python"
# $interface = "1.0"


def main():
    crt.Screen.IgnoreEscape = True
    crt.Screen.Synchronous = True
    crt.Screen.IgnoreCase = True
    crt.Screen.Clear()
    crt.Screen.Send('ls -al \n')
    crt.Screen.WaitForStrings(['#', '$'], 5)
    crt.Dialog.MessageBox("{} {} {} {}".format(
        crt.Screen.CurrentRow, crt.Screen.CurrentColumn, crt.Screen.Rows, crt.Screen.Columns))
    crt.Dialog.MessageBox(crt.Screen.Get(1, 1, 10, 56))
    crt.Screen.Send("ps aux\r\n")
main()
