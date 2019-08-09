#include <pcap.h>
#include <time.h>
#include <stdlib.h>
#include <stdio.h>
#include <netinet/in.h>
#include <ctype.h>
#include <arpa/inet.h>
#include <map>
#include <iostream>
#include <mutex>
#include <pthread.h>
#include <unistd.h>
#include <string.h>
#include "net.h"
#include "pcap.h"

using namespace std;

struct pcap_pkthdr *header;
std::mutex cmd_mutex;
map<string, int> cmdCount;
map<string, bool> validMethod;
void initMethod()
{
  validMethod["get"] = true;
  validMethod["set"] = true;
  validMethod["cluster"] = true;
  validMethod["del"] = true;
  validMethod["append"] = true;
  validMethod["hgetall"] = true;
  validMethod["zadd"] = true;
  validMethod["info"] = true;
}
/*
 * redis 方法统计
 */
void count_cmd(const u_char *payload, int len)
{
  cmd_mutex.lock();
  int i;
  const u_char *ch;
  ch = payload;
  char cmd[100];
  memset(cmd, 0, sizeof(cmd));
  int index = 0;
  if (len < 8 || len > 8 + 50)
    return;
  ch += 8;
  for (i = 7; i < len; i++)
  {
    if (isprint(*ch) && index < 100)
    {
      cmd[index] = ((char)*ch <= 'Z' && (char)*ch >= 'A') ? (char)*ch - ('Z' - 'z') : (char)*ch; // to lower
      index++;
    }
    else
    {
      break;
    }
    ch++;
  }
  if (validMethod[string(cmd)])
  {
    cmdCount[string(cmd)] = cmdCount[string(cmd)] + 1;
  }
  cmd_mutex.unlock();
  return;
}
/**
 * packet 处理函数
 */
void getPacket(u_char *argument, const struct pcap_pkthdr *packet_header,
               const u_char *packet)
{

  /* declare pointers to packet headers */
  const struct sniff_ethernet *ethernet; /* The ethernet header [1] */
  const struct ip_struct *ip;            /* The IP header */
  const struct tcp_struct *tcp;          /* The TCP header */
  const u_char *payload;                 /* Packet payload */

  int size_ip;
  int size_tcp;
  int size_payload;
  /* define ethernet header */
  ethernet = (struct sniff_ethernet *)(packet);

  /* define/compute ip header offset */
  ip = (struct ip_struct *)(packet + SIZE_ETHERNET);
  size_ip = IP_HL(ip) * 4;
  if (size_ip < 20)
  {
    printf("   * Invalid IP header length: %u bytes\n", size_ip);
    return;
  }
  /*
   *  This packet is TCP.
   */
  if (ip->ip_p != IPPROTO_TCP)
  {
    printf("   * Invalid TCP packet\n");
    return;
  }
  /* define/compute tcp header offset */
  tcp = (struct tcp_struct *)(packet + SIZE_ETHERNET + size_ip);
  size_tcp = TH_OFF(tcp) * 4;
  if (size_tcp < 20)
  {
    printf("   * Invalid TCP header length: %u bytes\n", size_tcp);
    return;
  }
  // printf("   Src port: %d\n", ntohs(tcp->th_sport));
  // printf("   Dst port: %d\n", ntohs(tcp->th_dport));
  /* define/compute tcp payload (segment) offset */
  payload = (u_char *)(packet + SIZE_ETHERNET + size_ip + size_tcp);

  /* compute tcp payload (segment) size */
  size_payload = ntohs(ip->ip_len) - (size_ip + size_tcp);

  /*
       * Print payload data; it might be binary, so don't just
       * treat it as a string.
       */
  if (size_payload > 0)
  {
    // printf("   Payload (%d bytes):\n", size_payload);
    count_cmd(payload, size_payload);
  }
  return;
}
/**
 * 重置
 */
void *clean(void *ptr)
{
  for (;;)
  {
    cmd_mutex.lock();
    for (map<string, int>::reverse_iterator iter = cmdCount.rbegin(); iter != cmdCount.rend(); iter++)
      cmdCount[iter->first] = 0;
    cmd_mutex.unlock();
    sleep(60); // 每分钟清除一次
  }
}
