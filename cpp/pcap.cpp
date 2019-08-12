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
#include "main.h"

using namespace std;

struct pcap_pkthdr *header;
std::mutex cmd_mutex;
map<string, int> cmdCount;
map<string, bool> validMethod;
void initMethod()
{
  // char[]
  validMethod["set"] = true;
  validMethod["setnx"] = true;
  validMethod["setex"] = true;
  validMethod["psetex"] = true;
  validMethod["get"] = true;
  validMethod["getset"] = true;
  validMethod["strlen"] = true;
  validMethod["append"] = true;
  validMethod["setrange"] = true;
  validMethod["getrange"] = true;
  validMethod["incr"] = true;
  validMethod["incrby"] = true;
  validMethod["incrbyfloat"] = true;
  validMethod["decr"] = true;
  validMethod["decrby"] = true;
  validMethod["mset"] = true;
  validMethod["msetnx"] = true;
  validMethod["mget"] = true;
  // hash table
  validMethod["hset"] = true;
  validMethod["hsetnx"] = true;
  validMethod["hget"] = true;
  validMethod["hexists"] = true;
  validMethod["hdel"] = true;
  validMethod["hlen"] = true;
  validMethod["hstrlen"] = true;
  validMethod["hincrby"] = true;
  validMethod["hincrbyfloat"] = true;
  validMethod["hmset"] = true;
  validMethod["hmget"] = true;
  validMethod["hkeys"] = true;
  validMethod["hvals"] = true;
  validMethod["hgetall"] = true;
  validMethod["hscan"] = true;
  // list
  validMethod["lpush"] = true;
  validMethod["lpushx"] = true;
  validMethod["rpush"] = true;
  validMethod["rpushx"] = true;
  validMethod["lpop"] = true;
  validMethod["rpop"] = true;
  validMethod["rpoplpush"] = true;
  validMethod["lrem"] = true;
  validMethod["llen"] = true;
  validMethod["lindex"] = true;
  validMethod["linsert"] = true;
  validMethod["lset"] = true;
  validMethod["lrange"] = true;
  validMethod["ltrim"] = true;
  validMethod["blpop"] = true;
  validMethod["brpop"] = true;
  validMethod["brpoplpush"] = true;
  // set
  validMethod["sadd"] = true;
  validMethod["sismember"] = true;
  validMethod["spop"] = true;
  validMethod["srandmember"] = true;
  validMethod["srem"] = true;
  validMethod["smove"] = true;
  validMethod["scard"] = true;
  validMethod["sasmembersdd"] = true;
  validMethod["sscan"] = true;
  validMethod["sinter"] = true;
  validMethod["sinterstore"] = true;
  validMethod["sunion"] = true;
  validMethod["sunionstore"] = true;
  validMethod["sdiff"] = true;
  validMethod["sdiffstore"] = true;
  // zip list
  validMethod["zadd"] = true;
  validMethod["zscore"] = true;
  validMethod["zincrby"] = true;
  validMethod["zcard"] = true;
  validMethod["zcount"] = true;
  validMethod["zrange"] = true;
  validMethod["zrevrange"] = true;
  validMethod["zrangebyscore"] = true;
  validMethod["zrevrangebyscore"] = true;
  validMethod["zrank"] = true;
  validMethod["zrevrank"] = true;
  validMethod["zrem"] = true;
  validMethod["zremrangebyrank"] = true;
  validMethod["zremrangebyscore"] = true;
  validMethod["zrangebylex"] = true;
  validMethod["zlexcount"] = true;
  validMethod["zremrangebylex"] = true;
  validMethod["zscan"] = true;
  validMethod["zunionstore"] = true;
  validMethod["zinterstore"] = true;
  // bit map
  validMethod["setbit"] = true;
  validMethod["getbit"] = true;
  validMethod["bitcount"] = true;
  validMethod["bitpos"] = true;
  validMethod["bitop"] = true;
  validMethod["bitfield"] = true;
  // db
  validMethod["exists"] = true;
  validMethod["type"] = true;
  validMethod["rename"] = true;
  validMethod["renamenx"] = true;
  validMethod["move"] = true;
  validMethod["del"] = true;
  validMethod["randomkey"] = true;
  validMethod["dbsize"] = true;
  validMethod["keys"] = true;
  validMethod["scan"] = true;
  // ttl
  validMethod["expire"] = true;
  validMethod["expireat"] = true;
  validMethod["ttl"] = true;
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
  if (*ch != '*') // "*1\r\n$3\r\nget\r\n"
  {
    cmd_mutex.unlock();
    return;
  }
  int index = 0;
  for (int z = 0; z < 2; z++)
  {
    for (int i = 0; i < len; i++)
    {
      if (*ch == '\n')
      {
        ch++;
        index++;
        goto next;
      }
      else if (i == len)
      {
        break;
      }
      ch++;
      index++;
    }
  next:
    continue;
  }
  char cmd[100];
  memset(cmd, 0, 100);
  for (int i = 0; i < len - index - 2 && i < 100; i++)
  {
    if (*ch == '\r' || *ch == '\n')
    {
      break;
    }
    cmd[i] = ((char)*ch <= 'Z' && (char)*ch >= 'A') ? (char)*ch - ('Z' - 'z') : (char)*ch; // to lower
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
    // printf("   * Invalid IP header length: %u bytes\n", size_ip);
    return;
  }
  /*
   *  This packet is TCP.
   */
  if (ip->ip_p != IPPROTO_TCP)
  {
    // printf("   * Invalid TCP packet\n");
    return;
  }
  /* define/compute tcp header offset */
  tcp = (struct tcp_struct *)(packet + SIZE_ETHERNET + size_ip);
  size_tcp = TH_OFF(tcp) * 4;
  if (size_tcp < 20)
  {
    // printf("   * Invalid TCP header length: %u bytes\n", size_tcp);
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
    {
      cmdCount.erase(iter->first);
      iter++;
    }
    cmd_mutex.unlock();
    sleep(TIME); // 每分钟清除一次
  }
}
