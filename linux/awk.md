1. Print Matched pattern. 

awk '{
	for(i=1; i<=NF; i++) {
		tmp=match($i, /[0-9]..?.?[^A-Za-z0-9]/)
			if(tmp) {
				print $i
			}
	}
}' $1

2. Re-organize ss output 
ss -p | awk -F ' ' '{ 
	if ($1 == "u_str") printf("%40s:%-16s %s %17s:%-10s %-10s\n", $5, $6, $2, $7, $8, $9); 
	if ($1 == "tcp" || $1 == "udp") 	printf("%44s %18s %22s\n", $5,  $2, $6); 
}' 

3. Print regexp matched

awk -F '=' '{
	for(i=1; i<=NF; i++) {
		tmp=match($i, /{[0-9a-z]+\,/)
			if(tmp) {
				print substr($i, RSTART+1, RLENGTH-2)
			}
		}
	}' $1
