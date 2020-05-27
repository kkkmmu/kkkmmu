cint_reset();

int unit = 0;
int i;
bcm_vrf_t       vrf = 0;
bcm_if_t egr_obj;  
bcm_l3_egress_t l3_egress;
bcm_l3_route_t route;
bcm_mpls_exp_map_t exp_map;

bcm_ip_t ip1 = 0x02020202; /* 2.2.2.2 */
bcm_ip_t ip2 = 0x04040404; /* 4.4.4.4 */
bcm_mac_t mac_1 = {0x00, 0x00, 0x00, 0x00, 0x00, 0x01};
bcm_mac_t mac_2 = {0x00, 0x00, 0x00, 0x00, 0x00, 0x02};
bcm_mac_t local_mac = {0x00, 0x00, 0x00, 0x00, 0x11, 0x11};    
bcm_vlan_t      vid_1 = 21;
bcm_vlan_t      vid_2 = 22; 

bcm_port_t      port_1 = 1;  /* egress port */ 
bcm_port_t      port_2 = 3;  /* ingress port */ 

uint32 label_1 = 0x12345;
uint32 label_2 = 0x6789a;
uint8 ttl_1 = 16;
uint8 ttl_2 = 1;

bcm_gport_t     gport_1, gport_2;    

bcm_port_gport_get(unit, port_1, &gport_1);
bcm_port_gport_get(unit, port_2, &gport_2);

/* DSCP map to internal priority and cng setup */
int ing_map, egr_map;
bcm_qos_map_t map;
bcm_qos_map_t_init(&map);
print bcm_qos_map_create(unit, BCM_QOS_MAP_INGRESS | BCM_QOS_MAP_L3, &ing_map);
for (i=0;i<8;i++) {
	map.int_pri=i;
	map.dscp=i;
	print bcm_qos_map_add(unit, BCM_QOS_MAP_INGRESS | BCM_QOS_MAP_L3, &map, ing_map);
}
print bcm_qos_port_map_set(unit, gport_1, ing_map, -1);

/* Internal priority and cng to DSCP, */ 
for(i=0;i<8;i++) {
	print bcm_port_dscp_unmap_set(unit, gport_1, i, 0, i);
}
print bcm_port_control_set(unit, gport_1, bcmPortControlEgressModifyDscp, 1);

/*
#===============================================================================
# L3 MPLS Decapsulation ----> Route ----> L3 MPLS Encapsulation
#
# Ingres port 3 - Ingress Packet, MPLS
# 00 00 00 00 11 11 00 00 00 00 00 02 81 00 00 16 
# 88 47 12 34 58 10 67 89 A5 01
# 45 1C 00 2E 00 00 00 00 3F 72 6F 37 04 04 04 04 02 02 02 02 
# 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 
# 00 00 00 00 00 00
#
# Egress port 1 - Captured Packet, Ethernet
# 00 00 00 00 00 02 00 00 00 00 11 11 81 00 00 16
# 88 47 67 89 AA 01 12 34 55 10 
# 45 1C 00 2E 00 00 00 00 0F 72 9F 37 04 04 04 04 02 02 02 02 
# 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 
# 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 
# E2 D4 1F A2
#===============================================================================
 */

/* VLAN setup */
print bcm_vlan_create(unit, vid_1);
print bcm_vlan_gport_add(unit, vid_1, gport_1, 0);
print bcm_vlan_gport_add(unit, vid_1, gport_2, 0);
print bcm_vlan_create(unit, vid_2);
print bcm_vlan_gport_add(unit, vid_2, gport_2, 0);
print bcm_vlan_gport_add(unit, vid_2, gport_1, 0);

/* Enable egress mode */
bcm_switch_control_set(unit, bcmSwitchL3EgressMode, 1);    

/* Create L3 interface for CUSTOMER port */
bcm_l3_intf_t l3_intf_1;
bcm_l3_intf_t_init(&l3_intf_1);
sal_memcpy(l3_intf_1.l3a_mac_addr, local_mac, 6);
l3_intf_1.l3a_vid = vid_1;            
l3_intf_1.l3a_vrf = vrf;
l3_intf_1.l3a_flags = BCM_L3_ADD_TO_ARL;
print bcm_l3_intf_create(unit, &l3_intf_1);

/* Create L3 interface for MPLS tunnel */
bcm_l3_intf_t l3_intf_2;
bcm_l3_intf_t_init(&l3_intf_2);
sal_memcpy(l3_intf_2.l3a_mac_addr, local_mac, 6);
l3_intf_2.l3a_vid = vid_2;
l3_intf_2.l3a_vrf = vrf;
l3_intf_2.l3a_flags = BCM_L3_ADD_TO_ARL;
print bcm_l3_intf_create(unit, &l3_intf_2);  

/* EXP remapping for tunnel label. Internal priority => EXP, 0=>7, 1=>6, ... */
int init_exp_map_id_tnl;
print bcm_mpls_exp_map_create(unit, BCM_MPLS_EXP_MAP_EGRESS, &init_exp_map_id_tnl);
for (i = 0; i < 8; i++) {
	bcm_mpls_exp_map_t_init(&exp_map);
	exp_map.exp = i;
	exp_map.priority = 8 - i - 1;
	print bcm_mpls_exp_map_set(unit, init_exp_map_id_tnl, &exp_map);
}

/* Set MPLS tunnel initiator */    
bcm_mpls_egress_label_t tun_label;
tun_label.flags = BCM_MPLS_EGRESS_LABEL_TTL_SET | BCM_MPLS_EGRESS_LABEL_EXP_REMARK;
tun_label.label = label_2;
tun_label.qos_map_id = init_exp_map_id_tnl;
tun_label.ttl = ttl_2;
tun_label.pkt_pri = 0;
tun_label.pkt_cfi = 0;
print bcm_mpls_tunnel_initiator_set(unit, l3_intf_2.l3a_intf_id, 1, &tun_label);

/* EXP remapping for vc label. Internal priority => EXP, 0=>0, 1=>1, ... */
int init_exp_map_id_vc;
print bcm_mpls_exp_map_create(unit, BCM_MPLS_EXP_MAP_EGRESS, &init_exp_map_id_vc);
for (i = 0; i < 8; i++) {
	bcm_mpls_exp_map_t_init(&exp_map);
	exp_map.exp = i;
	exp_map.priority = i;
	print bcm_mpls_exp_map_set(unit, init_exp_map_id_vc, &exp_map);
}

/* Create L3 egress object for MPLS tunnel */
bcm_l3_egress_t_init(&l3_egress);
l3_egress.intf = l3_intf_2.l3a_intf_id;
sal_memcpy(l3_egress.mac_addr, mac_2, 6);  
l3_egress.vlan = vid_2;
l3_egress.port = gport_1;
l3_egress.mpls_label = label_1;
l3_egress.mpls_qos_map_id = init_exp_map_id_vc;
l3_egress.mpls_flags |= BCM_MPLS_EGRESS_LABEL_TTL_SET | BCM_MPLS_EGRESS_LABEL_EXP_REMARK;
l3_egress.mpls_ttl = ttl_1;

print bcm_l3_egress_create(unit, BCM_L3_ROUTE_LABEL, &l3_egress, &egr_obj);    

/* Add route to SERVICE PROVIDER */
bcm_l3_route_t_init(&route);
route.l3a_subnet = ip1;
route.l3a_ip_mask = 0xFFFFFFFF;
route.l3a_intf = egr_obj;
route.l3a_vrf = vrf;

print bcm_l3_route_add(unit, &route);    


/* Create L3 VPN */
bcm_mpls_vpn_config_t vpn_info;
bcm_mpls_vpn_config_t_init(&vpn_info);  
vpn_info.flags = BCM_MPLS_VPN_L3;
vpn_info.lookup_id = vrf;
print bcm_mpls_vpn_id_create(unit, &vpn_info);

/* Create L3 egress object for Customer port */
bcm_l3_egress_t_init(&l3_egress);
l3_egress.intf = l3_intf_1.l3a_intf_id;
sal_memcpy(l3_egress.mac_addr, mac_1, 6);    
l3_egress.vlan = vid_1;
l3_egress.port = gport_2;
print bcm_l3_egress_create(unit, 0, &l3_egress, &egr_obj);

/*
/* Add L3 route */
bcm_l3_route_t_init(&route);
route.l3a_subnet = ip2;
route.l3a_ip_mask = 0xFFFFFFFF;
route.l3a_intf = egr_obj;
route.l3a_vrf = vrf;
print bcm_l3_route_add(unit, &route);      
*/

/* Install L2 tunnel MAC */
print bcm_l2_tunnel_add(unit, local_mac, vid_2);

int tnl_exp_map_id;
print bcm_mpls_exp_map_create(unit, BCM_MPLS_EXP_MAP_INGRESS, &tnl_exp_map_id);
for (i = 0; i < 8; i++) {
	bcm_mpls_exp_map_t_init(&exp_map);
	exp_map.exp = i;
	exp_map.priority = i;
	exp_map.dscp = i;
	print bcm_mpls_exp_map_set(unit, tnl_exp_map_id, &exp_map);
}

/* Add ILM for incoming MPLS tunnel label */
bcm_mpls_tunnel_switch_t info;
print bcm_mpls_tunnel_switch_t_init(&info);
info.flags = BCM_MPLS_SWITCH_INT_PRI_MAP | BCM_MPLS_SWITCH_INNER_EXP;
info.label = label_1;
info.port = BCM_GPORT_INVALID;
info.action = BCM_MPLS_SWITCH_ACTION_POP;
info.exp_map = tnl_exp_map_id;
info.vpn = vpn_info.vpn;
print bcm_mpls_tunnel_switch_add(unit, &info);

int vc_exp_map_id;
print bcm_mpls_exp_map_create(unit, BCM_MPLS_EXP_MAP_INGRESS, &vc_exp_map_id);
for (i = 0; i < 8; i++) {
	bcm_mpls_exp_map_t_init(&exp_map);
	exp_map.exp = i;
	exp_map.priority = i;
	exp_map.dscp = 8-i-1;
	print bcm_mpls_exp_map_set(unit, vc_exp_map_id, &exp_map);
}

/* Add ILM for incoming MPLS VC label */
print bcm_mpls_tunnel_switch_t_init(&info);
/* 
 * With BCM_MPLS_SWITCH_OUTER_EXP/BCM_MPLS_SWITCH_OUTER_TTL flag, 
 * DSCP from ING_MPLS_EXP_MAPPING. 
 * Otherwise, from EGR_DSCP_TABLE.
 * Without any mapping, the internal DSCP will be kept.  
 */
info.flags = BCM_MPLS_SWITCH_INT_PRI_MAP | BCM_MPLS_SWITCH_INNER_EXP; 
info.label = label_2;
info.port = BCM_GPORT_INVALID;
info.action = BCM_MPLS_SWITCH_ACTION_POP;
info.vpn = vpn_info.vpn;
info.exp_map = vc_exp_map_id;
info.ingress_if = l3_intf_1.l3a_intf_id; /* L3_IIF */
print bcm_mpls_tunnel_switch_add(unit, &info);
