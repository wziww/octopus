#include "net.h"
#include <string.h>
#include <mutex>
#include <errno.h>
#include <iostream>
#include "pcap.h"

#define PORT 9712 // TCP service port
#define CRLF "\r\n"

using std::string;
using std::to_string;

static int CLIENT_COUNT = 0;
const int CLIENT_MAX_COUNT = 2;
static std::mutex client_mutex;

void *tcpInit(void *ptr)
{
  tcp_server ts(PORT);
  ts.recv_msg();
  pthread_exit(0);
}
tcp_server::tcp_server(int listen_port)
{

  if ((socket_fd = socket(PF_INET, SOCK_STREAM, IPPROTO_TCP)) < 0)
  {
    printf("%s\n", strerror(errno));
    exit(1);
  }

  memset(&myserver, 0, sizeof(myserver));
  myserver.sin_family = AF_INET;
  myserver.sin_addr.s_addr = htonl(INADDR_ANY);
  myserver.sin_port = htons(listen_port);

  if (::bind(socket_fd, (sockaddr *)&myserver, sizeof(myserver)) < 0)
  {
    printf("%s\n", strerror(errno));
    exit(1);
  }

  if (listen(socket_fd, 10) < 0)
  {
    printf("%s\n", strerror(errno));
    exit(1);
  }
};
void _write(int __fd, const void *__buf, size_t __nbyte)
{
  if (write(__fd, __buf, __nbyte) <= 0)
  {
    printf("write error\n");
  }
}
void *handleAccept(void *ptr)
{
  int fd = *(int *)ptr;
  while (true)
  {
    char buffer[MAXSIZE];
    memset(buffer, 0, MAXSIZE);
    int len = (read(fd, buffer, MAXSIZE));
    if (len < 0)
    {
      printf("%s\n", strerror(errno));
      break;
    }
    else if (len == 0)
    {
      break;
    }
    else
    {
      if (string(buffer) == "ping\r\n")
      {
        _write(fd, "pong", 4);
        _write(fd, CRLF, sizeof(CRLF));
      }
      else if (string(buffer) == "get\r\n")
      {
        cmd_mutex.lock();
        for (map<string, int>::reverse_iterator iter = cmdCount.rbegin(); iter != cmdCount.rend(); iter++)
        {
          char buffer[sizeof(iter->first)];
          char v_buffer[sizeof(iter->second)];
          sprintf(v_buffer, "%d", iter->second);
          for (int i = 0; i < sizeof(iter->first); i++)
          {
            buffer[i] = iter->first[i];
          }
          _write(fd, buffer, sizeof(buffer));
          _write(fd, CRLF, sizeof(CRLF));
          _write(fd, v_buffer, sizeof(v_buffer));
          _write(fd, CRLF, sizeof(CRLF));
        }
        cmd_mutex.unlock();
      }
      else if (string(buffer) == "quit\r\n")
      {
        break;
      }
    }
  }
  close(fd);
  client_mutex.lock();
  CLIENT_COUNT--;
  client_mutex.unlock();
  pthread_exit(0);
};
int tcp_server::recv_msg()
{
  while (1)
  {
    int fd;
    pthread_t current;
    socklen_t sin_size = sizeof(struct sockaddr_in);
    if ((fd = accept(socket_fd, (struct sockaddr *)&remote_addr, &sin_size)) == -1)
    {
      printf("Accept error!\n");
      continue;
    }
    if (CLIENT_COUNT > CLIENT_MAX_COUNT)
    {
      close(fd);
      continue;
    }
    client_mutex.lock();
    CLIENT_COUNT++;
    client_mutex.unlock();
    printf("Received a connection from %s\n", (char *)inet_ntoa(remote_addr.sin_addr));
    if (pthread_create(&current, NULL, handleAccept, &fd) != 0)
    {
    }
  }
  return 0;
};