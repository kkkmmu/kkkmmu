#!/usr/bin/python
import re
import sys
import gdb
import traceback
import six
from collections import defaultdict
from curses.ascii import isgraph
from inspect import getframeinfo, stack
import socket


class HelloWorld(gdb.Command):
    def __init__(self):
        super(HelloWorld, self).__init__("hello-world", gdb.COMMAND_OBSCURE)

    def invoke(self, arguments, from_tty):
        print("Hello world")


HelloWorld()


class DeferencePointer(gdb.Command):
    def __init__(self):
        super(DeferencePointer, self).__init__("dp", gdb.COMMAND_OBSCURE)

    def invoke(self, arguments, from_tty):
        args = re.split('/', arguments)
        print(self.deref_pinter_by_type_name(
            args[0], args[1]))

    def deref_pinter_by_type_name(self, addr, type_name):
        type_pointer = gdb.lookup_type(type_name).pointer()
        addr = gdb.parse_and_eval(addr)
        val = gdb.Value(addr).cast(type_pointer).dereference()
        for field in val.type.keys():
            print(field, val[field])
        return gdb.Value(addr).cast(type_pointer).dereference()


DeferencePointer()


class DumpDatabase(gdb.Command):
    def __init__(self):
        super(DumpDatabase, self).__init__("xd", gdb.COMMAND_OBSCURE)

    def invoke(self, arguments, from_tty):
        arguments = re.sub(' +', ' ', arguments)
        args = arguments.split('/', -1)
        if len(args) < 2:
            print("Usuage xd POINTER/TYPE[/FIELDS]")
            return
        elif len(args) == 2:
            print(self.dump_database(args[0], args[1]))
        else:
            print(self.dump_database(args[0], args[1], args[2]))

    def dump_database(self, root, type_name, fields=None):
        # type_void_pointer = gdb.lookup_type("void").pointer()
        if fields != None:
            fields = re.split(',', fields)
        try:
            db = Listfy(root, type_name)
            for data in db:
                if data == 0x0:
                    continue
                out = ""
                if fields != None:
                    for field in fields:
                        if field.find('.'):
                            sub_fields = re.split('\.', field)
                            child = data
                            for sf in sub_fields:
                                if child.type.code == gdb.TYPE_CODE_PTR:
                                    child = child.dereference()[sf]
                                else:
                                    child = child[sf]
                            out += field + ":" + \
                                re.sub(",\s+'.*?' <repeats\s+\d+\s+times>",
                                       "", str(child))+", "
                        else:
                            out += field + ":" + \
                                re.sub(",\s+'.*?' <repeats\s+\d+\s+times>",
                                       "", str(data.dereference()[field]))+", "
                    print(out)
                else:
                    print(data.dereference())
        except:
            pass
            #    exc_info = sys.exc_info()
            #    traceback.print_exception(*exc_info)

    def branch(self):
        for branch in util.db_branch_names:
            for field in self.pNext.type.keys():
                if branch == field:
                    return branch

        return None

    def top(self):
        for rname in util.db_root_names:
            for field in self.pNext.type.keys():
                if field == rname:
                    self.pNext = self.pNext[field]
                return self.pNext
        return None


DumpDatabase()


class Listfy:
    def __init__(self, expr, data_type=None):
        if util.isstring(expr):
            if expr.isdigit():
                print("{} is not a valid variable name".format(root))
                return
            else:
                expr = gdb.parse_and_eval(expr)

            if expr.type.code == gdb.TYPE_CODE_PTR:
                self.pNext = expr
            else:
                self.pNext = expr.address
        else:
            if expr.type.code == gdb.TYPE_CODE_PTR:
                self.pNext = expr
            else:
                self.pNext = expr.address

        # if parent node of root is give, get the root node.
        self.top()
        self.first()

        br = self.branch()
        if br == None:
            print("Unknow database type")
            return

        self.count = 0
        if data_type != None:
            self.data_type = gdb.lookup_type(data_type).pointer()
        else:
            self.data_type = None

    def branch(self):
        for branch in util.db_branch_names:
            for field in self.pNext.dereference().type.keys():
                if branch == field:
                    return branch

        return None

    def top(self):
        for rname in util.db_root_names:
            for field in self.pNext.dereference().type.keys():
                if field == rname:
                    self.pNext = self.pNext.dereference()[field]
                    return self.pNext
        return None

    def first(self):
        if self.branch() in self.pNext.dereference().type.keys():
            # For binary tree we start from the smallest entry.
            if self.branch() == 'left':
                first = self.pNext.dereference()
                while(first['left']):
                    self.pNext = first['left']
                    first = self.pNext.dereference()
                return self.pNext

            # if self.branch() == 'link':
            #    first = self.pNext.dereference()
            #    while(first['link'][0]):
            #        self.pNext = first['link'][0]
            #        first = self.pNext.dereference()
            #    return self.pNext
        return None

    def __iter__(self):
        self.count = 0
        return self

    def __next__(self):
        if self.pNext == 0x0 or self.pNext == None:
            raise StopIteration

        try:
            curr_node = self.pNext.dereference()
            node = curr_node
            #self.pNext = 0
            if self.branch() == 'next':
                self.pNext = node['next']
            elif self.branch() == 'left':
                pnext = curr_node['right']
                if int(pnext) == 0:
                    pnext = curr_node['parent']
                    currt = curr_node
                    while int(pnext) != 0:
                        if int(pnext['left']) == int(currt.address):
                            break
                        else:
                            currt = pnext.dereference()
                            pnext = pnext['parent']
                    #curr_node = pnext.dereference()
                else:
                    ptmp = pnext['left']
                    if int(ptmp) != 0:
                        while int(pnext['left']) != 0:
                            pnext = pnext['left']

                self.pNext = pnext
            elif self.branch() == 'link':
                if node['link'][0]:
                    self.pNext = node['link'][0]
                elif node['link'][1]:
                    self.pNext = node['link'][1]
                else:
                    while int(node['parent']) != 0:
                        parent_node = node['parent'].dereference()
                        right = parent_node['link'][1]
                        if int(right) != 0 and right != node.address:
                            self.pNext = right
                            break
                        else:
                            node = parent_node

            if curr_node.address == self.pNext:
                self.pNext = None

            self.count += 1

            # If Data type is not give, we directly return the node.
            if self.data_type == None:
                # return curr_node.address
                return self.pNext

            # if Data type is given, return the data filed in this node.
            pdata = util.data_in_node(curr_node)
            if pdata != -1:
                pdata = pdata.cast(self.data_type)

            return pdata

        except:
            exc_info = sys.exc_info()
            traceback.print_exception(*exc_info)
            raise

    next = __next__


class util:
    db_branch_names = ['left', 'link', 'next']
    db_root_names = ['root', 'head', 'top']
    db_data_names = ['info', 'vinfo', 'data', 'entry']
    rib_pointer = gdb.lookup_type("struct rib").pointer()
    fifo_pointer = gdb.lookup_type("struct fifo").pointer()
    ip_vr_pointer = gdb.lookup_type("struct ipi_vr").pointer()
    ip_vrf_pointer = gdb.lookup_type("struct ipi_vrf").pointer()

    def isstring(expr):
        if (sys.version_info < (3, 0) and isinstance(expr, basestring)):
            return True

        if (sys.version_info >= (3, 0) and isinstance(expr, str)):
            return True
        return False

    def isdigit(expr):
        try:
            if util.ishex(expr):
                return True
            elif util.isstring(expr):
                return expr.isdigit()
            else:
                return false
        except:
            return False

    def ishex(expr):
        try:
            int(expr, 16)
            return True
        except:
            return False

    def tostring(a):
        if (sys.version_info < (3, 0) and isinstance(a, basestring)):
            return a.encode('ascii', 'ignore')

        if (sys.version_info >= (3, 0) and isinstance(a, str)):
            return a

    def data_in_node(node):
        fnames_in_node = node.type.keys()
        for dkey in util.db_data_names:
            if dkey in fnames_in_node:
                if dkey == 'vinfo':
                    pdata = node['vinfo'][0]
                else:
                    pdata = node[dkey]
                return pdata
        return -1

    def dump_raw(pointer, length):
        try:
            infer = gdb.selected_inferior()
            align = gdb.parameter('hxd-align')
            width = gdb.parameter('hxd-width')
            if width == 0:
                width = 16

            addr = int(gdb.parse_and_eval(pointer))
            mem = infer.read_memory(addr, int(length))
            offset = width

            if align:
                offset = width - (addr % width)
                addr -= addr % width

            for group in util.groups_of(mem, width, offset):
                out = ''
                group = group.tobytes()
                out += ('0x%x: ' % (addr) + '   '*(width - offset)) \
                    + ' '.join(['%02X' % g for g in group]) + \
                    ('   ' * (width - len(group) if offset == width else 0) + ' ') \
                    + (' '*(width - offset) + ''.join(
                        ['%c' % g if isgraph(g) or g == ' ' else '.' for g in group]))
                addr += width
                offset = width
                print(out)
        except:
            exc_info = sys.exc_info()
            traceback.print_exception(*exc_info)

    def groups_of(iterable, size, first=0):
        first = first if first != 0 else size
        chunk, iterable = iterable[:first], iterable[first:]
        while chunk:
            yield chunk
            chunk, iterable = iterable[:size], iterable[size:]

    def cast(pointer, type_name):
        type_pointer = gdb.lookup_type(type_name).pointer()
        pointer = gdb.parse_and_eval(pointer)
        return gdb.Value(pointer).cast(type_pointer)

    def dump_struct(pointer, type_name, fields=None):
        data = util.cast(pointer, type_name).dereference()
        if fields == None:
            print(data)
        else:
            fields = re.split(',', fields)
            out = ""
            for field in fields:
                if field.find('.'):
                    sub_fields = re.split('\.', field)
                    child = data
                    for sf in sub_fields:
                        if child.type.code == gdb.TYPE_CODE_PTR:
                            child = child.dereference()[sf]
                        else:
                            child = child[sf]
                            out += field + ":" + \
                                re.sub(",\s+'.*?' <repeats\s+\d+\s+times>",
                                       "", str(child))+", "
                else:
                    out += field + ":" + \
                        re.sub(",\s+'.*?' <repeats\s+\d+\s+times>",
                               "", str(data.dereference()[field]))+", "
            return out

    def get_hsl_db_and_type(fib, family, table):
        db = ""
        entry_type = ""
        if table == "route":
            if family == "i6":
                db = "p_hsl_fib_table->prefix6["+fib+"]"
            else:
                db = "p_hsl_fib_table->prefix["+fib+"]"
            entry_type = "struct hsl_prefix_entry"
        elif table == "nexthop":
            if family == "i6":
                db = "p_hsl_fib_table->nh6["+fib+"]"
            else:
                db = "p_hsl_fib_table->nh["+fib+"]"
            entry_type = "struct hsl_nh_entry"
        elif table == "interface":
            db = "p_hsl_if_db->if_tree"
            entry_type = "struct hsl_if"

        return (gdb.parse_and_eval(db), entry_type)

    def get_ribd_db_and_type(fib, family, table):
        db = ""
        entry_type = ""
        if table == "route" or table == "ptree":
            vr_addr = int(gdb.parse_and_eval(
                "rib_lib_globals->vr_vec.index[0]"))
            #vr_addr = gdb.Value(vr_addr).cast(util.ipi_vr_pointer)
            vrf_addr = int(gdb.parse_and_eval(
                "((struct ipi_vr*)"+str(vr_addr)+")->vrf_vec.index["+str(fib)+"]"))
            rib_vrf_addr = int(gdb.parse_and_eval(
                "((struct ipi_vrf*)"+str(vrf_addr)+")->proto"))

            if family == "i6":
                db = gdb.parse_and_eval(
                    "((struct rib_vrf*)"+str(rib_vrf_addr)+")->afi[2]->rib[1]")
            else:
                db = gdb.parse_and_eval(
                    "((struct rib_vrf*)"+str(rib_vrf_addr)+")->afi[1]->rib[1]")
            entry_type = "struct rib"
        elif table == "txlist":
            vr_addr = int(gdb.parse_and_eval(
                "rib_lib_globals->vr_vec.index[0]"))
            #vr_addr = gdb.Value(vr_addr).cast(util.ipi_vr_pointer)
            vrf_addr = int(gdb.parse_and_eval(
                "((struct ipi_vr*)"+str(vr_addr)+")->vrf_vec.index["+str(fib)+"]"))
            rib_vrf_addr = int(gdb.parse_and_eval(
                "((struct ipi_vrf*)"+str(vrf_addr)+")->proto"))

            if family == "i6":
                db = gdb.parse_and_eval(
                    "((struct rib_vrf*)"+str(rib_vrf_addr)+")->afi[2]->rib[1]->txlist")
            else:
                db = gdb.parse_and_eval(
                    "((struct rib_vrf*)"+str(rib_vrf_addr)+")->afi[1]->rib[1]->txlist")
            entry_type = "struct rib_ptree_node"
        elif table == "errlist":
            vr_addr = int(gdb.parse_and_eval(
                "rib_lib_globals->vr_vec.index[0]"))
            #vr_addr = gdb.Value(vr_addr).cast(util.ipi_vr_pointer)
            vrf_addr = int(gdb.parse_and_eval(
                "((struct ipi_vr*)"+str(vr_addr)+")->vrf_vec.index["+str(fib)+"]"))
            rib_vrf_addr = int(gdb.parse_and_eval(
                "((struct ipi_vrf*)"+str(vrf_addr)+")->proto"))

            if family == "i6":
                db = gdb.parse_and_eval(
                    "((struct rib_vrf*)"+str(rib_vrf_addr)+")->afi[2]->rib[1]->errlist")
            else:
                db = gdb.parse_and_eval(
                    "((struct rib_vrf*)"+str(rib_vrf_addr)+")->afi[1]->rib[1]->errlist")
            entry_type = "struct rib_ptree_node"
        elif table == "marker":
            vr_addr = int(gdb.parse_and_eval(
                "rib_lib_globals->vr_vec.index[0]"))
            #vr_addr = gdb.Value(vr_addr).cast(util.ipi_vr_pointer)
            vrf_addr = int(gdb.parse_and_eval(
                "((struct ipi_vr*)"+str(vr_addr)+")->vrf_vec.index["+str(fib)+"]"))
            rib_vrf_addr = int(gdb.parse_and_eval(
                "((struct ipi_vrf*)"+str(vrf_addr)+")->proto"))

            if family == "i6":
                db = gdb.parse_and_eval(
                    "((struct rib_vrf*)"+str(rib_vrf_addr)+")->afi[2]->rib[1]->marker")
            else:
                db = gdb.parse_and_eval(
                    "((struct rib_vrf*)"+str(rib_vrf_addr)+")->afi[1]->rib[1]->marker")
            entry_type = "struct rib_ptree_node"
        elif table == "interface":
            vr_addr = int(gdb.parse_and_eval(
                "rib_lib_globals->vr_vec.index[0]"))
            db = gdb.parse_and_eval(
                "((struct ipi_vr*)"+str(vr_addr)+")->ifm.if_table")
            entry_type = "struct interface"

        return (db, entry_type)

    def offsetof(type_name, member):
        type_pointer = gdb.lookup_type(type_name).pointer()
        zero = gdb.Value(0).cast(type_pointer)
        return int(zero[member].address)


class HexDumpAlign(gdb.Parameter):
    def __init__(self):
        super(HexDumpAlign, self).__init__('hxd-align',
                                           gdb.COMMAND_DATA,
                                           gdb.PARAM_BOOLEAN)

    set_doc = 'Determines if hxd always starts at an "aligned" address (see hxd-width'
    show_doc = 'Hex dump alignment is currently'


class HexDumpWidth(gdb.Parameter):
    def __init__(self):
        super(HexDumpWidth, self).__init__('hxd-width',
                                           gdb.COMMAND_DATA,
                                           gdb.PARAM_INTEGER)

    set_doc = 'Set the number of bytes per line of hxd'
    show_doc = 'The number of bytes per line in hxd is'


HexDumpAlign()
HexDumpWidth()


class PrefixPrinter:
    """Print struct ls_prefix and prefix."""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        try:
            prefix = self.val
            family = prefix['family']
            l = prefix['prefixlen']

            if 'prefix' in prefix.type.keys():
                p = prefix['prefix']
            else:
                p = prefix['u']['prefix']

            if family == 2:
                mem = gdb.selected_inferior().read_memory(p.address, 4)
                a = socket.inet_ntop(socket.AF_INET, mem)
                return ('%s/%d' % (a, l))
            elif family == 10:
                mem = gdb.selected_inferior().read_memory(p.address, 16)
                a = socket.inet_ntop(socket.AF_INET6, mem)
                return ('%s/%d' % (a, l))
            else:
                return ('%d:%s/%d' % (family, p, l))
        except Exception:
            exc_info = sys.exc_info()
            traceback.print_exception(*exc_info)


class MacPrinter:
    """Print mac address"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        mac = self.val
        return ("%02x:%02x:%02x:%02x:%02x:%02x" % (int(mac[0]), int(mac[1]), int(mac[2]), int(mac[3]), int(mac[4]), int(mac[5])))


class InAddrPrinter:
    """Print struct in_addr(ipv4)."""

    def __init__(self, val):
        self.val = val['s_addr']

    def to_string(self):
        n_addr = self.val
        mem = gdb.selected_inferior().read_memory(n_addr.address, n_addr.type.sizeof)
        return socket.inet_ntop(socket.AF_INET, mem)


class In6AddrPrinter:
    """Print struct in6_addr."""

    def __init__(self, val):
        self.val = val['__in6_u']['__u6_addr16']

    def to_string(self):
        n_addr = self.val
        mem = gdb.selected_inferior().read_memory(n_addr.address, n_addr.type.sizeof)
        return socket.inet_ntop(socket.AF_INET6, mem)


class HslIpv6AddressPrinter:
    """Print IPv6 address."""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        try:
            address = self.val.address
            length = self.val.type.sizeof
            mem = gdb.selected_inferior().read_memory(address, length)
            return '%s' % (socket.inet_ntop(socket.AF_INET6, mem))
        except:
            exc_info = sys.exc_info()
            traceback.print_exception(*exc_info)
            return


class HslIpAddressPrinter:
    """Print IP address."""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        try:
            address = self.val.address
            length = self.val.type.sizeof
            mem = gdb.selected_inferior().read_memory(address, length)
            return '%s' % (socket.inet_ntop(socket.AF_INET, mem))
        except:
            exc_info = sys.exc_info()
            traceback.print_exception(*exc_info)


class HslNexthopPrinter:
    """Print Hsl Nexthop address"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        nh = self.val
        mac = nh["mac"]
        flags = int(nh["flags"])
        nh_type = int(nh["nh_type"])
        ext_flags = int(nh["ext_flags"])
        system_info = int(nh["system_info"])
        refcnt = int(nh["refcnt"])
        rulecnt = int(nh["rulecnt"])
        retrycnt = int(nh["retrycnt"])
        aliveCounter = int(nh["aliveCounter"])
        out = "### HSL Nexthop Entry " + str(nh.address)+" ###\n"
        if int(nh['l2_ifp']) != 0:
            out += (' ifp %s, ifp2 %s nh_type 0x%x' % (str(nh["ifp"]["name"]).split(
                ',')[0], str(nh["l2_ifp"]["name"]).split(',')[0], nh_type))
        else:
            out += (' ifp %s, ifp2 %s nh_type 0x%x' %
                    (str(nh["ifp"]["name"]).split(',')[0], str(nh["l2_ifp"]), nh_type))
        out += (' mac %02x:%02x:%02x:%02x:%02x:%02x' %
                (int(mac[0]), int(mac[1]), int(mac[2]), int(mac[3]), int(mac[4]), int(mac[5])))
        out += (' flags 0x%04x, ext_flags: %04x refcnt %d, rulecnt: %d, retrycnt: %d aliveCounter: %d' %
                (flags, ext_flags, refcnt, rulecnt, retrycnt, aliveCounter))
        if system_info != 0:
            type_pointer = gdb.lookup_type(
                "struct hsl_bcm_nh_system_info").pointer()
            sysinfo = gdb.Value(system_info).cast(type_pointer).dereference()
            out += (' lport 0x%x egress_obj_id 0x%x' %
                    (int(sysinfo["lport"]), int(sysinfo["egress_obj_id"])))
        return out


class HslRouteNodePrinter:
    """Print Hsl Route Node"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        try:
            rn = self.val
            link0 = int(rn["link"][0])
            link1 = int(rn["link"][1])
            table = int(rn["table"])
            prefix = rn["p"]
            parent = int(rn["parent"])
            lock = int(rn["lock"])
            info_count = int(rn["info_count"])
            rn_type = int(rn["is_ecmp"])
            info = int(rn["info"])
            out = "### HSL Route Node " + str(rn.address)+" ###\n"
            out += (' p %s, type %d, lock %d, info 0x%x, info_count %d, table 0x%x, parent 0x%x, link[0] 0x%x, link[1] 0x%x' % (
                prefix, rn_type, lock, info, info_count, table, parent, link0, link1))
            return out
        except:
            exc_info = sys.exc_info()
            traceback.print_exception(*exc_info)


class HslPrefixEntryPrinter:
    """Print Hsl Prefix Entry"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        pe = self.val
        flags = int(pe["flags"])
        system_info = int(pe["system_info"])
        nhcount = int(pe["nhcount"])
        nhlist = int(pe["nhlist"])
        nhlist_hash = int(pe["nhlist_hash"])
        ifcount = int(pe["ifcount"])
        iflist = int(pe["iflist"])
        tun_list = int(pe["tun_list"])
        out = "### HSL Prefix Entry " + str(pe.address) + " ###\n"
        out += (' flags 0x%x, nhcount %d, ifcount %d, nhlist 0x%x, nhlist_hash 0x%x, iflist 0x%x, tun_list 0x%x, system_info 0x%x' %
                (flags, nhcount, ifcount, nhlist, nhlist_hash, iflist, tun_list, system_info))
        if system_info != 0:
            type_pointer = gdb.lookup_type(
                "struct hsl_bcm_mpath_sys_info").pointer()
            sysinfo = gdb.Value(system_info).cast(type_pointer).dereference()
            out += (' mpath_obj_id 0x%x egr_cnt %d' %
                    (int(system_info["mpath_obj_id"]), int(system_info["egr_cnt"])))
        if nhlist != 0:
            nhs = Listfy(pe["nhlist"], "struct hsl_nh_entry")
            out += '\n'
            for nh in nhs:
                if int(nh) != 0:
                    out += ('# %s \n' % nh.dereference())
        return out


class HslNhlistHashEntryPrinter:
    """Print Hsl Ecmp Group Entry"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        nhe = self.val
        num = int(nhe["num"])
        family = int(nhe["family"])
        ifids = nhe["ifindex"]
        sysinfo = int(nhe["system_info"])
        if family == socket.AF_INET6:
            nexthops = nhe["nexthopIP"]
        else:
            nexthops = nhe["nexthopIP6"]

        out = "#(ECMP) ifindex"
        for i in range(num):
            out += (' %d ' % (ifids[i]))

        for i in range(num):
            out += (' %s ' % (nexthos[i]))

        if system_info != 0:
            type_pointer = gdb.lookup_type(
                "struct hsl_bcm_mpath_sys_info").pointer()
            sysinfo = gdb.Value(system_info).cast(type_pointer).dereference()
            out += (' mpath_obj_id 0x%x egr_cnt %d' %
                    (int(system_info["mpath_obj_id"]), int(system_info["egr_cnt"])))
        return out


class HslIfPrinter:
    """Print Hsl Interface Entry"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        ifp = self.val
        out = '####[HSL Inteface ' + str(ifp.address) + ']### \n'
        out += ('fib_id: %d, name: %s, ifindex: %d, local_index: %d, type: %d, sub_type: %d, ' % (int(ifp['fib_id']),
                                                                                                  str(ifp["name"]).split(
                                                                                                      ',')[0],
                                                                                                  int(ifp['ifindex']), int(ifp['local_index']), int(ifp['type']), int(ifp['sub_type'])))
        out += ('operCnt: %d, flags: 0x%x, if_flags: 0x%x, pkt_flags: 0x%x, if_property: 0x%x, ' %
                (int(ifp['operCnt']), int(ifp['flags']), int(ifp['if_flags']), int(ifp['pkt_flags']), int(ifp['if_property'])))
        out += ('os_info: 0x%x, os_info_ext: 0x%x, system_info: 0x%x, ' %
                (int(ifp['os_info']), int(ifp['os_info_ext']), int(ifp['system_info'])))
        out += ('agg_member_cnt: %d, agg_ifindex: %d, ' %
                (int(ifp['agg_member_cnt']), int(ifp['agg_ifindex'])))
        out += ('if_vlan_tree: 0x%x' % (int(ifp['if_vlan_tree'])))
        out += ('slot_number: %d, if_number: %d, if_type: %d, if_media_type: %d, ' %
                (int(ifp['slot_number']), int(ifp['if_number']), int(ifp['if_type']), int(ifp['if_media_type'])))
        out += ('link_up_count: %d, link_down_count: %d, ' %
                (int(ifp['link_up_count']), int(ifp['link_down_count'])))
        for key in ifp.type.keys():
            if "type" in ifp.type.keys():
                if int(ifp['type']) == 2:  # HSL_IF_TYPE_IP
                    if key == 'u':
                        out += (' %s ' % ifp[key]["ip"])
                elif int(ifp['type']) == 3:  # HSL_IF_TYPE_L2_ETHERNET
                    if key == 'u':
                        out += (' %s ' % ifp[key]["l2_ethernet"])
                elif int(ifp['type']) == 4:  # HSL_IF_TYPE_MPLS
                    if key == 'u':
                        out += (' %s ' % ifp[key]["mpls"])
                else:
                    continue
        out += '####[HSL Inteface ' + str(ifp.address) + ' Finished]### \n'
        return out


class HslGeneralfPrinter:
    """Print Struct Entry"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        data = self.val
        out = ''
        for key in data.type.keys():
            if data[key].type.code == gdb.TYPE_CODE_STRUCT:
                out += ('\n%s: %s\n' % (key, str(data[key])))
            elif data[key].type.code == gdb.TYPE_CODE_ARRAY:
                out += ('\n %s {' % data[key].type)
                for i in range(data[key].type.range()[1]):
                    out += (' %s(%d) : ' % (key, i)) + \
                        '{},'.format(data[key][i])
                out += '}\n'
                #print("@@@@@ %d @@@@@" % data[key].type.sizeof)
            else:
                out += (' %s: ' % (key)) + '{}'.format(data[key])
        return out


class HslIgnorePrinter:
    """Ignore a type when do print"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        data = self.val
        out = ''
        return out


class HslInt64Printer:
    """Print Unsigned int64"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        data = self.val
        out = '%d' % (int(data['ll']))
        return out


class HslIpIfPrinter:
    """Print struct _hsl_ip_if hsl_ipv4If_t hsl_ipv6If_t"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        data = self.val
        out = ''
        for key in data.type.keys():
            if key == "ucAddr":
                uaddrs = Listfy(data['ucAddr'])
                for uaddr in uaddrs:
                    if int(uaddr) == 0:
                        continue
                    out += ('\n   <prefix %s, flags: %d, system_info: 0x%x>\n' %
                            (str(uaddr['prefix']), int(uaddr['flags']), int(uaddr['system_info'])))
            else:
                out += (' %s %s,' % (key, data[key]))
        return out


class ZebOSVectorPrinter:
    """Print zebos vector struct"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        data = self.val
        out = ''
        cap = int(data['max'])
        alloc = int(data['alloced'])
        index = data['index']
        for idx in range(alloc):
            out += (' %d: 0x%x ' % (idx, int(index[idx])))
        return out


class RibPtreeNodePrinter:
    """Print zebos rib_ptree_node struct"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        data = self.val
        tree = int(data['tree'])
        family = int(data['tree']["ptype"])
        link0 = int(data['link'][0])
        link1 = int(data['link'][1])
        parent = int(data['parent'])
        lock = int(data['lock'])
        info = int(data['info'])
        nhlist = data['nh_list']
        txlink = data['txlink']
        seq = int(data['seq'])
        req = int(data['req'])
        new_req = int(data['new_req'])
        fib_notify = int(data['fib_notify'])
        flags = int(data['flags'])
        key_len = data['key_len']

        out = '\n### RIBD Ptree Node ' + str(data.address)+'###\n'

        if family == socket.AF_INET:
            mem = gdb.selected_inferior().read_memory(
                int(data['key'].address), 4)
        if family == socket.AF_INET6:
            mem = gdb.selected_inferior().read_memory(
                int(data['key'].address), 16)
        a = socket.inet_ntop(family, mem)

        out += (' %s/%d, ' % (a, key_len))
        out += ('   \t\nlink0: 0x%x, link1: 0x%x, nhlist: 0x%x, parent: 0x%x, tree: 0x%x, info: 0x%x, ' %
                (link0, link1, int(nhlist), parent, tree, info))
        out += ('   \t\nseq: %d, req: %d, new_req: %d, fib_notify: %d, flags: 0x%x, ' %
                (seq, req, new_req, fib_notify, flags))

        out += ('   \t\nnh: ')
        if int(nhlist) != 0:
            nhiter = Listfy(nhlist)
            for nh in nhiter:
                if nh != None and int(nh) != 0:
                    out += (' %s, ' % nh)

        out += '\n### RIBD Ptree Node ' + str(data.address)+' Finshed ###\n'
        return out


class RibPrinter:
    """Print zebos rib struct"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        data = self.val
        out = '\n### RIBD rib ' + str(data.address) + '###\n'
        backup = []
        for key in data.type.keys():
            # for different distance route entry to same prefix,
            # use the first entry to get the list of other entry except the first one.
            if key == 'next':
                out += (' next: 0x%x, ' % int(data['next']))
                if int(data['next']) != 0 and int(data['prev']) == 0:
                    riter = Listfy(data['next'])
                    for r in riter:
                        if r != None and int(r) != 0:
                            backup.append(r.dereference())
            elif key == 'nexthop' and int(data['nexthop']) != 0:
                out += ' {} '.format(data['nexthop'].dereference())
            elif key == 'prev':
                out += (' prev: 0x%x, ' % int(data['prev']))
                if int(data['prev']) == 0:
                    type_pointer = gdb.lookup_type(
                        "struct rib_ptree_node").pointer()
                    node_pointer = gdb.Value(int(data.address) - util.offsetof('struct rib_ptree_node',
                                                                               "info")).cast(type_pointer)

                    out += ' \n@prefix@: {} \n'.format(
                        node_pointer.dereference())
            else:
                out += (' %s: ' % (key)) + '{}'.format(data[key])
        out += '\n### RIBD rib ' + str(data.address) + ' Finished###\n'
        for b in backup:
            out += '{}'.format(b)
        return out


class NexthopPrinter:
    """Print zebos nexthop struct"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        data = self.val
        out = '\n### Nexthop ' + str(data.address)+' ###\n'
        nexthops = []
        for key in data.type.keys():
            if key == 'next':
                out += (' next: 0x%x, ' % int(data['next']))
                if int(data['next']) != 0 and int(data['prev']) == 0:
                    riter = Listfy(data)
                    for r in riter:
                        if r != None and int(r) != 0:
                            nexthops.append(r.dereference())
            else:
                out += (' %s: %s, ' % (key, re.sub(",\s+'.*?' <repeats\s+\d+\s+times>",
                                                   '', str(data[key]))))
        out += '\n### Nexthop ' + str(data.address)+' Finished###\n'
        for nh in nexthops:
            out += '{}'.format(nh)
        return out


class InterfacePrinter:
    """Print zebos interface struct"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        data = self.val
        out = '\n### Interface ' + str(data.address) + '###\n'
        for key in data.type.keys():
            if key == 'next' and int(data['next']) != 0:
                riter = Listfy(data)
                for r in riter:
                    if r != None and int(r) != 0:
                        out += (' {} '.format(r.dereference()))
            else:
                out += (' %s: %s, ' % (key, re.sub(",\s+'.*?' <repeats\s+\d+\s+times>",
                                                   '', str(data[key]))))
        out += '\n### Interface ' + str(data.address) + ' Finished###\n'
        return out


class RibVrfPrinter(HslGeneralfPrinter):
    """Print rib_vrf struct"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        data = self.val
        out = '\n### RIB VRF ' + str(data.address)+' ###\n'
        out += super(RibVrfPrinter, self).to_string()
        out += '\n### RIB VRF ' + str(data.address)+' Finished###\n'
        return out


class RibPtreeTablePrinter:
    """Print rib_ptree_table struct"""

    def __init__(self, val):
        self.val = val

    def to_string(self):
        data = self.val
        out = '\n### RIB PTREE Table ' + str(data.address)+' ###\n'
        for key in data.type.keys():
            out += (' %s: %s, ' % (key, re.sub(",\s+'.*?' <repeats\s+\d+\s+times>",
                                               '', str(data[key]))))
        out += '\n### RIB PTREE Table ' + str(data.address)+' Finished###\n'
        return out


class DeleteBreakPoints(gdb.Command):
    """Remove all breakpoints"""

    def __init__(self):
        super(DeleteBreakPoints, self).__init__("cb", gdb.COMMAND_OBSCURE)

    def invoke(self, arguments, from_tty):
        bps = gdb.breakpoints()
        for bp in bps:
            print('breakpoint: %d:%s is deleted' % (bp.number, bp.location))
            bp.delete()


DeleteBreakPoints()


class ListBreakPoints(gdb.Command):
    """List all breakpoints"""

    def __init__(self):
        super(ListBreakPoints, self).__init__("lb", gdb.COMMAND_OBSCURE)

    def invoke(self, arguments, from_tty):
        bps = gdb.breakpoints()
        for bp in bps:
            print('breakpoint: %d, %s, %s, %s, %s ' %
                  (bp.number, bp.location, bp.type, bp.task, bp.thread))


ListBreakPoints()


class StackFold(gdb.Command):
    def __init__(self):
        super(StackFold, self).__init__("stackfold", gdb.COMMAND_DATA)

    def invoke(self, arg, from_tty):
        # An inferior is the 'currently running applications'. In this case we only have one.
        stack_maps = {}
        # This creates a dict where each element is keyed by backtrace.
        # Then each backtrace contains an array of "frames"
        #
        inferiors = gdb.inferiors()
        for inferior in inferiors:
            for thread in inferior.threads():
                # Change to our threads context
                thread.switch()
                # Get the thread IDS
                (tpid, lwpid, tid) = thread.ptid
                gtid = thread.num
                # Take a human readable copy of the backtrace, we'll need this for display later.
                o = gdb.execute('bt', to_string=True)
                # Build the backtrace for comparison
                backtrace = []
                gdb.newest_frame()
                cur_frame = gdb.selected_frame()
                while cur_frame is not None:
                    backtrace.append(cur_frame.name())
                    cur_frame = cur_frame.older()
                # Now we have a backtrace like ['pthread_cond_wait@@GLIBC_2.3.2', 'lazy_thread', 'start_thread', 'clone']
                # dicts can't use lists as keys because they are non-hashable, so we turn this into a string.
                # Remember, C functions can't have spaces in them ...
                s_backtrace = ' '.join(backtrace)
                # Let's see if it exists in the stack_maps
                if s_backtrace not in stack_maps:
                    stack_maps[s_backtrace] = []
                # Now lets add this thread to the map.
                stack_maps[s_backtrace].append(
                    {'gtid': gtid, 'tpid': tpid, 'bt': o})
        # Now at this point we have a dict of traces, and each trace has a "list" of pids that match. Let's display them
        for smap in stack_maps:
            # Get our human readable form out.
            o = stack_maps[smap][0]['bt']
            for t in stack_maps[smap]:
                # For each thread we recorded
                print("Thread %s (LWP %s))" % (t['gtid'], t['tpid']))
            print(o)


StackFold()

# this is x86 only


class FramePrinter:
    """Make ASCII art from a stack frame"""

    def __init__(self, frame):
        self._frame = frame
        self._decorator = gdb.FrameDecorator.FrameDecorator(self._frame)

    def __str__(self):
        if not self._frame.is_valid():
            return "<invalid>"
        result = ""
        # some basic frame stats
        if self._frame.function() is not None:
            result = result + "in " + self._frame.function().name
            if self._frame.type() is gdb.INLINE_FRAME:
                # recursively show inlining until we find a "real" parent frame
                result = result + "\ninlined with" + \
                    str(FramePrinter(self._frame.older()))
        else:
            result = result + "<unknown function>"
        if (self._frame.type() != gdb.NORMAL_FRAME):
            # IDK what else to do
            return result

        locls = self.__stackmap(self._decorator.frame_locals())
        args = self.__stackmap(self._decorator.frame_args())

        # assuming we are built with -fno-omit-frame-pointer here.  Not sure how to access
        # debug info that could tell us more, otherwise. More info is clearly present in C
        # (otherwise "info frame" could not do its job).

        # Display args
        yellow = "\u001b[33m"
        reset_color = "\u001b[0m"

        # find the address range of our args
        # from there to *(rbp+0x8), exclusive, is the range of possible args
        if args.keys():
            # the one with the highest address
            first_arg_addr = max(args.keys())
            result = result + self.__subframe_display(first_arg_addr,
                                                      self._frame.read_register(
                                                          'rbp')+0x8,
                                                      args,
                                                      yellow)

        # *(rbp+0x8) is the stored old IP
        cyan = "\u001b[36m"
        result = result + "\n" + \
            str(self._frame.read_register('rbp')+0x8) + " return address"
        voidstarstar = gdb.lookup_type("void").pointer().pointer()
        old_ip = (self._frame.read_register('rbp') +
                  0x8).cast(voidstarstar).dereference()
        result = result + cyan + " (" + str(old_ip) + ")" + reset_color

        # *(rbp) is the old RBP
        result = result + "\n" + \
            str(self._frame.read_register('rbp')+0x0) + " saved rbp"

        # print rest of stack, displaying locals
        green = "\u001b[32m"
        result = result + self.__subframe_display(self._frame.read_register('rbp')-0x8,
                                                  self._frame.read_register(
                                                      'sp')-0x8,
                                                  locls,
                                                  green)

        result = result + cyan + " <<< top of stack" + reset_color

        return result

    # display a range of stack addresses with colors, and compression of unknown contents as "stuff"
    def __subframe_display(self,
                           start, end,   # range of addresses to display
                           frame_items,  # map from addresses to lists of symbols
                           col):         # color to use for the symbols
        magenta = "\u001b[35m"
        reset_color = "\u001b[0m"
        empty_start = None
        result = ""
        for addr in range(start, end, -0x8):
            addr_hex = '0x{:02x}'.format(addr)
            if addr in frame_items:
                if empty_start:
                    # we just completed an empty range
                    if empty_start != (addr+0x8):
                        result = result + magenta + \
                            ' (through 0x{:02x})'.format(
                                addr+0x8) + reset_color
                    empty_start = None
                result = result + "\n" + addr_hex
                result = result + " " + col + \
                    ",".join([sym.name for sym in frame_items[addr]]
                             ) + reset_color
            elif empty_start is None:
                # we are starting an empty range
                empty_start = addr
                result = result + "\n" + addr_hex + magenta + " stuff" + reset_color

        if empty_start and (empty_start != end+0x8):
            # the empty range has more than one dword and extended through the end of the subframe
            result = result + magenta + \
                ' (through ' + str(end+0x8) + ')' + reset_color

        return result

    # produce a dict mapping addresses to symbol lists
    # for a given list of items (args or locals)

    def __stackmap(self, frame_items):
        symbolmap = defaultdict(list)
        if not frame_items:
            return symbolmap

        for i in frame_items:
            name = i.symbol().name
            addr = self._frame.read_var(name).address
            if not addr == None:
                # gdb.Value is not "hashable"; keys must be something else
                # so here we use addr converted to int
                sz = i.symbol().type.sizeof
                # mark all dwords in the stack with this symbol
                addr = addr.cast(gdb.lookup_type(
                    "void").pointer())  # cast to void*
                # handle sub-dword quantities by just listing everything that overlaps
                for saddr in range(addr, addr+sz, 0x8):
                    symbolmap[int(saddr)].append(i.symbol())
        return symbolmap

# Now create a gdb command that prints the current stack:


class PrintFrame (gdb.Command):
    """Display the stack memory layout for the current frame"""

    def __init__(self):
        super(PrintFrame, self).__init__("pframe", gdb.COMMAND_STACK)

    def invoke(self, arg, from_tty):
        try:
            print(FramePrinter(gdb.newest_frame()))
        except gdb.error:
            print("gdb got an error. Maybe we are not currently running?")


PrintFrame()


class BP(gdb.Breakpoint):
    def __init__(self):
        gdb.Breakpoint.__init__(self, "hsl_ifmgr_L2_link_down")
        self.silent = True

    def stop(self):
        print("##### breakpoint")

        gdb.execute("info locals")
        gdb.execute("print ifp->name")
        return False  # Do not the execution at this point


class PacketBreakPoints(gdb.Command):
    """Enable breakpoints for hsl packet process API"""
    hsl_bcm_pkt_process_bp = None
    hsl_bcm_pkt_send_bp = None

    def __init__(self):
        super(PacketBreakPoints, self).__init__("hpbs", gdb.COMMAND_STACK)

    def invoke(self, arguments, from_tty):
        PacketBreakPoints.hsl_bcm_pkt_process_bp = gdb.Breakpoint(
            "hsl_bcm_pkt_process")
        PacketBreakPoints.hsl_bcm_pkt_send_bp = gdb.Breakpoint(
            "hsl_bcm_pkt_send")
        gdb.events.stop.connect(PacketBreakPoints.hsl_bcm_pkt_stop_handler)

    def hsl_bcm_pkt_stop_handler(event):
        for b in event.breakpoints:
            if b == PacketBreakPoints.hsl_bcm_pkt_process_bp or b == PacketBreakPoints.hsl_bcm_pkt_send_bp:
                pkt = gdb.parse_and_eval("pkt")
                data = gdb.parse_and_eval("pkt->pkt_data[0].data")
                length = gdb.parse_and_eval("pkt->pkt_len")
                data = str(data).strip().split(' ')[0]
                gdb.execute("xs " + str(pkt) +
                            "/bcm_pkt_t/unit,cos,vlan,src_port")
                gdb.execute("xx " + str(data) + "/" + str(length))
            else:
                print("Unknown breakpoints: ", b.location)
        gdb.execute("continue")


PacketBreakPoints()


class RibBreakPoints(gdb.Command):
    """Enable breakpoints for rib route process API"""
    rib_process_route_ack_bp = None
    _rib_process_route_ack_bp = None
    rib_process_ipv6_route_ack_bp = None

    def __init__(self):
        super(RibBreakPoints, self).__init__("rrbs", gdb.COMMAND_STACK)

    def invoke(self, arguments, from_tty):
        RibBreakPoints.rib_process_route_ack_bp = gdb.Breakpoint(
            "rib_process_route_ack")
        RibBreakPoints._rib_process_route_ack_bp = gdb.Breakpoint(
            "_rib_process_route_ack")
        RibBreakPoints.rib_process_ipv6_route_ack_bp = gdb.Breakpoint(
            "rib_process_ipv6_route_ack")
        gdb.events.stop.connect(RibBreakPoints.ribd_rib_process_stop_handler)

    def ribd_rib_process_stop_handler(event):
        for b in event.breakpoints:
            if b == RibBreakPoints.rib_process_route_ack_bp:
                fib_id = int(gdb.parse_and_eval("fib_id"))
                seq_no = int(gdb.parse_and_eval("seq_no"))
                masklen = int(gdb.parse_and_eval("masklen"))
                resp = int(gdb.parse_and_eval("resp"))
                route_type = int(gdb.parse_and_eval("route_type"))
                addr = gdb.parse_and_eval("addr").dereference()
                print('fib_id: %d, seq_no: %d, masklen: %d, resp: %d, route_type: %d, addr: %s' %
                      (fib_id, seq_no, masklen, resp, route_type, str(addr)))
            elif b == RibBreakPoints._rib_process_route_ack_bp:
                vrf = gdb.parse_and_eval("vrf")
                rn = gdb.parse_and_eval("rn")
                resp = int(gdb.parse_and_eval("resp"))
                route_type = int(gdb.parse_and_eval("route_type"))
                print(' vrf: %s ' % str(vrf.dereference()))
                print(' rn: %s ' % str(rn.dereference()))
                print(' resp: %d, route_type: %d ' % (resp, route_type))
            elif b == RibBreakPoints.rib_process_ipv6_route_ack_bp:
                fib_id = int(gdb.parse_and_eval("fib_id"))
                seq_no = int(gdb.parse_and_eval("seq_no"))
                masklen = int(gdb.parse_and_eval("masklen"))
                resp = int(gdb.parse_and_eval("resp"))
                route_type = int(gdb.parse_and_eval("route_type"))
                addr = gdb.parse_and_eval("addr").dereference()
                print('fib_id: %d, seq_no: %d, masklen: %d, resp: %d, route_type: %d, addr: %s' %
                      (fib_id, seq_no, masklen, resp, route_type, str(addr)))
            else:
                print("Unknown breakpoints: ", b.location)
        gdb.execute("continue")


RibBreakPoints()


class NexthopBreakPoints(gdb.Command):
    """Enable breakpoints for hsl nexthop process API"""
    hsl_fib_nh_add_bp = None
    hsl_fib_nh_del_bp = None

    def __init__(self):
        super(NexthopBreakPoints, self).__init__("hnbs", gdb.COMMAND_STACK)

    def invoke(self, arguments, from_tty):
        NexthopBreakPoints.hsl_fib_nh_add_bp = gdb.Breakpoint(
            "_hsl_fib_nh_add")
        NexthopBreakPoints.hsl_fib_nh_del_bp = gdb.Breakpoint(
            "_hsl_fib_nh_delete")
        gdb.events.stop.connect(NexthopBreakPoints.hsl_fib_nh_stop_handler)

    def hsl_fib_nh_stop_handler(event):
        for b in event.breakpoints:
            if b == NexthopBreakPoints.hsl_fib_nh_add_bp:
                ifp = gdb.parse_and_eval("ifp")
                rnh = gdb.parse_and_eval("rnh")
                info = gdb.parse_and_eval("rnh->info")
                print("ifp : ", str(ifp), util.dump_struct(
                    str(ifp), "struct hsl_if", "name,ifindex,type"))
                print(rnh.dereference())
                #print("rnp : ", str(rnh), util.dump_struct(str(rnh), "struct hsl_route_node", "p,is_ecmp,type,info_count,info" ))
                if int(info) != 0:
                    print(gdb.parse_and_eval("*(struct hsl_nh_entry*)"+str(info)))
            elif b == NexthopBreakPoints.hsl_fib_nh_del_bp:
                rnh = gdb.parse_and_eval("rnh")
                info = int(gdb.parse_and_eval("rnh->info"))
                print(rnh.dereference())
                if info != 0:
                    type_pointer = gdb.lookup_type(
                        "struct hsl_nh_entry").pointer()
                    sysinfo = gdb.Value(info).cast(type_pointer).dereference()
                    print(sysinfo)
            else:
                print("Unknown breakpoints: ", b.location)
        gdb.execute("continue")


NexthopBreakPoints()


class RouteBreakPoints(gdb.Command):
    """Enable breakpoints for hsl route process API"""
    hsl_fib_add_to_hw_bp = None
    hsl_fib_delete_from_hw_bp = None

    def __init__(self):
        super(RouteBreakPoints, self).__init__("hrbs", gdb.COMMAND_STACK)

    def invoke(self, arguments, from_tty):
        RouteBreakPoints.hsl_fib_add_to_hw_bp = gdb.Breakpoint(
            "hsl_fib_add_to_hw")
        RouteBreakPoints.hsl_fib_delete_from_hw_bp = gdb.Breakpoint(
            "hsl_fib_delete_from_hw")
        gdb.events.stop.connect(RouteBreakPoints.hsl_fib_route_stop_handler)

    def hsl_fib_route_stop_handler(event):
        for b in event.breakpoints:
            if b == RouteBreakPoints.hsl_fib_add_to_hw_bp or b == RouteBreakPoints.hsl_fib_delete_from_hw_bp:
                rnp = gdb.parse_and_eval("rnp")
                nh = gdb.parse_and_eval("nh")
                info = int(gdb.parse_and_eval("rnp->info"))
                print(rnp.dereference())
                if info != 0:
                    type_pointer = gdb.lookup_type(
                        "struct hsl_prefix_entry").pointer()
                    sysinfo = gdb.Value(info).cast(type_pointer).dereference()
                    print(sysinfo)
                print(nh.dereference())
            else:
                print("Unknown breakpoints: ", b.location)
        gdb.execute("continue")


RouteBreakPoints()


class CurrentFrame(gdb.Command):
    """Dump current frame detailed information"""

    def __init__(self):
        super(CurrentFrame, self).__init__("cframe", gdb.COMMAND_STACK)

    def invoke(self, arguments, from_tty):
        try:
            print(gdb.newest_frame())
            print(gdb.selected_frame().name())
            print(gdb.selected_frame().type())
            print(gdb.selected_frame().pc())
            print(gdb.selected_frame().unwind_stop_reason())
            print(gdb.selected_frame().older())
            print(gdb.selected_frame().newer())
            print(gdb.selected_frame().function())
            # print(gdb.selected_frame().block())
            print(gdb.selected_frame().architecture())
            print(gdb.selected_frame().architecture().name())
            print(gdb.current_progspace())
            print(gdb.current_progspace().filename)
            print(gdb.current_progspace().pretty_printers)
            print(gdb.current_progspace().type_printers)
            print(gdb.current_progspace().frame_unwinders)
            print(gdb.current_progspace().frame_filters)
            print(gdb.inferiors())
            print(gdb.inferiors()[0].is_valid())
            for th in gdb.inferiors()[0].threads():
                print(th.name, th.num, th.global_num, th.ptid, th.is_exited(
                ), th.is_running(), th.is_stopped(), th.is_valid(), th.inferior)

            print(gdb.selected_frame().find_sal())
            symtab = gdb.selected_frame().find_sal().symtab
            if symtab != None:
                ltab = symtab.linetable()
                for line in ltab:
                    print(line.line, line.pc)
        except:
            exc_info = sys.exc_info()
            traceback.print_exception(*exc_info)


CurrentFrame()


class PythonDocument(gdb.Command):
    """Print the helper for GDB python API"""

    def __init__(self):
        super(PythonDocument, self).__init__("pd", gdb.COMMAND_STACK)

    def invoke(self, arguments, from_tty):
        print(help(gdb))


PythonDocument()


class DumpRaw(gdb.Command):
    """Dump memory in hex mode"""

    def __init__(self):
        super(DumpRaw, self).__init__("xx", gdb.COMMAND_STACK)

    def invoke(self, arguments, from_tty):
        args = re.split('/', arguments)
        if len(args) == 2:
            print(util.dump_raw(args[0], args[1]))
        else:
            print("Usuage: xx ADDERSS/LENGTH")


DumpRaw()


class DumpStruct(gdb.Command):
    """Dump a pointer as specified type"""

    def __init__(self):
        super(DumpStruct, self).__init__("xs", gdb.COMMAND_STACK)

    def invoke(self, arguments, from_tty):
        args = re.split('/', arguments)
        if len(args) == 2:
            print(util.dump_struct(args[0], args[1]))
        elif len(args) == 3:
            print(util.dump_struct(args[0], args[1], args[2]))
        else:
            print("Usuage: xs ADDERSS/TYPE[/FIELDS]")


DumpStruct()


class DumpHsl(gdb.Command):
    """Dump Hsl L3 related tables"""

    def __init__(self):
        super(DumpHsl, self).__init__("xh", gdb.COMMAND_STACK)

    def invoke(self, arguments, from_tty):
        arguments = arguments.strip()
        if not arguments.isspace() and len(arguments) >= 3:
            args = re.split('/', arguments)
            if args[0] != "i4" and args[0] != "i6":
                print(args[0], "is not a valid address family")
                print("Usuage: xh <i4/i6>/FIB_ID/<route|nexthop|interface>")
                return
            elif not args[1].isdigit():
                print(args[1], "is not a valid FIB ID")
                print("Usuage: xh <i4/i6>/FIB_ID/<route|nexthop|interface>")
                return
            elif not args[2] in ["route", "nexthop", "interface"]:
                print(args[2], "is not supported")
                print("Usuage: xh <i4/i6>/FIB_ID/<route|nexthop|interface>")
                return
            if len(args) == 3:
                self.dump_table(args[1], args[0], args[2])
            else:
                self.dump_table(args[1], args[0], args[2], args[3])
        else:
            print("Usuage: xh <i4/i6>/FIB_ID/<route|nexthop|interface>")

    def dump_table(self, fib, family, table, fields=None):
        db, type_name = util.get_hsl_db_and_type(fib, family, table)
        #dbiter = Listfy(db)
        #type_pointer = gdb.lookup_type(type_name).pointer()
        # for item in dbiter:
        #    if item != None:
        #        print(item.dereference())
        #        if int(item['info']) != 0:
        #            print(gdb.Value(int(item['info'])).cast(
        #                type_pointer).dereference())
        if fields == None:
            gdb.execute("xd " + "("+str(db.type)+")"+str(db) + "/" + type_name)
        else:
            gdb.execute("xd " + "("+str(db.type)+")"+str(db) +
                        "/" + type_name + "/" + fields)


DumpHsl()


class DumpRibd(gdb.Command):
    """Dump RIB related tables"""

    def __init__(self):
        super(DumpRibd, self).__init__("xr", gdb.COMMAND_STACK)

    def invoke(self, arguments, from_tty):
        try:
            arguments = arguments.strip()
            if not arguments.isspace() and len(arguments) >= 3:
                args = re.split('/', arguments)
                if args[0] != "i4" and args[0] != "i6":
                    print(args[0], "is not a valid address family")
                    print(
                        "Usuage: xr <i4/i6>/FIB_ID/<route|interface|txlist|errlist|marker|ptree>")
                    return
                elif not args[1].isdigit():
                    print(args[1], "is not a valid FIB ID")
                    print(
                        "Usuage: xr <i4/i6>/FIB_ID/<route|interface|txlist|errlist|marker|ptree>")
                    return
                elif not args[2] in ["route", "txlist", "errlist", "marker", "interface", "ptree"]:
                    print(args[2], "is not supported")
                    print(
                        "Usuage: xr <i4/i6>/FIB_ID/<route|interface|txlist|errlist|marker|ptree>")
                    return
                if len(args) == 3:
                    self.dump_table(args[1], args[0], args[2])
                else:
                    self.dump_table(args[1], args[0], args[2], args[3])
            else:
                print(
                    "Usuage: xr <i4/i6>/FIB_ID/<route|interface|txlist|errlist|marker|ptree>")
        except:
            exc_info = sys.exc_info()
            traceback.print_exception(*exc_info)

    def dump_table(self, fib, family, table, fields=None):
        db, type_name = util.get_ribd_db_and_type(fib, family, table)
        dbiter = Listfy(db)
        type_pointer = gdb.lookup_type(type_name).pointer()
        for rib in dbiter:
            if rib.type != util.fifo_pointer:
                break

            if rib != None and int(rib) != 0:
                if rib.type == util.fifo_pointer and int(rib['next'].address) != int(db.address):
                    print(gdb.Value(int(rib.address) - util.offsetof(type_name,
                                                                     "txlink")).cast(type_pointer).dereference())
                elif rib.type == util.fifo_pointer and int(rib['next'].address) == int(db.address):
                    break

        if rib.type != util.fifo_pointer:
            if table == 'ptree':
                dbiter = Listfy(db)
                for rib in dbiter:
                    if rib != None and int(rib) != 0:
                        print(rib.dereference())
            elif fields == None:
                gdb.execute("xd " + "("+str(db.type)+")" +
                            str(db) + "/" + type_name)
            else:
                gdb.execute("xd " + "("+str(db.type)+")" +
                            str(db) + "/" + type_name + "/" + fields)

        # if int(rib['info']) != 0:
        #    print(gdb.Value(int(rib['info'])).cast(
        #        type_pointer).dereference())


DumpRibd()

printers = gdb.printing.RegexpCollectionPrettyPrinter('alpha')
printers.add_printer('hsl_prefix_t', '^hsl_prefix_t$', PrefixPrinter)
printers.add_printer('hsl_mac_address_t', '^hsl_mac_address_t$', MacPrinter)
printers.add_printer('hsl_nh_entry', '^.*hsl_nh_entry$', HslNexthopPrinter)
printers.add_printer(
    'hsl_route_node', '^.*hsl_route_node$', HslRouteNodePrinter)
printers.add_printer('hsl_prefix_entry',
                     '^.*hsl_prefix_entry$', HslPrefixEntryPrinter)
printers.add_printer('struct in_addr', '^.*in_addr$', InAddrPrinter)
printers.add_printer('struct in6_addr', '^.*in6_addr$', In6AddrPrinter)
printers.add_printer('hsl_ipv4Address_t',
                     '^.*hsl_ipv4Address_t$', HslIpAddressPrinter)
printers.add_printer('hsl_ipv6Address_t',
                     '^.*hsl_ipv6Address_t$', HslIpv6AddressPrinter)
printers.add_printer('hsl_if', '^.*hsl_if$', HslIfPrinter)
printers.add_printer('hsl_ifIP_t', '^.*hsl_ifIP_t$', HslGeneralfPrinter)
printers.add_printer('hsl_ifL2_ethernet_t',
                     '^.*hsl_ifL2_ethernet_t$', HslGeneralfPrinter)
printers.add_printer('hsl_ifMPLS_t', '^.*hsl_ifMPLS_t$', HslGeneralfPrinter)
printers.add_printer('pal_if_stats', '^.*pal_if_stats$', HslGeneralfPrinter)
printers.add_printer('if_stats', '^.*if_stats$', HslGeneralfPrinter)
printers.add_printer('arp_params', '^.*arp_params$', HslGeneralfPrinter)
printers.add_printer('label_space_data',
                     '^.*label_space_data$', HslGeneralfPrinter)
printers.add_printer('label_range_data',
                     '^.*label_range_data$', HslGeneralfPrinter)
printers.add_printer('hsl_efm_err_frame',
                     '^.*hsl_efm_err_frame$', HslIgnorePrinter)
printers.add_printer(
    'hal_if_counters', '^.*hal_if_counters$', HslIgnorePrinter)
printers.add_printer('ut_int64_t', '^.*ut_int64_t$', HslInt64Printer)
printers.add_printer('st_int64_t', '^.*st_int64_t$', HslInt64Printer)
printers.add_printer('hsl_ip_if', '^.*hsl_ip_if$', HslIpIfPrinter)
printers.add_printer('hsl_IfIPv4_t', '^.*hsl_IfIPv4_t$', HslIpIfPrinter)
printers.add_printer('hsl_IfIPv6_t', '^.*hsl_IfIPv6_t$', HslIpIfPrinter)
printers.add_printer('hsl_nhlist_hash_entry',
                     '^.*hsl_nhlist_hash_entry$', HslNhlistHashEntryPrinter)
printers.add_printer('vector', '^vector$', ZebOSVectorPrinter)
printers.add_printer('struct _vector', '^struct _vector$', ZebOSVectorPrinter)
printers.add_printer(
    'rib_ptree_node', '^.*rib_ptree_node$', RibPtreeNodePrinter)
printers.add_printer('struct rib', '^rib$', RibPrinter)
printers.add_printer('struct nexthop', '^.*nexthop$', NexthopPrinter)
printers.add_printer('struct interface', '^.*interface$', InterfacePrinter)
printers.add_printer('struct rib_vrf', '^.*rib_vrf$', RibVrfPrinter)
printers.add_printer('struct rib_ptree_table',
                     '^.*rib_ptree_table$', RibPtreeTablePrinter)
printers.add_printer('struct counter', '^.*counter$', HslGeneralfPrinter)
printers.add_printer('struct rib_vrf_afi',
                     '^.*rib_vrf_afi$', HslGeneralfPrinter)
gdb.printing.register_pretty_printer(gdb.current_objfile(), printers)

gdb.execute("set print pretty")
gdb.execute("set pagination off")
gdb.execute("set height 0")
