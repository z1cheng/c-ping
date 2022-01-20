# c-ping
c-ping is a very simple and small ping tool that sends ICMP Echo datagram to a host.

# ICMP Echo or Echo reply message

The structure of ICMP Echo or Echo reply message:

```
  0               1               2               3
  0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 8 9 A B C D E F
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |     Type      |     Code      |          Checksum             |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |           Identifier          |        Sequence Number        |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |  	                    Data(Optional)...                    |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

```

1. Type: **8** for echo message, and **0** for echo reply message
2. Code: for echo or echo reply, it should be **0**
3. Checksum: 16 bits, it will be computed by the checksum algorithm
4. Identifier: may be used like a port in TCP or UDP to **identify a session**
5. Sequence Number: may be **incremented** on each echo request sent

The destination returns the **same** Identifier and Sequence Number in the reply.


# Checksum algorithm
This algorithm based on [RFC1071](https://datatracker.ietf.org/doc/html/rfc1071) is below:
1. Set checksum field of request to 0
2. Add each 16 bits from request to sum
3. Add the remaining byte if it’s odd
4. Add the high 16 bits to the low 16 bits util the high 16 bits are 0
5. The complement of sum is the checksum

To check a checksum, just pass the datagram with the checksum to this method, then check if the result is 0 

# Example
Ensure that you have installed the go environment first, then just run following commands:

```shell
# generate c-ping
go build ping.go

# start to ping
./ping  github.com
```
