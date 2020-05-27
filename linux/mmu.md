Page and Buffer Cache

Performance and efficiency are two factors to which great importance is attached during kernel
development. The kernel relies not only on a sophisticated overall concept of interaction between
its individual components, but also on an extensive framework of buffers and caches designed to
boost system speed.

Buffering and caching make use of parts of system RAM to ensure that the most important and the
most frequently used data of block devices can be manipulated not on the slow devices themselves
but in main memory. RAM memory is also used to store the data read in from block devices so that
the data can be subsequently accessed directly in fast RAM when it is needed again rather than
fetching it from external devices.

Of course, this is done transparently so that the applications do not and cannot notice any difference
as to from where the data originate.

Data are not written back after each change but after a specific interval whose length depends
on a variety of factors such as the free RAM capacity, the frequency of usage of the data held in
RAM, and so on. Individual write requests are bundled and collectively take less time to perform.
Consequently, delaying write operations improves system performance as a whole.
However, caching has its downside and must be employed judiciously by the kernel:
	Usually there is far less RAM capacity than block device capacity so that only carefully
	selected data may be cached.

	The memory areas used for caching are not exclusively reserved for ‘‘normal‘‘ application
	data. This reduces the RAM capacity that is effectively available.

	If the system crashes (owing to a power outage, e.g.), the caches may contain data that have
	not been written back to the underlying block device. Such data are irretrievably lost.
	However, the advantages of caching outweigh the disadvantages to such an extent that caches are
	permanently integrated into the kernel structures.
