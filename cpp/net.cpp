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
  int flag = 1;
  if (setsockopt(socket_fd, SOL_SOCKET, SO_REUSEADDR, &flag, sizeof(flag)) < 0)
  {
    printf("socket setsockopt error=%d(%s)!!!\n", errno, strerror(errno));
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
      string _cmd = string(buffer);
      if (_cmd == "ping\r\n")
      {
        _write(fd, "pong", 4);
        _write(fd, CRLF, sizeof(CRLF));
      }
      else if (_cmd == "get\r\n")
      {
        cmd_mutex.lock();
        long long int mlen = cmdCount.size();
        string mlenstr = to_string(mlen);
        char len_buffer[mlenstr.size()];
        strcpy(len_buffer, mlenstr.c_str());
        _write(fd, len_buffer, sizeof(len_buffer));
        _write(fd, CRLF, sizeof(CRLF));
        for (map<string, int>::reverse_iterator iter = cmdCount.rbegin(); iter != cmdCount.rend(); iter++)
        {
          char v_buffer[sizeof(iter->second)];
          memset(v_buffer, 0, sizeof(v_buffer));
          sprintf(v_buffer, "%d", iter->second);
          _write(fd, iter->first.c_str(), iter->first.size());
          _write(fd, CRLF, sizeof(CRLF));
          _write(fd, v_buffer, sizeof(v_buffer));
          _write(fd, CRLF, sizeof(CRLF));
        }
        cmd_mutex.unlock();
      }
      else if (_cmd == "quit\r\n")
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