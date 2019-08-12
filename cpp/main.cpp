#include <pcap.h>
#include <time.h>
#include <stdlib.h>
#include <stdio.h>
#include <netinet/in.h>
#include <ctype.h>
#include <arpa/inet.h>
#include <map>
#include <iostream>
#include <pthread.h>
#include <unistd.h>
#include "pcap.h"
#include "net.h"

using namespace std;

/* max size packet to catch */
#define MAX_TCP_PACKET_SIZE 65535
/* time to wait to return packets */
#define TIME_DURATION 1000 * 3
string FILTER_STRING = "dst port 6379";
/* ethernet headers are always exactly 14 bytes [1] */
string DEVICE = "eth0";
unsigned int TIME = 60;
int main(int argc, char *argv[])
{
  int ch;
  while ((ch = getopt(argc, argv, "d:f:t:")) != -1)
  {
    switch (ch)
    {
    case 'd': // set device
      DEVICE = optarg;
      break;
    case 'f': // set fukter string
      FILTER_STRING = optarg;
      break;
    case 't': // set time to clear dat
      TIME = (unsigned int)atoi(optarg);
      break;
    default:
      break;
    }
  }
  initMethod();
  pthread_t pClean, pTcpService;
  if (pthread_create(&pClean, NULL, clean, NULL) != 0)
  {
    printf("create pthread error!\n");
    exit(1);
  }
  if (pthread_create(&pTcpService, NULL, tcpInit, NULL) != 0)
  {
    printf("create pthread error!\n");
    exit(1);
  }
  char errBuf[PCAP_ERRBUF_SIZE];
  pcap_t *device = pcap_open_live(DEVICE.c_str(), MAX_TCP_PACKET_SIZE, 1, TIME_DURATION, errBuf);

  if (!device)
  {
    printf("错误: pcap_open_live(): %s\n", errBuf);
    exit(1);
  }
  struct bpf_program filter;
  pcap_compile(device, &filter, FILTER_STRING.c_str(), 1, 0);
  pcap_setfilter(device, &filter);

  pcap_loop(device, -1, getPacket, NULL);
  pcap_close(device);

  return 0;
}