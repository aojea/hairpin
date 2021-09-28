# hairpin

Hairpin creates a synchronous, in-memory, half-duplex connection.

Hairpin implements golang net.Conn interface and is very useful for protocol testing.

PacketHairpin implements golang net.PacketConn interface.


## How it works

The data written in the connection is executed by the Packet Process Handler and written to a buffer.

If there is no Packet Process Handler defined, the data is written directly to the buffer.

Operations are serialized, once a packet is processed, Writes() are blocked until the processed packet is Read().

Partial Reads are allowed, but Writes() will not be unblocked until the Read() buffer is fully drained.

