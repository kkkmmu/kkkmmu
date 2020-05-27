Virtual eXtensible Local Area Network (VXLAN): 
VXLAN creates LAN segments using a MAC in IP encapsulation. The encapsulation carries the original L2
frame received from a host to the destination in another server using IP tunnels. The endpoints of the virtualized tunnel
formed using VXLAN are called VTEPs (VXLAN Tunnel EndPoints). This technology allows the network to support
several tenants with minimum changes in the network. The VTEPs carry tenant data in L3 tunnels over the network. The
tenant data is not used in routing or switching. This aids in tenant machine movement and allows the tenants to have
the same IP or MAC addresses on end devices, hosts/VM’s.

VNI VXLAN Network Identifier (or VXLAN Segment ID)
VTEP VXLAN Tunnel End Point. An entity that originates and/or terminates VXLAN tunnels
VXLAN Virtual eXtensible Local Area Network
VXLAN Segment VXLAN Layer 2 overlay network over which VMs communicate
VXLAN Gateway An entity that forwards traffic between VXLANs
