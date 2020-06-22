import paramiko
ssh_client = paramiko.SSHClient()
ssh_client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
ssh_client.connect(hostname="10.71.1.3", username="tsl", password="tsl")

stdin, stdout, stderr = ssh_client.exec_command("ls")
#print(stdout.readlines())

scp = ssh_client.open_sftp()
scp.get("1.txt", "1.txt")
scp.close()
