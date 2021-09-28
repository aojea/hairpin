# hairpin

Hairpin creates a synchronous, in-memory, half-duplex connection, and allows to set
a custom function to process the data that is going through the connection.

The data written in the connection is executed by the custom Handler and sent back to the connection.
Operations are serialized, once a packet is processed, Writes() are blocked until the processed packet is Read().

Partial Reads are allowed, but Writes() will not be unblocked until the Read() buffer is fully drained.

Hairpin implements golang net.Conn interface and is very useful for protocol testing.

## Example

Implement an in-memory DNS resolver using a custom handler that is able to process DNS requests

https://github.com/aojea/mem-resolver