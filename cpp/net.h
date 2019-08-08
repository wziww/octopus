#ifndef __OCTOPUS_NET_h
#define __OCTOPUS_NET_h
#include <unistd.h>
#include <iostream>
#include <sys/socket.h>
#include <arpa/inet.h>

#define MAXSIZE 1024
#define PORT 9712 // TCP service port
void *tcpInit(void *ptr);
class tcp_server
{
private:
  int socket_fd;
  sockaddr_in myserver;
  sockaddr_in remote_addr;

public:
  tcp_server(int listen_port);
  int recv_msg();
};
#endif