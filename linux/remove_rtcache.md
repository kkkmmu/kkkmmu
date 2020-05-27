author	David S. Miller <davem@davemloft.net>	2012-07-22 17:04:15 -0700
committer	David S. Miller <davem@davemloft.net>	2012-07-22 17:04:15 -0700
commit	5e9965c15ba88319500284e590733f4a4629a288 (patch)
tree	ab76263b9f43fb75048a50141d199f445f5fdd2d
parent	3ba97381343b271296487bf073eb670d5465a8b8 (diff)
parent	2860583fe840d972573363dfa190b2149a604534 (diff)
download	net-next-5e9965c15ba88319500284e590733f4a4629a288.tar.gz

Merge branch 'kill_rtcache'
The ipv4 routing cache is non-deterministic, performance wise, and is
subject to reasonably easy to launch denial of service attacks.

The routing cache works great for well behaved traffic, and the world
was a much friendlier place when the tradeoffs that led to the routing
cache's design were considered.

What it boils down to is that the performance of the routing cache is
a product of the traffic patterns seen by a system rather than being a
product of the contents of the routing tables.  The former of which is
controllable by external entitites.

Even for "well behaved" legitimate traffic, high volume sites can see
hit rates in the routing cache of only ~%10.

The general flow of this patch series is that first the routing cache
is removed.  We build a completely new rtable entry every lookup
request.

Next we make some simplifications due to the fact that removing the
routing cache causes several members of struct rtable to become no
longer necessary.

Then we need to make some amends such that we can legally cache
pre-constructed routes in the FIB nexthops.  Firstly, we need to
invalidate routes which are hit with nexthop exceptions.  Secondly we
have to change the semantics of rt->rt_gateway such that zero means
that the destination is on-link and non-zero otherwise.

Now that the preparations are ready, we start caching precomputed
routes in the FIB nexthops.  Output and input routes need different
kinds of care when determining if we can legally do such caching or
not.  The details are in the commit log messages for those changes.

The patch series then winds down with some more struct rtable
simplifications and other tidy ups that remove unnecessary overhead.

On a SPARC-T3 output route lookups are ~876 cycles.  Input route
lookups are ~1169 cycles with rpfilter disabled, and about ~1468
cycles with rpfilter enabled.

These measurements were taken with the kbench_mod test module in the
net_test_tools GIT tree:

git://git.kernel.org/pub/scm/linux/kernel/git/davem/net_test_tools.git

That GIT tree also includes a udpflood tester tool and stresses
route lookups on packet output.

For example, on the same SPARC-T3 system we can run:

	time ./udpflood -l 10000000 10.2.2.11

with routing cache:
real    1m21.955s       user    0m6.530s        sys     1m15.390s

without routing cache:
real    1m31.678s       user    0m6.520s        sys     1m25.140s

Performance undoubtedly can easily be improved further.

For example fib_table_lookup() performs a lot of excessive
computations with all the masking and shifting, some of it
conditionalized to deal with edge cases.

Also, Eric's no-ref optimization for input route lookups can be
re-instated for the FIB nexthop caching code path.  I would be really
pleased if someone would work on that.

In fact anyone suitable motivated can just fire up perf on the loading
of the test net_test_tools benchmark kernel module.  I spend much of
my time going:

bash# perf record insmod ./kbench_mod.ko dst=172.30.42.22 src=74.128.0.1 iif=2
bash# perf report

Thanks to helpful feedback from Joe Perches, Eric Dumazet, Ben
Hutchings, and others.

Signed-off-by: David S. Miller <davem@davemloft.net>
Diffstat
-rw-r--r--	include/net/dst.h	15	
-rw-r--r--	include/net/flow.h	1	
-rw-r--r--	include/net/inet_connection_sock.h	3	
-rw-r--r--	include/net/ip_fib.h	3	
-rw-r--r--	include/net/route.h	40	
-rw-r--r--	net/core/dst.c	4	
-rw-r--r--	net/dccp/ipv4.c	2	
-rw-r--r--	net/decnet/dn_route.c	4	
-rw-r--r--	net/ipv4/arp.c	5	
-rw-r--r--	net/ipv4/fib_frontend.c	5	
-rw-r--r--	net/ipv4/fib_semantics.c	4	
-rw-r--r--	net/ipv4/inet_connection_sock.c	9	
-rw-r--r--	net/ipv4/ip_fragment.c	4	
-rw-r--r--	net/ipv4/ip_gre.c	2	
-rw-r--r--	net/ipv4/ip_input.c	4	
-rw-r--r--	net/ipv4/ip_output.c	2	
-rw-r--r--	net/ipv4/ipip.c	2	
-rw-r--r--	net/ipv4/ipmr.c	9	
-rw-r--r--	net/ipv4/netfilter/ipt_MASQUERADE.c	5	
-rw-r--r--	net/ipv4/route.c	1329	
-rw-r--r--	net/ipv4/tcp_ipv4.c	4	
-rw-r--r--	net/ipv4/xfrm4_input.c	4	
-rw-r--r--	net/ipv4/xfrm4_policy.c	9	
-rw-r--r--	net/ipv6/route.c	4	
-rw-r--r--	net/sctp/transport.c	2	
-rw-r--r--	net/xfrm/xfrm_policy.c	23	
26 files changed, 292 insertions, 1206 deletions

