// g++ udp_recvfrom.cpp -o udp_recvfrom �ő��v�����Ag++ -std=c++11�Ƃ��Ȃ��Ɠ����Ȃ�(���Ƃ�����)
// Windows �̏ꍇ�́Ag++ -g udp_recvfrom.cpp -o udp_recvfrom -lwsock32 -lws2_32

#include <stdio.h>
#include <string.h>
#include "simple_udp.h"

simple_udp udp0("0.0.0.0",12345);

int main(int argc, char **argv){
  udp0.udp_bind();
  while (1){
    std::string rdata=udp0.udp_recv();
    printf("recv:%s\n", rdata.c_str());
  }
  return 0;
}
