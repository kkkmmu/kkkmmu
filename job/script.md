How to push up to 3 labels for MPLS

Up to 3 labels (1 VC label and 2 tunnel labels) can be pushed on.

/*
# 1. input packet on GE1: L3 packet
#       00 00 00 00 04 44 00 00 00 00 00 01 81 00 00 0B 08 00 45 00 00 32 00 00 00 00 40 72 13 44 C6 13 01 01 0A 01 96 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
#    request: packet egress out of GE2 with 1 VC label (0x666) + 2 MPLS label (0x444 and 0x555)
#       00 00 00 00 02 22 00 00 00 00 04 44 81 00 00 0A 88 47 00 44 40 10 00 55 50 10 00 66 61 3F 45 00 00 32 00 00 00 00 3F 72 14 44 C6 13 01 01 0A 01 96 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
 */
cint_reset();
void l3_mpls_add_3_labels()
{
	int     unit = 0;
	int     vrf = 0;
	int     i;
	bcm_mpls_vpn_config_t vpn_info;

	bcm_port_t port_1 = 2;  /* GE1: customer port    */
	bcm_port_t port_2 = 3;  /* GE2: provider port */
	bcm_gport_t gport_1, gport_2;
	bcm_vlan_t vid_1 = 11; /* customer VLAN */
	bcm_vlan_t vid_2 = 10; /* service provider VLAN */

	int  intf_id = 2;
	int  ttl = 16;
	int  num_labels = 2;

	uint32 tunnel_label_init_1 = 0x555;  /* 1365 */
	uint32 tunnel_label_init_2 = 0x444;  /* 1092 */
	uint32 vc_label_term     = 0x777;  /* 1911 */
	uint32 vc_label_init     = 0x666;  /* 1638 */

	bcm_l3_egress_t          l3_egress;
	bcm_if_t                 l3_egr_obj;
	bcm_l3_intf_t            l3_intf;
	bcm_mpls_egress_label_t  mpls_egress_label[2];
	bcm_mac_t remote_mac = {00, 00, 00, 00, 0x2, 0x22};
	bcm_mac_t local_mac  = {00, 00, 00, 00, 0x4, 0x44};

	bcm_mpls_port_t          cust_mpls_port;       /* customer port */
	bcm_mpls_port_t          service_mpls_port;  /* provider port */

	bcm_mpls_tunnel_switch_t info;


	/*
	 * Initialize gport values
	 */
	bcm_port_gport_get(unit, port_1, &gport_1);
	printf("GE1: gport_1=0x%x\n", gport_1);
	bcm_port_gport_get(unit, port_2, &gport_2);
	printf("GE2: gport_1=0x%x\n", gport_2);

	/* Create vlan 1 */
	bcm_vlan_create(unit, vid_1);
	print bcm_vlan_gport_add(unit, vid_1, gport_1, 0);

	/* Create vlan 2 */
	bcm_vlan_create(unit, vid_2);
	print bcm_vlan_gport_add(unit, vid_2, gport_2, 0);

	/* VRF for customer VLAN */
	bcm_vlan_control_vlan_t vlan_control;
	print bcm_vlan_control_vlan_get(unit, vid_1, &vlan_control);
	vlan_control.vrf = vrf;
	print bcm_vlan_control_vlan_set(unit, vid_1, vlan_control);

	/* Create L3 interface for L3 route */
	bcm_l3_intf_t_init(&l3_intf);
	l3_intf.l3a_flags = BCM_L3_WITH_ID | BCM_L3_ADD_TO_ARL;
	l3_intf.l3a_intf_id = vid_1;
	sal_memcpy(l3_intf.l3a_mac_addr, local_mac, 6);
	l3_intf.l3a_vid = vid_1;
	l3_intf.l3a_vrf = vrf;
	print bcm_l3_intf_create(unit, &l3_intf);

	/*
	 * Enable L3 egress mode & VLAN translation
	 */
	bcm_switch_control_set(unit, bcmSwitchL3EgressMode, 1);
	bcm_vlan_control_set(unit, bcmVlanTranslate, 1);
	bcm_switch_control_set(unit, bcmSwitchL2StaticMoveToCpu, 1);

	/*
	 * Create tunnel
	 */

	/*
	 * Create L3 MPLS VPN
	 */
	bcm_mpls_vpn_config_t_init(&vpn_info);
	vpn_info.flags = BCM_MPLS_VPN_L3;
	print bcm_mpls_vpn_id_create(unit, &vpn_info);
	printf("vpn_id=%d\n", vpn_info.vpn);

	/* Create L3 interface for MPLS tunnel */
	bcm_l3_intf_t_init(&l3_intf);
	l3_intf.l3a_flags = BCM_L3_WITH_ID | BCM_L3_ADD_TO_ARL;
	l3_intf.l3a_intf_id = intf_id;
	sal_memcpy(l3_intf.l3a_mac_addr, local_mac, 6);
	l3_intf.l3a_vid = vid_2;
	l3_intf.l3a_vrf = vrf;
	print bcm_l3_intf_create(unit, &l3_intf);

	/* Set MPLS tunnel initiator */
	bcm_mpls_egress_label_t_init(&mpls_egress_label[0]);
	bcm_mpls_egress_label_t_init(&mpls_egress_label[1]);
	mpls_egress_label[0].flags = BCM_MPLS_EGRESS_LABEL_TTL_SET;
	mpls_egress_label[0].label = tunnel_label_init_1;
	mpls_egress_label[0].ttl     = ttl;
	mpls_egress_label[1].flags = BCM_MPLS_EGRESS_LABEL_TTL_SET;
	mpls_egress_label[1].label = tunnel_label_init_2;
	mpls_egress_label[1].ttl     = ttl;
	print bcm_mpls_tunnel_initiator_set(unit, intf_id, num_labels, mpls_egress_label);

	/* Create L3 egress object for MPLS tunnel */
	bcm_l3_egress_t_init(&l3_egress);
	l3_egress.intf = intf_id;
	l3_egress.mpls_label = vc_label_init;
	l3_egress.port = gport_2;
	l3_egress.vlan = vid_2;    /* VLAN field not used, but API requires it to be a valid VLAN */
	sal_memcpy(l3_egress.mac_addr, remote_mac, 6);
	print bcm_l3_egress_create(unit, 0, &l3_egress, &l3_egr_obj);

	/* Add IP route into DEFIP table */
	bcm_l3_route_t l3_route_info;
	bcm_ip_t        ip_addr = 0x0A019601;
	bcm_ip_t        ip_mask = 0xffffff00;
	bcm_l3_route_t_init(&l3_route_info);
	l3_route_info.l3a_vrf = vrf;
	l3_route_info.l3a_subnet  = ip_addr;
	l3_route_info.l3a_ip_mask = ip_mask;
	l3_route_info.l3a_intf = l3_egr_obj;
	print bcm_l3_route_add(unit, l3_route_info);

	/* My station MAC */
	/*print bcm_l2_tunnel_add(unit, local_mac, vid_2);*/

}

l3_mpls_add_3_labels();
