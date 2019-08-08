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
void *handleAccept(void *ptr)
{
  int fd = *(int *)ptr;
  while (true)
  {
    char buffer[MAXSIZE];
    memset(buffer, 0, MAXSIZE);
    if ((read(fd, buffer, MAXSIZE)) < 0)
    {
      printf("%s\n", strerror(errno));
      exit(1);
    }
    else
    {
      printf("%s\n", string(buffer).c_str());
      if (string(buffer) == "get\r\n")
      {
        cmd_mutex.lock();
        for (map<string, int>::reverse_iterator iter = cmdCount.rbegin(); iter != cmdCount.rend(); iter++)
        {
          write(fd, (void *)&iter->first, sizeof(iter->first));
          write(fd, CRLF, sizeof(CRLF));
          string count = to_string((long long int)iter->second);
          write(fd, (void *)&count, sizeof(count));
          write(fd, CRLF, sizeof(CRLF));
        }
        cmd_mutex.unlock();
      }
      if (string(buffer) == "quit\r\n")
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