package main

import (
	"command"
	"context"
	"flag"
	"fmt"
	"rut"
)

var (
	IP       = flag.String("ip", "10.71.20.182", "IP address of the remote device")
	Server   = flag.String("server", "10.71.1.3", "IP address of file server")
	Port     = flag.String("port", "23", "Port to connect")
	Host     = flag.String("hostname", "SWITCH", "Host name of the remote device")
	Protocol = flag.String("prot", "telnet", "Login protocol(ssh/telnet)")
	User     = flag.String("u", "admin", "Username of the remote device")
	Password = flag.String("p", "Dasan123456", "Passwrod of the remote device")

	Process = flag.String("process", "hsl", "process name to debug")
	Bin     = flag.String("bin", "hsl", "binary file name path")

	CTX = context.Background()
)

var Script = `
cint
cint_reset();
bcm_error_t
setup_IPv4_route_entries(int unit, int routeIPv4TableCount, bcm_if_t if_id,
bcm_ip_t base_local_ipv4, int keep_trying)
{
	bcm_error_t rv;
	bcm_l3_route_t info;
	int i;
	int missed = 0;
	int first_miss = -1;
	int added;
	bcm_ip_t ipv4mask = ~0xF; /* Prefix length = 28 */

	base_local_ipv4 &= ipv4mask;
	printf("Set up %d IPv4 route entries at %d.%d.%d.%d\n", routeIPv4TableCount,
	(base_local_ipv4 >> 24) & 0xFF,
	(base_local_ipv4 >> 16) & 0xFF,
	(base_local_ipv4 >> 8) & 0xFF, (base_local_ipv4 >> 0) & 0xFF);
	for (i = 0; i < routeIPv4TableCount; i++) {
		bcm_ip_t local_ipv4 = base_local_ipv4 + (i << 4);

		/* Add L3 ipV4 route entry. */
		bcm_l3_route_t_init(&info);
		info.l3a_subnet = local_ipv4;
		info.l3a_ip_mask = ipv4mask;
		info.l3a_intf = if_id;
		info.l3a_vrf = -2;
		if (BCM_FAILURE(rv = bcm_l3_route_add(unit, &info))) {
			if (keep_trying && (rv == BCM_E_FULL)) {
				/* Keep trying to add new entries even after first table full */
				if (first_miss < 0) {
					first_miss = i + 1;
				}
				missed++;
			} else {
				printf(" bcm_l3_route_add() failed at iteration %d: %s\n", i + 1,
				bcm_errmsg(rv));
				missed = routeIPv4TableCount - i;
				break;
			}
		}
	}
	added = routeIPv4TableCount - missed;
	printf(" Added %d IPv4 route entries;", added);
	if (keep_trying) {
		printf(" %d missed;", missed);
		if (missed) {
			printf(" first miss: %d;\n", first_miss);
		}
	}
	printf("\n");
	return ((rv != BCM_E_NONE) && (rv != BCM_E_FULL)) ? rv :
	(keep_trying && (added > 0)) ||
	(!keep_trying && (added == routeIPv4TableCount)) ? BCM_E_NONE : BCM_E_FULL;
}

bcm_error_t
create_vlan(int unit, bcm_vlan_t vlan, bcm_port_t port)
{
	bcm_error_t rv;
	bcm_pbmp_t port_list;
	bcm_pbmp_t untagged;

	BCM_PBMP_PORT_SET(port_list, port);
	BCM_PBMP_CLEAR(untagged); /* Never untagged */

	rv = bcm_vlan_create(unit, vlan);
	if (BCM_FAILURE(rv) && (rv != BCM_E_EXISTS)) {
		return rv;
	}

	BCM_IF_ERROR_RETURN(bcm_vlan_port_add(unit, vlan, port_list, untagged));

	return BCM_E_NONE;
}

bcm_error_t
l3_info(int unit)
{
	bcm_l3_info_t l3info;

	BCM_IF_ERROR_RETURN(bcm_l3_info(unit, &l3info));

	printf("L3 INFO:\n");
	printf(" L3 host table size (unit is IPv4 unicast): %d\n", l3info.l3info_max_host);
	printf(" L3 host entries used: %d (%d remaining, %d%% utilization)\n",
	l3info.l3info_used_host, l3info.l3info_max_host - l3info.l3info_used_host,
	((l3info.l3info_used_host * 100) +
	(l3info.l3info_used_host / 2)) / l3info.l3info_max_host);
	printf(" L3 route table size (unit is IPv4 route): %d\n", l3info.l3info_max_route);
	printf(" L3 route entries used: %d\n", l3info.l3info_used_route);
	printf(" NextHops used: %d\n", l3info.l3info_used_nexthop);
	printf(" L3 interfaces used: %d\n", l3info.l3info_used_intf);
	printf(" LPM blocks used: %d\n", l3info.l3info_used_lpm_block);
	printf(" Maximum ECMP paths allowed: %d\n", l3info.l3info_max_ecmp);
	printf(" Maximum IPV4 tunnels that can be initiated: %d\n",
	l3info.l3info_max_tunnel_init);
	printf(" Maximum IPV4 tunnels that can be terminated: %d\n",
	l3info.l3info_max_tunnel_term);
	printf(" Maximum L3 interface groups the chip supports: %d\n",
	l3info.l3info_max_intf_group);
	printf(" Maximum L3 interfaces the chip supports: %d\n", l3info.l3info_max_intf);
	printf(" Maximum LPM blocks: %d\n", l3info.l3info_max_lpm_block);
	printf(" Maximum NextHops: %d\n", l3info.l3info_max_nexthop);
	printf(" Maximum number of virtual routers allowed: %d\n", l3info.l3info_max_vrf);
	printf(" Number of active IPV4 tunnels initiated: %d\n",
	l3info.l3info_used_tunnel_init);
	printf(" Number of active IPV4 tunnels terminated: %d\n",
	l3info.l3info_used_tunnel_term);
	printf(" Number of virtual routers created so far: %d\n", l3info.l3info_used_vrf);
	return BCM_E_NONE;
}



bcm_error_t
create_route_entries(int unit, int routeIPv4TableCount, bcm_vlan_t local_vid, bcm_mac_t local_mac, 
bcm_vlan_t remote_vid, bcm_port_t remote_port, 
bcm_mac_t remote_mac, bcm_ip_t base_route_ipv4, int keep_trying)
{
	const int leastFull = TRUE;

	/* Program variables */
	bcm_error_t rv;
	bcm_gport_t remote_gport;
	bcm_if_t if_id;
	bcm_l3_egress_t l3_egress;
	bcm_l3_intf_t intf;
	int i;

	BCM_IF_ERROR_RETURN(bcm_port_gport_get(unit, remote_port, &remote_gport));

	BCM_IF_ERROR_RETURN(create_vlan(unit, remote_vid, remote_port));

	/* Create L3 Interface */
	bcm_l3_intf_t_init(&intf);
	intf.l3a_flags = BCM_L3_ADD_TO_ARL;
	intf.l3a_mac_addr = local_mac;
	intf.l3a_vid = local_vid;
	BCM_IF_ERROR_RETURN(bcm_l3_intf_create(unit, &intf));

	/* Create L3 Egress Object */
	bcm_l3_egress_t_init(&l3_egress);
	l3_egress.mac_addr = remote_mac;
	l3_egress.intf = intf.l3a_intf_id;
	l3_egress.vlan = remote_vid;
	l3_egress.port = remote_gport;

	BCM_IF_ERROR_RETURN(bcm_l3_egress_create(unit, 0, &l3_egress, &if_id));

	/* STEP 1: Create IPv4 route entries */
	if (routeIPv4TableCount) {
		if (BCM_FAILURE(rv = setup_IPv4_route_entries(unit, routeIPv4TableCount, if_id,
		base_route_ipv4, keep_trying))) {
			BCM_IF_ERROR_RETURN(l3_info(unit));
			return rv;
		}
	}
	return 0;
}


bcm_error_t
ipv4_lpm_test(int unit)
{
	/* Program constants */
	const bcm_mac_t local_mac = { 0x00, 0x00, 0x00, 0x00, 0x21, 0x21 };
	const bcm_vlan_t local_vid1 = 21;
	const bcm_vlan_t local_vid2 = 22;

	const bcm_port_t remote_port1 = 5;
	const bcm_port_t remote_port2 = 6;
	const bcm_mac_t remote_mac1 = { 0x00, 0x00, 0x00, 0x00, 0x31, 0x31 };
	const bcm_mac_t remote_mac2 = { 0x00, 0x00, 0x00, 0x00, 0x32, 0x32 };
	const bcm_vlan_t remote_vid1 = 31;
	const bcm_vlan_t remote_vid2 = 32;

	const bcm_ip_t base_route_ipv41 = 191 << 24 | 00 << 16 | 0 << 8 | 0 << 0;
	const bcm_ip_t base_route_ipv42 = 192 << 24 | 00 << 16 | 0 << 8 | 0 << 0;

	BCM_IF_ERROR_RETURN(bcm_switch_control_set(unit, bcmSwitchL3EgressMode, TRUE));

	create_route_entries(0, 20000, local_vid1, local_mac, remote_vid1, remote_port1, remote_mac1, base_route_ipv41, 1);
	create_route_entries(0, 20000, local_vid2, local_mac, remote_vid2, remote_port2, remote_mac2, base_route_ipv42, 1);

	printf("TEST COMPLETE\n");
	BCM_IF_ERROR_RETURN(l3_info(unit));
	return BCM_E_NONE;
}
ipv4_lpm_test(0);
exit;
d nocache chg L3_DEFIP_ALPM_IPV4
`

func main() {
	flag.Parse()
	dev, err := rut.New(&rut.RUT{
		Name:     "SWITCH",
		Device:   "V5",
		Protocol: *Protocol,
		IP:       *IP,
		Port:     *Port,
		Username: *User,
		Hostname: *Host,
		Password: *Password,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Start connect to %s!\n", *IP)
	err = dev.Init()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Connect to %s sucessfull!\n", *IP)

	ctx := context.Background()

	fmt.Println(dev.CurrentMode())
	if dev.CurrentMode() == "normal" {
		_, err = dev.RunCommand(ctx, &command.Command{
			Mode: "normal",
			CMD:  "q sh -l",
		})

		if err != nil {
			panic(err)
		}
	}

	data, err := dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  "cd /etc/sdk",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  "bcm.user",
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(data)

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "bcmshell",
		CMD:  "init all",
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(data)

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "bcmshell",
		CMD:  Script,
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}
