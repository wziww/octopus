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
using namespace std;
std::mutex cmd_mutex;
map<string, int> cmdCount;
map<string, bool> validMethod;
void initMethod()
{
  validMethod["get"] = true;
  validMethod["set"] = true;
  validMethod["cluster"] = true;
  validMethod["del"] = true;
}
/* max size packet to catch */
#define MAX_TCP_PACKET_SIZE 65535
/* time to wait to return packets */
#define TIME_DURATION 1000 * 3

/* ethernet headers are always exactly 14 bytes [1] */
#define SIZE_ETHERNET 14

/* IP header */
struct ip_struct
{
  u_char ip_vhl;                 /* version << 4 | header length >> 2 */
  u_char ip_tos;                 /* type of service */
  u_short ip_len;                /* total length */
  u_short ip_id;                 /* identification */
  u_short ip_off;                /* fragment offset field */
#define IP_RF 0x8000             /* reserved fragment flag */
#define IP_DF 0x4000             /* dont fragment flag */
#define IP_MF 0x2000             /* more fragments flag */
#define IP_OFFMASK 0x1fff        /* mask for fragmenting bits */
  u_char ip_ttl;                 /* time to live */
  u_char ip_p;                   /* protocol */
  u_short ip_sum;                /* checksum */
  struct in_addr ip_src, ip_dst; /* source and dest address */
};
#define IP_HL(ip) (((ip)->ip_vhl) & 0x0f)
#define IP_V(ip) (((ip)->ip_vhl) >> 4)

/* TCP header */
typedef u_int tcp_seq;
struct tcp_struct
{
  u_short th_sport; /* source port */
  u_short th_dport; /* destination port */
  tcp_seq th_seq;   /* sequence number */
  tcp_seq th_ack;   /* acknowledgement number */
  u_char th_offx2;  /* data offset, rsvd */
#define TH_OFF(th) (((th)->th_offx2 & 0xf0) >> 4)
  u_char th_flags;
#define TH_FIN 0x01
#define TH_SYN 0x02
#define TH_RST 0x04
#define TH_PUSH 0x08
#define TH_ACK 0x10
#define TH_URG 0x20
#define TH_ECE 0x40
#define TH_CWR 0x80
#define TH_FLAGS (TH_FIN | TH_SYN | TH_RST | TH_ACK | TH_URG | TH_ECE | TH_CWR)
  u_short th_win; /* window */
  u_short th_sum; /* checksum */
  u_short th_urp; /* urgent pointer */
};

/*
 * redis 方法统计
 */
void count_cmd(const u_char *payload, int len)
{
  cmd_mutex.lock();
  int i;
  const u_char *ch;
  ch = payload;
  char cmd[100] = "";
  int index = 0;
  if (len < 8)
    return;
  ch += 8;
  for (i = 7; i < len; i++)
  {
    if (isprint(*ch) && index < 100)
    {
      cmd[index] = (char)*ch;
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
struct pcap_pkthdr *header;
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
 * 统计输出
 */
void *printResult(void *ptr)
{
  for (;;)
  {
    cmd_mutex.lock();
    for (map<string, int>::reverse_iterator iter = cmdCount.rbegin(); iter != cmdCount.rend(); iter++)
      cout << iter->first << "  " << iter->second << endl;
    cmd_mutex.unlock();
    sleep(1);
  }
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
int main()
{
  initMethod();

  pthread_t pCount, pClean;
  int i, ret;
  //创建子线程，线程id为pId
  ret = pthread_create(&pCount, NULL, printResult, NULL);
  ret = pthread_create(&pClean, NULL, clean, NULL);

  if (ret != 0)
  {
    printf("create pthread error!\n");
    exit(1);
  }
  char errBuf[PCAP_ERRBUF_SIZE];
  pcap_t *device = pcap_open_live("en0", MAX_TCP_PACKET_SIZE, 1, TIME_DURATION, errBuf);

  if (!device)
  {
    printf("错误: pcap_open_live(): %s\n", errBuf);
    exit(1);
  }
  struct bpf_program filter;
  pcap_compile(device, &filter, "host 10.0.6.29 and dst port 6379", 1, 0);
  pcap_setfilter(device, &filter);

  pcap_loop(device, -1, getPacket, NULL);
  pcap_close(device);

  return 0;
}