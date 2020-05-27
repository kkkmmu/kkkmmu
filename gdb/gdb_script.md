set $og = (struct lib_globals *)0x1019b6dc
set $vr = (struct ipi_vr *)*$og->vr_vec->index
set $om = (struct ospf_master *)$vr->proto
set $ospf_list = $om->ospf
set $ospf_top = (struct ospf *)$ospf_list->head->data
set $nbr_table = $ospf_top->nbr_table
 
printf "ospf router_id %08x\n", $ospf_top->router_id.s_addr
define walk_nbr
  set $rn = (struct ls_node *)$arg1
  set $left_$arg0 = $rn->link[0]
  set $right_$arg0 = $rn->link[1]
  set $nbr = (struct ospf_neighbor *)$rn->vinfo[0]
  if $nbr
        printf "\nnbr %p state %d ostate %d change %6d router_id %08x dr %08x bdr %08x dd.recv.mtu %d dd.recv.seqnum %d\n", \
        $nbr, $nbr->state, $nbr->ostate, $nbr->state_change, $nbr->ident.router_id.s_addr, $nbr->ident.d_router.s_addr, $nbr->ident.bd_router.s_addr, $nbr->dd.recv.mtu, $nbr->dd.recv.seqnum
        set $oi = $nbr->oi
        printf "  oi %s if_mtu %d full_nbr_cnt %d hello_in_out(%d %d) lsupd_in_out(%d %d) lsack_in_out(%d %d) discarded %d change %d\n", \
        $oi->u.ifp->name, $oi->u.ifp->mtu, $oi->full_nbr_count, $oi->hello_in, $oi->hello_out, $oi->ls_upd_in, $oi->ls_upd_out, $oi->ls_ack_in, $oi->ls_ack_out, $oi->discarded, $oi->state_change
#    pt *$oi
#    p *$oi
  end
  if $left_$arg0
        walk_nbr left_$arg0 $left_$arg0
  end
  if $right_$arg0
        walk_nbr right_$arg0 $right_$arg0
  end
end
walk_nbr top $nbr_table->top
  
============= 출력 결과 =====
(gdb)  source ospf_nbr.gs
ospf router_id 0aba2b23
Redefine command "walk_nbr"? (y or n) [answered Y; input not from terminal]
nnbr 0x102369e4 state 8 ostate 8 change 135232 router_id 0a482b0c dr 0abb9252 bdr 00000000 dd.recv.mtu 0 dd.recv.seqnum 0
  oi br20 if_mtu 2000 full_nbr_cnt 1 hello_in_out(683017 135228) lsupd_in_out(154381 7190) lsack_in_out(1 137894) discarded 1 change 683018
nnbr 0x102374ac state 4 ostate 4 change      3 router_id 0aba2b1a dr 0abb9252 bdr 00000000 dd.recv.mtu 0 dd.recv.seqnum 0
  oi br20 if_mtu 2000 full_nbr_cnt 1 hello_in_out(683017 135228) lsupd_in_out(154381 7190) lsack_in_out(1 137894) discarded 1 change 683018
nnbr 0x10238534 state 4 ostate 4 change      3 router_id 0aba2b1e dr 0abb9252 bdr 00000000 dd.recv.mtu 0 dd.recv.seqnum 0
  oi br20 if_mtu 2000 full_nbr_cnt 1 hello_in_out(683017 135228) lsupd_in_out(154381 7190) lsack_in_out(1 137894) discarded 1 change 683018
nnbr 0x10236134 state 4 ostate 4 change      3 router_id 0aba2b4d dr 0abb9252 bdr 00000000 dd.recv.mtu 0 dd.recv.seqnum 0
  oi br20 if_mtu 2000 full_nbr_cnt 1 hello_in_out(683017 135228) lsupd_in_out(154381 7190) lsack_in_out(1 137894) discarded 1 change 683018
nnbr 0x10237cbc state 4 ostate 4 change      3 router_id 0aba2b4e dr 0abb9252 bdr 00000000 dd.recv.mtu 0 dd.recv.seqnum 0
  oi br20 if_mtu 2000 full_nbr_cnt 1 hello_in_out(683017 135228) lsupd_in_out(154381 7190) lsack_in_out(1 137894) discarded 1 change 683018
